package main

import (
	_ "embed"
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	//go:embed vertex_triangle.glsl
	VertexShaderSource string
	//go:embed fragment_triangle.glsl
	FragmentShaderSource string
)

var vertices = []float32{
	-0.5, -0.5, 0.0,
	0.5, -0.5, 0.0,
	0.0, 0.5, 0.0,
}

type ShaderKind uint32

const (
	VertexShaderKind   ShaderKind = gl.VERTEX_SHADER
	FragmentShaderKind ShaderKind = gl.FRAGMENT_SHADER
)

type Shader struct {
	Id     uint32
	Source string
	Kind   ShaderKind
}

func NewShader(kind ShaderKind, source string) (shader *Shader, err error) {
	shader = new(Shader)
	shader.Id = gl.CreateShader(uint32(kind))
	shader.Source = source
	shader.Kind = kind
	// Compile
	cSources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader.Id, 1, cSources, nil)
	free()
	gl.CompileShader(shader.Id)
	// Check compilation
	var success int32
	gl.GetShaderiv(shader.Id, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader.Id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader.Id, logLength, nil, gl.Str(log))
		err = fmt.Errorf("failed to compile shader: %v", log)
		shader.Destroy()
	}
	return
}

func (shader *Shader) Destroy() {
	gl.DeleteShader(shader.Id)
	shader.Id = 0
}

type ShaderProgram struct {
	Id      uint32
	Shaders []*Shader
}

func NewShaderProgram(shaders ...*Shader) (program *ShaderProgram, err error) {
	program = new(ShaderProgram)
	program.Id = gl.CreateProgram()
	program.Shaders = shaders
	// Build program
	for _, shader := range program.Shaders {
		if shader.Id == 0 {
			err = fmt.Errorf("destroyed shader")
			program.Destroy()
			return
		}
		gl.AttachShader(program.Id, shader.Id)
	}
	gl.LinkProgram(program.Id)
	// Check build
	var success int32
	gl.GetProgramiv(program.Id, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.Id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.Id, logLength, nil, gl.Str(log))
		err = fmt.Errorf("failed to link shader program: %v", log)
		program.Destroy()
	}
	return
}

func (program *ShaderProgram) Execute(argsVAO uint32) {
	gl.BindVertexArray(argsVAO)
	gl.UseProgram(program.Id)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (program *ShaderProgram) Destroy() {
	gl.DeleteProgram(program.Id)
	program.Id = 0
}

func InitTriangle() (program *ShaderProgram, vao, vbo uint32) {
	vertexShader, err := NewShader(VertexShaderKind, VertexShaderSource)
	if err != nil {
		log.WithError(err).Fatal("Failed to compile vertex shader")
	}

	fragmentShader, err := NewShader(FragmentShaderKind, FragmentShaderSource)
	if err != nil {
		log.WithError(err).Fatal("Failed to compile fragment shader")
	}

	triangleProgram, err := NewShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		log.WithError(err).Fatal("Failed to link triangle shader program")
	}
	//vertexShader.Destroy()
	//fragmentShader.Destroy()

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	//vertAttrib := uint32(gl.GetAttribLocation(triangleProgram.Id, gl.Str("vert\x00")))
	//gl.EnableVertexAttribArray(vertAttrib)
	//gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 3*4, 0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	program, vao, vbo = triangleProgram, VAO, VBO
	return
}
