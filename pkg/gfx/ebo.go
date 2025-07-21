package gfx

import "github.com/go-gl/gl/v4.1-core/gl"

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
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, U32.SizeOf(len(indices)), gl.Ptr(indices), gl.STATIC_DRAW)
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
