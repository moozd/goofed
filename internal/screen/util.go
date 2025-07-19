package screen

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func slice[T any](s []T, start, end int) []T {
	if start < 0 {
		start = 0
	}
	if start >= len(s) {
		return []T{}
	}
	if end > len(s) {
		end = len(s)
	}
	if end < start {
		return []T{}
	}
	return s[start:end]
}

func assert[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}

func getGlErrorCode(err uint32) string {
	switch err {
	case gl.NO_ERROR:
		return "NO_ERROR"
	case gl.INVALID_ENUM:
		return "INVALID_ENUM"
	case gl.INVALID_VALUE:
		return "INVALID_VALUE"
	case gl.INVALID_OPERATION:
		return "INVALID_OPERATION"
	case gl.OUT_OF_MEMORY:
		return "OUT_OF_MEMORY"
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		return "INVALID_FRAMEBUFFER_OPERATION"
	default:
		return fmt.Sprintf("Unknown error: 0x%x", err)
	}
}

type shader struct {
	id uint32
}

func (s *shader) Id() uint32 {
	return s.id
}

func newShader(vertSrc, fragSrc string) *shader {
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

	return &shader{id: shaderProgram}
}

func (s *shader) use() {
	gl.UseProgram(s.id)
}
