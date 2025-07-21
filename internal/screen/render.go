package screen

import (
	_ "embed"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/moozd/goofed/internal/gfx"
)

var (
	//go:embed assets/frag.glsl
	fragShaderSrc string

	//go:embed assets/vert.glsl
	vertShaderSrc string
)

func (self *Screen) Render() {
	surface := gfx.NewSurface()

	shader := gfx.NewShader(vertShaderSrc, fragShaderSrc)

	vertices := []float32{
		-0.5, float32(-0.5 * math.Sqrt(3) * 1 / 3), 0.0, 0.8, 0.3, 0.02,
		0.5, float32(-0.5 * math.Sqrt(3) * 1 / 3), 0.0, 0.8, 0.3, 0.02,
		0.0, float32(0.5 * math.Sqrt(3) * 2 / 3), 0.0, 1.0, 0.6, 0.32,
		-0.25, float32(0.5 * math.Sqrt(3) * 1 / 6), 0.0, 0.9, 0.45, 0.17,
		0.25, float32(0.5 * math.Sqrt(3) * 1 / 6), 0.0, 0.9, 0.45, 0.17,
		0.0, float32(-0.5 * math.Sqrt(3) * 1 / 3), 0.0, 0.8, 0.3, 0.02,
	}

	indices := []uint32{
		0, 3, 5,
		3, 2, 4,
		5, 4, 1,
	}

	vao := gfx.NewVAO()
	vbo := gfx.NewVBO(vertices)
	ebo := gfx.NewEBO(indices)

	vao.Define(vbo, gfx.F32, 0, 3, 6*gfx.F32.Size, 0)
	vao.Define(vbo, gfx.F32, 1, 3, 6*gfx.F32.Size, 3*gfx.F32.Size)

	vao.Unbind()
	vbo.Unbind()
	ebo.Unbind()

	surface.Loop(func() {

		gl.ClearColor(0.07, 0.13, 0.17, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shader.Use()
		vao.Bind()

		vao.Draw(ebo)

	})

	shader.Delete()
	vbo.Delete()
	vao.Delete()
	ebo.Delete()

}
