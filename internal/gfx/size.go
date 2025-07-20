package gfx

import "github.com/go-gl/gl/v4.1-core/gl"

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
var glTypeToSize = map[GLType]int{
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
var glTypeToGLEnum = map[GLType]uint32{
	GLbyte:   gl.BYTE,
	GLubyte:  gl.UNSIGNED_BYTE,
	GLshort:  gl.SHORT,
	GLushort: gl.UNSIGNED_SHORT,
	GLint:    gl.INT,
	GLuint:   gl.UNSIGNED_INT,
	GLfloat:  gl.FLOAT,
	GLdouble: gl.DOUBLE,
}

func GetType(t GLType) (size int, glType uint32) {
	size = glTypeToSize[t]
	glType = glTypeToGLEnum[t]
	return
}
