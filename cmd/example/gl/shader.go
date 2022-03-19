package gl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderKind uint32

const (
	VertexShaderKind   ShaderKind = gl.VERTEX_SHADER
	FragmentShaderKind ShaderKind = gl.FRAGMENT_SHADER
)

type Shader struct {
	Id     uint32
	Kind   ShaderKind
	Source string
}

func NewShader(kind ShaderKind, source string) (shader *Shader, err error) {
	shader = CreateShader(kind)
	shader.SetSources(source)
	shader.Compile()

	if success := shader.GetCompileStatus(); !success {
		logLength := shader.GetLogInfoLength()
		log := shader.GetInfoLog(logLength)
		shader = nil
		err = fmt.Errorf("failed to compile shader %v: %v", source, log)
	}
	return
}

func (shader *Shader) Exists() bool {
	return shader.Id != 0
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.Id)
	shader.Id = 0
}

// GL Functions

func CreateShader(kind ShaderKind) (shader *Shader) {
	shader = new(Shader)
	shader.Id = gl.CreateShader(uint32(kind))
	return
}

func (shader *Shader) Compile() {
	gl.CompileShader(shader.Id)
}

func (shader *Shader) SetSources(sources ...string) {
	glSources, free := glStrings(sources...)
	gl.ShaderSource(shader.Id, int32(len(sources)), glSources, nil)
	free()
}

type ShaderProperty uint32

const (
	CompileStatusShaderProperty ShaderProperty = gl.COMPILE_STATUS
	LogInfoLengthShaderProperty ShaderProperty = gl.INFO_LOG_LENGTH
)

func (shader *Shader) Getiv(property ShaderProperty) (value int32) {
	gl.GetShaderiv(shader.Id, uint32(property), &value)
	return
}

func (shader *Shader) GetCompileStatus() (success bool) {
	success = toBool(shader.Getiv(CompileStatusShaderProperty))
	return
}

func (shader *Shader) GetLogInfoLength() (length int) {
	length = int(shader.Getiv(LogInfoLengthShaderProperty))
	return
}

func (shader *Shader) GetInfoLog(maxLength int) (log string) {
	logIn, logOut := glStrOut(maxLength)
	gl.GetShaderInfoLog(shader.Id, int32(maxLength), nil, logIn())
	log = logOut()
	return
}
