package gfx

import "github.com/go-gl/gl/v4.1-core/gl"

type VAO struct {
	id uint32
}

func NewVAO() *VAO {
	vao := &VAO{}

	gl.GenVertexArrays(1, &vao.id)
	diagnose()

	vao.Bind()

	return vao
}

func (v *VAO) Define(vbo *VBO, kind *GFXType, layout uint32, numComponents int32, stride int, offset int) {

	normalized := false

	if kind == U8 || kind == I8 {
		normalized = true
	}

	vbo.Bind()

	gl.VertexAttribPointerWithOffset(layout, numComponents, kind.glType, normalized, int32(stride), uintptr(offset))
	diagnose()
	gl.EnableVertexAttribArray(layout)
	diagnose()

	vbo.Unbind()
}

func (v *VAO) Draw(ebo *EBO) {
	gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(ebo.indices)), gl.UNSIGNED_INT, 0)
	diagnose()
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
