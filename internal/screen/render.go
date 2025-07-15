package screen

import (
	_ "embed"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/moozd/goofed/internal/gpu/shader"
	"github.com/moozd/goofed/internal/screen/grid"
)

type gpuContext struct {
	vao      uint32
	vbo      uint32
	shader   *shader.Shader
	vertices []float32
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

	self.gpuContext = &gpuContext{
		vao:    vao,
		vbo:    vbo,
		shader: shader.New(vertSrc, fragSrc),
	}

	self.handleSize()
	self.grid.ResetViewOffset()
	self.updateGridVerts()

	self.setupGridAttrs()
	self.gpuContext.shader.Use()

}

func (self *Screen) Render() {
	self.handleSize()

	self.updateGridVerts()

	gl.BindVertexArray(self.gpuContext.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(self.gpuContext.vertices)/5))
}

func (self *Screen) updateGridVerts() {

	var vertices []float32
	cols := self.grid.Size.Cols
	rows := self.grid.Size.Rows
	cw := float32(self.grid.CellSize.Width) / float32(cols) // OpenGL NDC: -1 to 1
	ch := float32(self.grid.CellSize.Height) / float32(rows)

	self.grid.GetView(grid.ViewDirty, func(x, y int, cell *grid.Cell) {
		cell.Dirty = false
		x0 := -1 + float32(x)*cw
		y0 := -1 + float32(y)*ch
		x1 := x0 + cw
		y1 := y0 + ch

		r, g, b := float32(0), float32(0), float32(0)
		if y%2 == 0 {
			r, g, b = float32(0), float32(0), float32(0)
			if x%2 == 0 {
				r, g, b = float32(0xff), float32(0xff), float32(0xff)
			}
		} else {

			r, g, b = float32(0), float32(0), float32(0)
			if x%2 != 0 {
				r, g, b = float32(0xff), float32(0xff), float32(0xff)
			}
		}

		// two triangles
		vertices = append(vertices,
			x0, y0, r, g, b,
			x1, y0, r, g, b,
			x1, y1, r, g, b,
			x0, y0, r, g, b,
			x1, y1, r, g, b,
			x0, y1, r, g, b,
		)
	},
	)

	self.gpuContext.vertices = vertices

	gl.BindVertexArray(self.gpuContext.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, self.gpuContext.vbo) // you'll need to store vbo
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
}

func (self *Screen) setupGridAttrs() {
	gl.BindVertexArray(self.gpuContext.vao)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(8))
	gl.EnableVertexAttribArray(1)
}

func (self *Screen) handleSize() {
	w, h := self.Window().GetSize()
	self.grid.Resize(w, h, 10, 10)
}
