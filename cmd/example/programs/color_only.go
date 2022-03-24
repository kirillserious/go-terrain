package programs

import (
	_ "embed"
	"terrain/cmd/example/gl"
)

var (
	//go:embed with_normal.vertex.glsl
	WithNormalVertexSource string
	//go:embed simple.vertex.glsl
	SimpleVertexSource string
	//go:embed color_only.fragment.glsl
	ColorOnlyFragmentSource string
	//go:embed phong.fragment.glsl
	PhongFragmentSource string
)

func LasyProgram(m map[gl.ShaderKind]string) func() *gl.Program {
	return func() func() *gl.Program {
		var program *gl.Program
		return func() *gl.Program {
			if program == nil {
				var err error
				program, err = gl.NewProgram(m)
				if err != nil {
					panic(err)
				}
			}
			return program
		}
	}()
}

var ColorOnly = LasyProgram(map[gl.ShaderKind]string{
	gl.VertexShaderKind:   SimpleVertexSource,
	gl.FragmentShaderKind: ColorOnlyFragmentSource,
})

var Phong = LasyProgram(map[gl.ShaderKind]string{
	gl.VertexShaderKind:   WithNormalVertexSource,
	gl.FragmentShaderKind: PhongFragmentSource,
})
