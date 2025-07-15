package gpu

import (
	"context"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Gpu struct {
	ctx    context.Context
	config *Config
	window *sdl.Window
	glc    sdl.GLContext
	models []Renderable
}

func New(ctx context.Context, c Config) *Gpu {
	return &Gpu{
		ctx:    ctx,
		config: &c,
	}
}

func (self *Gpu) Run() error {
	err := self.initWindow()
	if err != nil {
		return err
	}

	defer self.cleanup()

	err = self.mainLoop()
	if err != nil {
		return err
	}

	return nil
}

func (self *Gpu) Add(m Renderable) {
	self.models = append(self.models, m)
}

func (self *Gpu) initWindow() error {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	// Configure OpenGL version and profile
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	win, err := sdl.CreateWindow(
		self.config.WindowTitle,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(self.config.WindowWidth),
		int32(self.config.WindowHeight),
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)

	if err != nil {
		return err
	}

	glContext, err := win.GLCreateContext()
	if err != nil {
		return err
	}

	self.glc = glContext
	self.window = win

	return nil

}

func (self *Gpu) mainLoop() error {

	if err := gl.Init(); err != nil {
		return err
	}

	gl.Viewport(0, 0, int32(self.config.WindowWidth), int32(self.config.WindowHeight))

	for _, model := range self.models {
		model.Init(self)
		model.Setup()
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyboardEvent := event.(*sdl.KeyboardEvent)

				for _, model := range self.models {
					model.OnKeyboardEvent(keyboardEvent)
				}
			}
		}

		for _, model := range self.models {
			model.Render()
		}

		self.window.GLSwap()

	}

	return nil
}

func (self *Gpu) cleanup() {
	self.window.Destroy()
	sdl.GLDeleteContext(self.glc)
	sdl.Quit()
}
