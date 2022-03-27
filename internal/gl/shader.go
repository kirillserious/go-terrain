package gl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"regexp"
	"strings"
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
	shader = new(Shader)
	shader.Id, shader.Kind, shader.Source = generateShaderId(kind), kind, source
	shader.setSources(source)
	shader.compile()

	if success := shader.GetCompileStatus(); !success {
		logLength := shader.GetLogInfoLength()
		log := shader.GetInfoLog(logLength)
		shader = nil
		err = fmt.Errorf("failed to compile shader %v: %v", source, log)
	}
	return
}

var (
	uniformArgRegexp = regexp.MustCompile("^uniform .*;")
	inArgRegexp      = regexp.MustCompile("^in .*;")
	outArgRegexp     = regexp.MustCompile("^out .*;")
)

func (shader *Shader) Args() (uniform []string, in []string, out []string) {
	lines := strings.Split(shader.Source, "\n")

	tool := func(line string, slice *[]string, re *regexp.Regexp) {
		if found := re.FindString(line); found != "" {
			split := strings.Split(found, " ")
			last := split[len(split)-1]
			*slice = append(*slice, last[:len(last)-1])
		}
	}
	for _, line := range lines {
		tool(line, &uniform, uniformArgRegexp)
		tool(line, &in, inArgRegexp)
		tool(line, &out, outArgRegexp)
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

func generateShaderId(kind ShaderKind) (id uint32) {
	id = gl.CreateShader(uint32(kind))
	return
}

func (shader *Shader) compile() {
	gl.CompileShader(shader.Id)
}

func (shader *Shader) setSources(sources ...string) {
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
