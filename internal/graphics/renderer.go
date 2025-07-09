package graphics

import (
	"context"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Renderable interface {
	Init(r *Renderer)
	Draw(r *Renderer)
}

type Renderer struct {
	height      int
	width       int
	title       string
	ctx         context.Context
	renderables []Renderable
}

func NewRenderer(ctx context.Context, width, height int) *Renderer {
	return &Renderer{
		ctx:    ctx,
		width:  width,
		height: height,
		title:  "goofed",
	}
}

func (r *Renderer) Add(rnd Renderable) {
	r.renderables = append(r.renderables, rnd)

}

func (r *Renderer) Start() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(r.width, r.height, r.title, nil, nil)

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	for _, rnd := range r.renderables {
		rnd.Init(r)
	}
	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.2, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		window.SwapBuffers()
		glfw.PollEvents()

		for _, rnd := range r.renderables {
			rnd.Draw(r)
		}

	}

}

func (r *Renderer) Close() {}
