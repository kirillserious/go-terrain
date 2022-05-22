// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main // import "github.com/go-gl/terrain/gl41core-cube"

import (
	_ "embed"
	"fmt"
	"go/build"
	"terrain/internal"
	gl2 "terrain/internal/gl"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var windowWidth = 800
var windowHeight = 600

func debugCb(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {

	msg := fmt.Sprintf("[GL_DEBUG] source %d gltype %d id %d severity %d length %d: %s", source, gltype, id, severity, length, message)
	if severity == gl.DEBUG_SEVERITY_HIGH {
		panic(msg)
	}
	fmt.Println(msg)
}

func main() {
	window, terminate := internal.NewWindow(windowWidth, windowHeight)
	defer terminate()

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(debugCb, unsafe.Pointer(nil))

	program, err := gl2.NewProgram(map[gl2.ShaderKind]string{
		gl2.VertexShaderKind:   vertexShader,
		gl2.FragmentShaderKind: fragmentShader,
	})
	if err != nil {
		panic(err)
	}

	model := mgl32.Ident4()

	// Load the texture
	texture := gl2.NewTextureFromFile("untitled.png")
	// Configure the vertex data
	vertexBuf := gl2.NewArrayBuffer(cubeVertices, 5, gl2.StaticDrawBufferUsage)

	angle := 0.0

	coord := internal.NewCoord()

	heightMap := internal.GenerateHeightMap(1000, 1000)
	terrain := internal.NewTerrain2(heightMap, texture)
	window.Render(func() {
		coord.Draw(window.Camera().ViewMatrix(), window.Projection())
		terrain.Draw(window.Camera(), window.Projection())

		angle += window.DeltaTime()
		model = mgl32.Translate3D(-5, 1, 0).Mul4(mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0}))

		//camera = mainCamera.ViewMatrix()
		//cursorRayPos, cursorRayDir := mgl32.UnProject()

		//simpleModel := mgl32.Translate3D(5, 1, 0)
		// Render

		//texture.Bind()
		program.MustDraw(gl2.TrianglesDrawMode,
			vertexBuf,
			texture,
			map[string]interface{}{
				"vert":         gl2.BufferBind{Size: 3, Offset: 0},
				"vertTexCoord": gl2.BufferBind{Size: 2, Offset: 3},
				"model":        model,
				"camera":       window.Camera().ViewMatrix(),
				"projection":   window.Projection(),
				"tex":          0,
			})

		/*
			programs.ColorOnly().MustDraw(gl2.TrianglesDrawMode, vertexBuf, map[string]interface{}{
				"vert":       gl2.BufferBind{Size: 3, Offset: 0},
				"model":      simpleModel,
				"camera":     camera,
				"projection": projection,
				"Color":      mgl32.Vec4{1.0, 0.0, 0.0, 0.0},
			})*/

	})
}

var vertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
`

var fragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord);
}
`

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
/*func init() {
	dir, err := importPathToDir("github.com/go-gl/terrain/gl41core-cube")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}
*/
// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
