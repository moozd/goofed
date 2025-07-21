package gfx

import "github.com/go-gl/gl/v4.1-core/gl"

type VBO struct {
	id       uint32
	vertices []float32
}

func NewVBO(vertices []float32) *VBO {
	vbo := &VBO{}
	vbo.vertices = vertices

	gl.GenBuffers(1, &vbo.id)
	diagnose()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.id)
	diagnose()
	gl.BufferData(gl.ARRAY_BUFFER, F32.SizeOf(len(vertices)), gl.Ptr(vertices), gl.STATIC_DRAW)
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
