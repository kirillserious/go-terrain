package gl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Program struct {
	Id uint32

	args map[string]ProgramArgumentLocation
	vaos map[*ArrayBuffer]*vao
}

// ProgramArgumentLocation is location of the program argument
type ProgramArgumentLocation int32

func NewProgram(shaders ...*Shader) (program *Program, err error) {
	program = CreateProgram()
	// Link shaders
	for _, shader := range shaders {
		program.AttachShader(shader)
	}
	program.Link()
	// Check if everything compiles correct
	if success := program.GetLinkStatus(); !success {
		logLength := program.GetInfoLogLength()
		infoLog := program.GetInfoLog(logLength)
		program, err = nil, fmt.Errorf("failed to link program: %v", infoLog)
	}
	return
}

func (program *Program) SetArgumentByLocation(loc ProgramArgumentLocation, value interface{}) {
	asWas := program.use()
	defer asWas()

	switch cast := value.(type) {
	case mgl32.Mat4:
		gl.UniformMatrix4fv(int32(loc), 1, false, &cast[0])
	case int:
		gl.Uniform1i(int32(loc), int32(cast))
	case MyBufferArg:
		vao, ok := program.vaos[cast.ArrayBuffer]
		if !ok {
			program.vaos[cast.ArrayBuffer] = NewVAO(cast.ArrayBuffer)
			vao = program.vaos[cast.ArrayBuffer]
		}
		vao.LinkWithLocation(loc, cast.Size, cast.Offset)
	}
}

type MyBufferArg struct {
	*ArrayBuffer
	Size, Offset int
}

func (program *Program) SetArgument(name string, value interface{}) {
	argLoc, ok := program.args[name]
	if !ok {
		argLoc = program.UniformArgLocation(name)
		program.args[name] = argLoc
	}

	var loc ProgramArgumentLocation
	switch value.(type) {
	case int, mgl32.Mat4:
		loc = program.UniformArgLocation(name)
	case MyBufferArg:
		loc = program.InputArgLocation(name)
	}
	program.SetArgumentByLocation(loc, value)
}

type ArrayArgument struct {
	Name         string
	Size, Offset int
}

func (program *Program) DrawArray(buf *ArrayBuffer) (err error) {
	asWas := program.use()
	defer asWas()

	vao, ok := program.vaos[buf]
	if !ok {
		err = fmt.Errorf("array is unset as input argument")
		return
	}
	vao.Bind()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(buf.Buffer)/buf.Stride))
	UnbindVAO()
	return
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

func CreateProgram() (program *Program) {
	program = new(Program)
	program.Id = gl.CreateProgram()
	program.args = make(map[string]ProgramArgumentLocation)
	program.vaos = make(map[*ArrayBuffer]*vao)

	outIn := glStrIn("outputColor")
	gl.BindFragDataLocation(program.Id, 0, outIn())
	return
}

// UniformArgLocation returns the location of a uniform program argument
//
// An argument is uniform if every video-card processor have the same value of the argument
func (program *Program) UniformArgLocation(name string) (loc ProgramArgumentLocation) {
	nameIn := glStrIn(name)
	loc = ProgramArgumentLocation(gl.GetUniformLocation(program.Id, nameIn()))
	return
}

// InputArgLocation returns the location of an input program argument
//
// An argument is input if every video-card processor have its own argument (for example,
// with vertex array binding).
func (program *Program) InputArgLocation(name string) (loc ProgramArgumentLocation) {
	nameIn := glStrIn(name)
	loc = ProgramArgumentLocation(gl.GetAttribLocation(program.Id, nameIn()))
	return
}

func (program *Program) AttachShader(shader *Shader) {
	gl.AttachShader(program.Id, shader.Id)
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

func (program *Program) Link() {
	gl.LinkProgram(program.Id)
}
