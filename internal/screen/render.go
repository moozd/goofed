package screen

import (
	_ "embed"
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func (self *Screen) Loop() {
	g := gState{scr: self}
	g.init()
	defer g.cleanUp()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		}

		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		g.render()
		g.win.GLSwap()
	}
}

type gState struct {
	scr    *Screen
	win    *sdl.Window
	gctx   sdl.GLContext
	shader *shader
}

//go:embed assets/frag.glsl
var fragShaderSrc string

//go:embed assets/vert.glsl
var vertShaderSrc string

func (gs *gState) render() {

	gs.shader.use()
}

func (gs *gState) init() {
	title := "goofed"
	width := int32(800)
	height := int32(600)

	assert(0, sdl.Init(sdl.INIT_VIDEO))
	log.Printf("SDL initilized")

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	log.Printf("SDL gl attrs initilized")

	win := assert(sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE))
	log.Printf("SDL window created.")

	ctx := assert(win.GLCreateContext())
	log.Printf("SDL window gl context created.")

	assert(0, win.GLMakeCurrent(ctx))
	log.Printf("SDL window switch to gl conext.")

	gl.Init()
	gs.diagnose()

	gl.Enable(gl.BLEND)
	gs.diagnose()

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gs.diagnose()

	gs.win = win
	gs.gctx = ctx
	gs.shader = newShader(vertShaderSrc, fragShaderSrc)

}

func (gs *gState) cleanUp() {
	log.Printf("Cleaning up....")
	// gl.DeleteProgram(r.Shader)
	// gl.DeleteVertexArrays(1, &r.VAO)
	// gl.DeleteBuffers(1, &r.VBO)
	// gl.DeleteTextures(1, &r.Texture)
	sdl.GLDeleteContext(gs.gctx)
	gs.win.Destroy()
	sdl.Quit()
}

func (gs *gState) diagnose() {
	errCode := gl.GetError()
	if errCode == gl.NO_ERROR {
		return
	}

	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("OpenGL error: 0x%x (caller info unavailable)\n", errCode)
		return
	}

	fn := runtime.FuncForPC(pc)
	fnName := "unknown"
	if fn != nil {
		fnName = fn.Name()
	}

	fmt.Printf("OpenGL error 0x%x (%s) at %s:%d (in %s)\n",
		errCode, getGlErrorCode(errCode), file, line, fnName)
}
