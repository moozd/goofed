package shader

import "github.com/go-gl/gl/v4.1-core/gl"

type Shader struct {
	id uint32
}

func (self *Shader) Id() uint32 {
	return self.id
}

func New(vertSrc, fragSrc string) *Shader {
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
