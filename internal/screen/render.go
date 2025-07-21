package screen

import (
	_ "embed"
	"image/color"
	"log"

	"github.com/moozd/goofed/pkg/gfx"
)

var (
	//go:embed assets/frag.glsl
	fragShaderSrc string

	//go:embed assets/vert.glsl
	vertShaderSrc string
)

func (self *Screen) Render() {
	surface := gfx.NewSurface()
	surface.SetBackground(color.RGBA{R: 0x0, G: 0x0, B: 0x0})

	shader := gfx.NewShader(vertShaderSrc, fragShaderSrc)

	w, h := surface.Size()

	vertices, indices := createGrid(w, h, 20, 30)

	vao := gfx.NewVAO(gfx.F32.SizeOf(3 + 2 + 3 + 3))
	vbo := gfx.NewVBO(vertices)
	ebo := gfx.NewEBO(indices)

	vao.Define(vbo, gfx.F32, 0, 3, 0)                 // pos
	vao.Define(vbo, gfx.F32, 1, 2, gfx.F32.SizeOf(3)) // uv
	vao.Define(vbo, gfx.F32, 2, 3, gfx.F32.SizeOf(5)) // fg
	vao.Define(vbo, gfx.F32, 3, 3, gfx.F32.SizeOf(8)) // bg

	vao.Unbind()
	vbo.Unbind()
	ebo.Unbind()

	surface.OnResize(func(w, h int32) {
		log.Printf("Window size changed to %dx%d", w, h)

		shader.Use()
		shader.SetMat4("projection", surface.Projection)
	})

	surface.Loop(func() {
		shader.Use()
		vao.Bind()
		vao.Draw(ebo)
	})

	shader.Delete()
	vbo.Delete()
	vao.Delete()
	ebo.Delete()

}

func createGrid(ww, wh int32, m, n int) (vertices []float32, indices []uint32) {

	tc := uint32(0)
	cw := float32(ww) / float32(m) // width per column
	ch := float32(wh) / float32(n) // height per row
	for y := range n {
		for x := range m {

			l := float32(x) * cw
			r := float32(x+1) * cw
			b := float32(y+1) * ch
			t := float32(y) * ch

			fgr, fgg, fgb := 0.5, 0.5, 0.5
			bgr, bgg, bgb := 0.1, 0.1, 0.1

			vertices = append(vertices, []float32{
				// pos   // uv  								 // fg           												   // bg
				l, b, 0, 0.0, 0.0, float32(fgr), float32(fgg), float32(fgb), float32(bgr), float32(bgg), float32(bgb),
				r, b, 0, 1.0, 0.0, float32(fgr), float32(fgg), float32(fgb), float32(bgr), float32(bgg), float32(bgb),
				l, t, 0, 0.0, 1.0, float32(fgr), float32(fgg), float32(fgb), float32(bgr), float32(bgg), float32(bgb),
				r, t, 0, 1.0, 1.0, float32(fgr), float32(fgg), float32(fgb), float32(bgr), float32(bgg), float32(bgb),
			}...,
			)

			indices = append(indices, []uint32{
				tc, tc + 1, tc + 2,
				tc + 1, tc + 2, tc + 3,
			}...)
			tc += 4

		}
	}

	return
}
