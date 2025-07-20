package screen

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

//go:embed assets/frag.glsl
var fragShaderSrc string

//go:embed assets/vert.glsl
var vertShaderSrc string

func (self *Screen) Render() {
	surface := NewSurface()

	shader := NewShader(vertShaderSrc, fragShaderSrc)

	vertices := []float32{
		//               COORDINATES                  //     COLORS
		-0.5, float32(-0.5 * math.Sqrt(3) * 1 / 3), 0.0, 0.8, 0.3, 0.02, // Lower left corner
		0.5, float32(-0.5 * math.Sqrt(3) * 1 / 3), 0.0, 0.8, 0.3, 0.02, // Lower right corner
		0.0, float32(0.5 * math.Sqrt(3) * 2 / 3), 0.0, 1.0, 0.6, 0.32, // Upper corner
		-0.25, float32(0.5 * math.Sqrt(3) * 1 / 6), 0.0, 0.9, 0.45, 0.17, // Inner left
		0.25, float32(0.5 * math.Sqrt(3) * 1 / 6), 0.0, 0.9, 0.45, 0.17, // Inner right
		0.0, float32(-0.5 * math.Sqrt(3) * 1 / 3), 0.0, 0.8, 0.3, 0.02, // Inner down
	}

	indices := []uint32{
		0, 3, 5,
		3, 2, 4,
		5, 4, 1,
	}

	vao := NewVAO()
	vao.Bind()

	vbo := NewVBO(vertices)
	ebo := NewEBO(indices)

	vao.LinkVBO(vbo, 0, 3, GLfloat, 6, 0)
	vao.LinkVBO(vbo, 1, 3, GLfloat, 6, 3)

	vao.Unbind()
	vbo.Unbind()
	ebo.Unbind()

	surface.Loop(func() {

		gl.ClearColor(0.07, 0.13, 0.17, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shader.Use()
		vao.Bind()

		gl.DrawElements(gl.TRIANGLES, 9, gl.UNSIGNED_INT, gl.PtrOffset(0))

	})

	shader.Delete()
	vbo.Delete()
	vao.Delete()
	ebo.Delete()

}

// -----------GL HELPERS-------------

type Surface struct {
	win  *sdl.Window
	gctx sdl.GLContext
}

func NewSurface() *Surface {
	s := &Surface{}
	s.init()
	return s
}

func (s *Surface) Loop(cbl func()) {

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

		cbl()

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

func diagnose() {
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

// -------VBO---------

type VBO struct {
	id       uint32
	vertices []float32
}

func NewVBO(vertices []float32) *VBO {
	vbo := &VBO{}
	vbo.vertices = vertices
	size := GLTypeToSize[GLfloat]

	gl.GenBuffers(1, &vbo.id)
	diagnose()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)
	diagnose()
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*size, gl.Ptr(vertices), gl.STATIC_DRAW)
	diagnose()

	return vbo
}

func (v *VBO) ID() uint32 {
	return v.id
}

func (v *VBO) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, v.id)
	diagnose()
}
func (v *VBO) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	diagnose()
}
func (v *VBO) Delete() {
	gl.DeleteBuffers(1, &v.id)
	diagnose()
}

//----------

type VAO struct {
	id uint32
}

func NewVAO() *VAO {
	vao := &VAO{}

	gl.GenVertexArrays(1, &vao.id)
	diagnose()

	return vao
}

func (v *VAO) LinkVBO(vbo *VBO, layout uint32, numComponents int32, glType GLType, stride int32, offset int) {

	gsize := GLTypeToSize[glType]
	gtype := GLTypeToGLEnum[glType]

	vbo.Bind()

	gl.VertexAttribPointerWithOffset(layout, numComponents, gtype, false, stride*int32(gsize), uintptr(offset*gsize))
	diagnose()
	gl.EnableVertexAttribArray(layout)
	diagnose()

	vbo.Unbind()
}
func (v *VAO) Bind() {
	gl.BindVertexArray(v.id)
	diagnose()
}
func (v *VAO) Unbind() {
	gl.BindVertexArray(0)
	diagnose()
}
func (v *VAO) Delete() {
	gl.DeleteVertexArrays(1, &v.id)
	diagnose()
}

type EBO struct {
	id      uint32
	indices []uint32
}

func NewEBO(indices []uint32) *EBO {
	ebo := &EBO{}
	ebo.indices = indices

	gl.GenBuffers(1, &ebo.id)
	diagnose()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo.id)
	diagnose()
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	diagnose()

	return ebo
}

func (v *EBO) ID() uint32 {
	return v.id
}

func (v *EBO) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, v.id)
	diagnose()
}
func (v *EBO) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	diagnose()
}
func (v *EBO) Delete() {
	gl.DeleteBuffers(1, &v.id)
	diagnose()
}

// --------Shader -----------------

type Shader struct {
	id uint32
}

func NewShader(vertSrc, fragSrc string) *Shader {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	cSources, free := gl.Strs(vertSrc + "\x00")
	gl.ShaderSource(vertexShader, 1, cSources, nil)
	free()
	gl.CompileShader(vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cFrag, freeFrag := gl.Strs(fragSrc + "\x00")
	gl.ShaderSource(fragmentShader, 1, cFrag, nil)
	freeFrag()
	gl.CompileShader(fragmentShader)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &Shader{id: shaderProgram}
}

func (s *Shader) Use() {
	gl.UseProgram(s.id)
}

func (s *Shader) Id() uint32 {
	return s.id
}

func (s *Shader) Delete() {
	gl.DeleteProgram(s.id)
}

// GLType enum
type GLType int

const (
	GLbyte GLType = iota
	GLubyte
	GLshort
	GLushort
	GLint
	GLuint
	GLfloat
	GLdouble
)

// Map: GLType -> size in bytes
var GLTypeToSize = map[GLType]int{
	GLbyte:   1,
	GLubyte:  1,
	GLshort:  2,
	GLushort: 2,
	GLint:    4,
	GLuint:   4,
	GLfloat:  4,
	GLdouble: 8,
}

// Map: GLType -> OpenGL GLenum (e.g. gl.FLOAT)
var GLTypeToGLEnum = map[GLType]uint32{
	GLbyte:   gl.BYTE,
	GLubyte:  gl.UNSIGNED_BYTE,
	GLshort:  gl.SHORT,
	GLushort: gl.UNSIGNED_SHORT,
	GLint:    gl.INT,
	GLuint:   gl.UNSIGNED_INT,
	GLfloat:  gl.FLOAT,
	GLdouble: gl.DOUBLE,
}
