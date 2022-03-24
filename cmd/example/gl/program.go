package gl

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program struct {
	Id uint32

	InArgs, UniformArgs map[string]ProgramArgumentLocation
	VAO                 map[*ArrayBuffer]*vao
}

const OutputParameterColor = "outputColor"

// ProgramArgumentLocation is location of the program argument
type ProgramArgumentLocation int32

type DrawMode int

const (
	TrianglesDrawMode DrawMode = gl.TRIANGLES
	LinesDrawMod               = gl.LINES
)

func NewProgram(sources map[ShaderKind]string) (program *Program, err error) {
	// TODO: Check does it needed
	if _, ok := sources[VertexShaderKind]; !ok {
		err = fmt.Errorf("no vertex shader")
		return
	}
	if _, ok := sources[FragmentShaderKind]; !ok {
		err = fmt.Errorf("no fragment shader")
		return
	}

	program = new(Program)
	program.Id = gl.CreateProgram()
	program.VAO = make(map[*ArrayBuffer]*vao)
	outIn := glStrIn(OutputParameterColor)
	gl.BindFragDataLocation(program.Id, 0, outIn())

	shaders := make([]*Shader, 0, len(sources))
	defer func() {
		for _, shader := range shaders {
			shader.Delete()
		}
	}()

	for kind, source := range sources {
		shader, _err := NewShader(kind, source)
		if _err != nil {
			err = _err
			return
		}
		gl.AttachShader(program.Id, shader.Id)
		shaders = append(shaders, shader)
	}

	gl.LinkProgram(program.Id)
	if success := program.GetLinkStatus(); !success {
		logLength := program.GetInfoLogLength()
		infoLog := program.GetInfoLog(logLength)
		program, err = nil, fmt.Errorf("failed to link program: %v", infoLog)
	}

	program.InArgs, program.UniformArgs = make(map[string]ProgramArgumentLocation), make(map[string]ProgramArgumentLocation)
	for _, shader := range shaders {
		uniform, in, out := shader.Args()
		for _, arg := range uniform {
			program.UniformArgs[arg] = program.UniformArgLocation(arg)
		}

		if shader.Kind == VertexShaderKind {
			for _, arg := range in {
				program.InArgs[arg] = program.InputArgLocation(arg)
			}
		}

		if shader.Kind == FragmentShaderKind {
			if len(out) != 1 {
				err = fmt.Errorf("incorrect out arguments number")
				return
			}
			if out[0] != OutputParameterColor {
				err = fmt.Errorf("incorrect output parameter color")
				return
			}
		}
	}
	return
}

type BufferBind struct {
	Size, Offset int
}

func (program *Program) MustDraw(mode DrawMode, buffer *ArrayBuffer, args map[string]interface{}) {
	asWas := program.use()
	defer asWas()

	// TODO: Add check for args
	if len(args) != len(program.InArgs)+len(program.UniformArgs) {
		panic("incorrect args number")
	}

	vao, ok := program.VAO[buffer]
	if !ok {
		program.VAO[buffer] = NewVAO(buffer)
		vao = program.VAO[buffer]
		for name, value := range args {
			if cast, ok := value.(BufferBind); ok {
				vao.LinkWithLocation(program.InArgs[name], cast.Size, cast.Offset)
			}
		}
		program.VAO[buffer] = vao
	}
	for name, value := range args {
		loc := program.UniformArgs[name]
		switch cast := value.(type) {
		case mgl32.Mat4:
			gl.UniformMatrix4fv(int32(loc), 1, false, &cast[0])
		case mgl32.Vec4:
			gl.Uniform4fv(int32(loc), 1, &cast[0])
		case int:
			gl.Uniform1i(int32(loc), int32(cast))
		}
	}

	vao.Bind()
	gl.DrawArrays(uint32(mode), 0, int32(len(buffer.Buffer)/buffer.Stride))
	UnbindVAO()
}

func (program *Program) use() (asWas func()) {
	currentId := GetCurrentProgramId()
	if currentId != int(program.Id) {
		gl.UseProgram(program.Id)
		asWas = func() { gl.UseProgram(uint32(currentId)) }
		return
	}
	asWas = func() {}
	return
}

// GL Functions

func (program *Program) UniformArgLocation(name string) (loc ProgramArgumentLocation) {
	nameIn := glStrIn(name)
	loc = ProgramArgumentLocation(gl.GetUniformLocation(program.Id, nameIn()))
	return
}

func (program *Program) InputArgLocation(name string) (loc ProgramArgumentLocation) {
	nameIn := glStrIn(name)
	loc = ProgramArgumentLocation(gl.GetAttribLocation(program.Id, nameIn()))
	return
}

type ProgramProperty uint32

const (
	LinkStatusProgramProperty    ProgramProperty = gl.LINK_STATUS
	InfoLogLengthProgramProperty                 = gl.INFO_LOG_LENGTH
)

func (program *Program) Getiv(param ProgramProperty) (value int32) {
	gl.GetProgramiv(program.Id, uint32(param), &value)
	return
}

func (program *Program) GetInfoLogLength() (length int) {
	length = int(program.Getiv(InfoLogLengthProgramProperty))
	return
}

func (program *Program) GetInfoLog(maxLength int) (log string) {
	logIn, logOut := glStrOut(maxLength)
	gl.GetProgramInfoLog(program.Id, int32(maxLength), nil, logIn())
	log = logOut()
	return
}

func (program *Program) GetLinkStatus() (success bool) {
	success = toBool(program.Getiv(LinkStatusProgramProperty))
	return
}
