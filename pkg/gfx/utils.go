package gfx

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
)

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
		errCode, asGLErrorCode(errCode), file, line, fnName)
}

func assert[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}

func asGLErrorCode(err uint32) string {
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
