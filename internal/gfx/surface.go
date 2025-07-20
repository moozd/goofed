package gfx

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Surface struct {
	win  *sdl.Window
	gctx sdl.GLContext
}

func NewSurface() *Surface {
	s := &Surface{}
	s.init()
	return s
}

func (s *Surface) Loop(fn func()) {

	defer s.cleanUp()
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

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

		fn()

		s.win.GLSwap()
	}
}

func (gs *Surface) init() {
	title := "goofed"
	width := int32(800)
	height := int32(800)

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
	diagnose()

	gl.Enable(gl.BLEND)
	diagnose()

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	diagnose()

	gl.Viewport(0, 0, width, height)
	diagnose()

	gs.win = win
	gs.gctx = ctx

}

func (gs *Surface) cleanUp() {
	log.Printf("Gl Cleaning up....")
	sdl.GLDeleteContext(gs.gctx)
	gs.win.Destroy()
	sdl.Quit()
}
