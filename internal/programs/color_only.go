package programs

import (
	_ "embed"
	gl2 "terrain/internal/gl"
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

func LasyProgram(m map[gl2.ShaderKind]string) func() *gl2.Program {
	return func() func() *gl2.Program {
		var program *gl2.Program
		return func() *gl2.Program {
			if program == nil {
				var err error
				program, err = gl2.NewProgram(m)
				if err != nil {
					panic(err)
				}
			}
			return program
		}
	}()
}

var ColorOnly = LasyProgram(map[gl2.ShaderKind]string{
	gl2.VertexShaderKind:   SimpleVertexSource,
	gl2.FragmentShaderKind: ColorOnlyFragmentSource,
})

var Phong = LasyProgram(map[gl2.ShaderKind]string{
	gl2.VertexShaderKind:   WithNormalVertexSource,
	gl2.FragmentShaderKind: PhongFragmentSource,
})
