package gpu

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Renderable interface {
	Init(d *Gpu)
	Setup()
	OnKeyboardEvent(event *sdl.KeyboardEvent)
	Render()
}

type Model struct {
	gpu *Gpu
}

func (self *Model) Init(d *Gpu) {
	self.gpu = d
}

func (self *Model) Setup() {}

func (self *Model) Render() {}

func (self *Model) OnKeyboardEvent(event *sdl.KeyboardEvent) {}

func (self *Model) Config() *Config {
	return self.gpu.config
}

func (self *Model) Window() *sdl.Window {
	return self.gpu.window
}

// type props struct {
// 	vao             uint32
// 	count           int32
// 	cellSizeLoc     int32
// 	wireCellSizeLoc int32
// 	shader1         uint32
// 	shader2         uint32
// }
//
// const (
// 	cols = 80
// 	rows = 20
// )
//
// var vertexShaderSrc = `
// #version 410 core
// layout (location = 0) in vec2 aPos;
// layout (location = 1) in vec2 aOffset;
// uniform vec2 cellSize;
// void main() {
//     vec2 scaled = aPos * cellSize;
//     vec2 pos = scaled + aOffset;
//     gl_Position = vec4(pos, 0.0, 1.0);
// }
// ` + "\x00"
//
// var filledFragShader = `
// #version 410 core
// out vec4 FragColor;
// void main() {
//     FragColor = vec4(0.0, 0.0, 0.0, 1.0); // black fill
// }` + "\x00"
//
// var wireFragShader = `
// #version 410 core
// out vec4 FragColor;
// void main() {
//     FragColor = vec4(1.0); // white lines
// }` + "\x00"
//
// func (self *Display) setup() *props {
//
// 	shader1 := createProgram(vertexShaderSrc, filledFragShader)
// 	shader2 := createProgram(vertexShaderSrc, wireFragShader)
//
// 	vao, instanceCount := createQuadGrid(cols, rows)
//
// 	return &props{
// 		vao:     vao,
// 		count:   instanceCount,
// 		shader1: shader1,
// 		shader2: shader2,
// 	}
// }
//
// func (self *Display) render(p *props) {
// 	width, height := self.window.GetSize()
// 	aspect := float32(width) / float32(height)
//
// 	gl.Viewport(0, 0, int32(width), int32(height))
// 	gl.ClearColor(0.1, 0.1, 0.1, 1)
// 	gl.Clear(gl.COLOR_BUFFER_BIT)
//
// 	cellW := 2.0 / float32(cols)
// 	cellH := 2.0 / float32(rows)
// 	cellSize := [2]float32{
// 		cellW / aspect,
// 		cellH,
// 	}
//
// 	// Filled background
// 	gl.UseProgram(p.shader1)
// 	gl.Uniform2f(p.cellSizeLoc, cellSize[0], cellSize[1])
// 	gl.BindVertexArray(p.vao)
// 	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil, p.count)
//
// 	// Wireframe overlay
// 	gl.Enable(gl.POLYGON_OFFSET_LINE)
// 	gl.PolygonOffset(-1, -1)
// 	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
//
// 	gl.UseProgram(p.shader2)
// 	gl.Uniform2f(p.wireCellSizeLoc, cellSize[0], cellSize[1])
// 	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil, p.count)
//
// 	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
// 	gl.Disable(gl.POLYGON_OFFSET_LINE)
//
// }
