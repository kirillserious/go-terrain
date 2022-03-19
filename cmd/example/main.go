// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main // import "github.com/go-gl/example/gl41core-cube"

import (
	_ "embed"
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"runtime"
	gl2 "terrain/cmd/example/gl"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	log "github.com/sirupsen/logrus"
)

var (
	//go:embed simple.vertex.glsl
	SimpleVertexShader string

	//go:embed simple.fragment.glsl
	SimpleFragmentShader string
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var mainCamera = Camera{
	Position: mgl32.Vec3{0, -0, 0}, WorldUp: mgl32.Vec3{0, 1, 0},
	Yaw: 0, Pitch: 0, Zoom: 0,
	MovementSpeed: 10, RotationSpeed: 100,
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := CreateWindow()
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(keyCallback())

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	log.Infof("OpenGL version: %s", gl2.GetVersion())

	gl.Viewport(0, 0, windowWidth, windowHeight)
	// Configure the vertex and fragment shaders
	program, err := newProgram(SimpleVertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0)
	camera := mainCamera.ViewMatrix()
	model := mgl32.Ident4()

	program.SetArgument("projection", projection)
	program.SetArgument("camera", camera)
	program.SetArgument("model", model)
	program.SetArgument("tex", 0)

	simpleProgram, err := newProgram(SimpleVertexShader, SimpleFragmentShader)
	if err != nil {
		panic(err)
	}
	simpleProgram.SetArgument("projection", projection)
	simpleProgram.SetArgument("camera", camera)
	simpleProgram.SetArgument("model", model)

	// Load the texture
	texture, err := newTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	vertexBuf := gl2.NewArrayBuffer(cubeVertices, 5, gl2.StaticDrawBufferUsage)

	program.SetArgument("vert",
		gl2.MyBufferArg{vertexBuf, 3, 0},
	)
	program.SetArgument("vertTexCoord",
		gl2.MyBufferArg{vertexBuf, 2, 3},
	)
	simpleProgram.SetArgument("vert",
		gl2.MyBufferArg{vertexBuf, 3, 0},
	)
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 0.9, 0.8, 0.7)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		doMovement(elapsed)

		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		camera = mainCamera.ViewMatrix()
		program.SetArgument("camera", camera)
		simpleProgram.SetArgument("camera", camera)

		simpleModel := mgl32.Translate3D(-5, 0, 0).Mul4(model)
		// Render
		program.SetArgument("model", model)
		simpleProgram.SetArgument("model", simpleModel)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		program.DrawArray(vertexBuf)

		simpleProgram.DrawArray(vertexBuf)
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (*gl2.Program, error) {
	vertexShader, err := gl2.NewShader(gl2.VertexShaderKind, vertexShaderSource)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := gl2.NewShader(gl2.FragmentShaderKind, fragmentShaderSource)
	if err != nil {
		return nil, err
	}

	program, err := gl2.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}
	vertexShader.Delete()
	fragmentShader.Delete()

	return program, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

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
func init() {
	dir, err := importPathToDir("github.com/go-gl/example/gl41core-cube")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

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

var PressedKeys = map[glfw.Key]struct{}{}

var movementMap = map[glfw.Key]Direction{
	glfw.KeyW: ForwardDirection,
	glfw.KeyS: BackwardDirection,
	glfw.KeyA: LeftDirection,
	glfw.KeyD: RightDirection,
	glfw.KeyR: UpRotation,
	glfw.KeyF: DownRotation,
	glfw.KeyQ: LeftRotation,
	glfw.KeyE: RightRotation,
}

func doMovement(delta float64) {
	for key, direction := range movementMap {
		if _, ok := PressedKeys[key]; ok {
			mainCamera.Move(delta, direction)
		}
	}
}

func keyCallback() KeyCallback {
	var (
		fullScreen                                bool
		prevXPos, prevYPos, prevHeight, prevWidth int
	)
	return func(window *Window, key glfw.Key, action glfw.Action) {
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
			return
		}
		if key == glfw.KeyF11 && action == glfw.Press {
			if fullScreen {
				window.UnFullScreen(prevXPos, prevYPos, prevWidth, prevHeight)
			} else {
				prevXPos, prevYPos = window.GetPos()
				prevWidth, prevHeight = window.GetSize()
				window.FullScreen(GetPrimaryMonitor())
			}
			fullScreen = !fullScreen
		}
		switch action {
		case glfw.Press:
			PressedKeys[key] = struct{}{}
		case glfw.Release:
			delete(PressedKeys, key)
		}
	}
}
