package gfx

import "github.com/go-gl/gl/v4.1-core/gl"

type VAO struct {
	id uint32
}

func NewVAO() *VAO {
	vao := &VAO{}

	gl.GenVertexArrays(1, &vao.id)
	diagnose()

	return vao
}

func (v *VAO) Define(vbo *VBO, layout uint32, numComponents int32, glType GLType, stride int32, offset int) {

	gsize, gtype := GetType(glType)

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
