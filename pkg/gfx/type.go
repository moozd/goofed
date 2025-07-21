package gfx

import "github.com/go-gl/gl/v4.1-core/gl"

// Map: GLType -> size in bytes
var glTypeSizes = map[uint32]int{
	gl.BYTE:           1,
	gl.UNSIGNED_BYTE:  1,
	gl.SHORT:          2,
	gl.UNSIGNED_SHORT: 2,
	gl.INT:            4,
	gl.UNSIGNED_INT:   4,
	gl.FLOAT:          4,
	gl.DOUBLE:         8,
}

type GFXType struct {
	size   int
	glType uint32
}

func newGFXType(openGLType uint32) *GFXType {
	size := glTypeSizes[openGLType]
	return &GFXType{
		size:   size,
		glType: openGLType,
	}
}

func (t *GFXType) SizeOf(count int) int {
	return t.size * count
}

var (
	I8  = newGFXType(gl.BYTE)
	U8  = newGFXType(gl.UNSIGNED_BYTE)
	I16 = newGFXType(gl.SHORT)
	U16 = newGFXType(gl.UNSIGNED_SHORT)
	I32 = newGFXType(gl.INT)
	U32 = newGFXType(gl.UNSIGNED_INT)
	F32 = newGFXType(gl.FLOAT)
	F64 = newGFXType(gl.DOUBLE)
)
