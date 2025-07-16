package screen

import (
	_ "embed"
	"image"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/moozd/goofed/internal/gpu/math"
	"github.com/moozd/goofed/internal/gpu/shader"
	"github.com/moozd/goofed/internal/screen/grid"
)

type gpuContext struct {
	vao        uint32
	vbo        uint32
	wh         int32
	ww         int32
	tex        uint32
	shader     *shader.Shader
	vertices   []float32
	projection [16]float32
}

var (
	//go:embed assets/vert.glsl
	vertSrc string

	//go:embed assets/frag.glsl
	fragSrc string
)

func (self *Screen) Setup() {

	var vbo, vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	self.gctx = &gpuContext{
		vao:    vao,
		vbo:    vbo,
		shader: shader.New(vertSrc, fragSrc),
	}

	self.handleSize()
	self.grid.ResetViewOffset()

	self.setupTexture()
	self.setupGrid()

}

func (self *Screen) Render() {
	self.handleSize()

	self.updateGrid()

	gl.BindVertexArray(self.gctx.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(self.gctx.vertices)/5))
}

func (self *Screen) updateGrid() {

	uProjection := gl.GetUniformLocation(self.gctx.shader.Id(), gl.Str("uProjection\x00"))

	self.gctx.shader.Use()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, self.gctx.tex)

	gl.UniformMatrix4fv(uProjection, 1, false, &self.gctx.projection[0])

	uFontAtlas := gl.GetUniformLocation(self.gctx.shader.Id(), gl.Str("uFontAtlas\x00"))
	gl.Uniform1i(uFontAtlas, 0)

	if self.grid.IsClean() {
		return
	}
	var vertices []float32

	cw := self.grid.CellSize.Width
	ch := self.grid.CellSize.Height
	fnt := self.Config().Font

	self.grid.GetView(grid.ViewModeRender, func(x, y int, cell *grid.Cell) {

		g := fnt.Glyph(cell.Rune)

		fnt.Atlas().Bounds()
		log.Printf("x=%d, y=%d, rune=%c", x, y, cell.Rune)
		cw = g.Width
		ch = g.Height

		x0 := float32(x * cw)
		y0 := float32(y * ch)
		x1 := x0 + float32(cw)
		y1 := y0 + float32(ch)

		fgR, fgG, fgB := float32(1.0), float32(1.0), float32(1.0)
		bgR, bgG, bgB := float32(0), float32(0), float32(0)

		u0, v0, u1, v1 := g.U0, g.V0, g.U1, g.V1

		// two triangles
		vertices = append(vertices,
			x0, y0, u0, v0, fgR, fgG, fgB, bgR, bgG, bgB,
			x1, y0, u1, v0, fgR, fgG, fgB, bgR, bgG, bgB,
			x1, y1, u1, v1, fgR, fgG, fgB, bgR, bgG, bgB,
			x0, y0, u0, v0, fgR, fgG, fgB, bgR, bgG, bgB,
			x1, y1, u1, v1, fgR, fgG, fgB, bgR, bgG, bgB,
			x0, y1, u0, v1, fgR, fgG, fgB, bgR, bgG, bgB,
		)
	},
	)

	self.gctx.vertices = vertices

	gl.BindVertexArray(self.gctx.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, self.gctx.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
}

func (self *Screen) setupTexture() {
	img := self.Config().Font.Atlas().(*image.RGBA)
	gl.GenTextures(1, &self.gctx.tex)
	gl.BindTexture(gl.TEXTURE_2D, self.gctx.tex)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1) // in case row alignment matters

	w, h := int32(img.Rect.Dx()), int32(img.Rect.Dy())

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA, // internal format
		w,
		h,
		0,
		gl.RGBA, // source format
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix),
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

}

func (self *Screen) setupGrid() {

	self.updateGrid()

	gl.BindVertexArray(self.gctx.vao)

	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 40, 0)  // aPos x,y
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 40, 8)  // aUV  u,v
	gl.VertexAttribPointerWithOffset(2, 3, gl.FLOAT, false, 40, 16) // aFGColor fg
	gl.VertexAttribPointerWithOffset(3, 3, gl.FLOAT, false, 40, 28) // aBGColor bg

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
}

func (self *Screen) handleSize() {
	w, h := self.Window().GetSize()
	if w == self.gctx.ww && h == self.gctx.wh {
		return
	}
	self.gctx.ww = w
	self.gctx.wh = h
	gl.Viewport(0, 0, int32(w), int32(h))

	cw, ch := self.Config().Font.Measure('A')
	self.grid.Resize(w, h, int32(cw), int32(ch))
	self.gctx.projection = math.Ortho(0, float32(w), float32(h), 0, -1, 1)

	log.Printf("win size w=%d h=%d rows=%d cols=%d ", w, h, self.grid.Size.Rows, self.grid.Size.Cols)
}
