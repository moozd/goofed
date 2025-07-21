package gfx

import (
	"image/color"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type ResizeHandler = func(w, h int32)

type Surface struct {
	win           *sdl.Window
	gctx          sdl.GLContext
	resizeHandler ResizeHandler
	bg            color.RGBA
	Projection    mgl32.Mat4
}

func NewSurface() *Surface {
	s := &Surface{}
	s.init()
	return s
}

func (s *Surface) Size() (w, h int32) { w, h = s.win.GetSize(); return }

func (s *Surface) OnResize(fn ResizeHandler) { s.resizeHandler = fn }

func (s *Surface) Loop(fn func()) {

	defer s.cleanUp()
	s.handleResize()

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
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					s.handleResize()
				}

			}
		}

		s.drawBackground()
		fn()
		s.win.GLSwap()
	}
}

func (s *Surface) SetBackground(c color.RGBA) {
	s.bg = c
}

func (s *Surface) drawBackground() {
	r, g, b, a := toOpenGLColor(s.bg)
	gl.ClearColor(r, g, b, a)
	diagnose()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	diagnose()
}

func (s *Surface) handleResize() {
	w, h := s.win.GetSize()
	log.Printf("SDL window size changed %dx%d\n", w, h)

	gl.Viewport(0, 0, w, h)
	s.Projection = mgl32.Ortho2D(0, float32(w), float32(h), 0)
	if s.resizeHandler != nil {
		s.resizeHandler(w, h)
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
	log.Printf("SDL bye.")
	sdl.GLDeleteContext(gs.gctx)
	gs.win.Destroy()
	sdl.Quit()
}
