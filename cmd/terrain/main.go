// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main // import "github.com/go-gl/terrain/gl41core-cube"

import (
	_ "embed"
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"runtime"
	"terrain/internal"
	gl2 "terrain/internal/gl"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	log "github.com/sirupsen/logrus"
)

var windowWidth = 800
var windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var mainCamera = internal.Camera{
	Position: mgl32.Vec3{10, 10, 10}, WorldUp: mgl32.Vec3{0, 1, 0},
	Yaw: 0, Pitch: 0, Zoom: 0,
	MovementSpeed: 30, RotationSpeed: 100,
}

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
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := internal.CreateWindow(windowWidth, windowHeight)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(keyCallback())
	window.SetCursorPosCallback(cursorCallback)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	log.Infof("OpenGL version: %s", gl2.GetVersion())

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(debugCb, unsafe.Pointer(nil))

	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))

	/*
		program, err := gl2.NewProgram(map[gl2.ShaderKind]string{
			gl2.VertexShaderKind:   vertexShader,
			gl2.FragmentShaderKind: fragmentShader,
		})
		if err != nil {
			panic(err)
		}
	*/
	near, far := float32(0.1), float32(100)
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), near, far)
	camera := mainCamera.ViewMatrix()
	//model := mgl32.Ident4()
	/*
		// Load the texture
		texture, err := newTexture("square.png")
		if err != nil {
			log.Fatalln(err)
		}
	*/
	// Configure the vertex data
	//vertexBuf := gl2.NewArrayBuffer(cubeVertices, 5, gl2.StaticDrawBufferUsage)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//gl.ClearColor(1.0, 0.9, 0.8, 0.7)

	angle := 0.0
	previousTime := glfw.GetTime()

	coord := internal.NewCoord()

	heightMap := internal.GenerateHeightMap(1000, 1000)
	terrain := internal.NewTerrain(heightMap)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		windowWidth, windowHeight = window.GetSize()
		projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/float32(windowHeight), near, far)
		gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		doMovement(elapsed)

		angle += elapsed
		//model = mgl32.Translate3D(-5, 1, 0).Mul4(mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0}))

		camera = mainCamera.ViewMatrix()
		//cursorRayPos, cursorRayDir := mgl32.UnProject()

		//simpleModel := mgl32.Translate3D(5, 1, 0)
		// Render
		/*
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, texture)
			program.MustDraw(gl2.TrianglesDrawMode, vertexBuf, map[string]interface{}{
				"vert":         gl2.BufferBind{Size: 3, Offset: 0},
				"vertTexCoord": gl2.BufferBind{Size: 2, Offset: 3},
				"model":        model,
				"camera":       camera,
				"projection":   projection,
				"tex":          0,
			})
		*/
		/*
			programs.ColorOnly().MustDraw(gl2.TrianglesDrawMode, vertexBuf, map[string]interface{}{
				"vert":       gl2.BufferBind{Size: 3, Offset: 0},
				"model":      simpleModel,
				"camera":     camera,
				"projection": projection,
				"Color":      mgl32.Vec4{1.0, 0.0, 0.0, 0.0},
			})*/

		coord.Draw(camera, projection)
		terrain.Draw(mainCamera, projection)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
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

var movementMap = map[glfw.Key]internal.Direction{
	glfw.KeyW: internal.ForwardDirection,
	glfw.KeyS: internal.BackwardDirection,
	glfw.KeyA: internal.LeftDirection,
	glfw.KeyD: internal.RightDirection,
	glfw.KeyR: internal.UpRotation,
	glfw.KeyF: internal.DownRotation,
	glfw.KeyQ: internal.LeftRotation,
	glfw.KeyE: internal.RightRotation,
}

func doMovement(delta float64) {
	for key, direction := range movementMap {
		if _, ok := PressedKeys[key]; ok {
			mainCamera.Move(delta, direction)
		}
	}
}

var PressedKeys = map[glfw.Key]struct{}{}

func keyCallback() internal.KeyCallback {
	var (
		fullScreen                                bool
		prevXPos, prevYPos, prevHeight, prevWidth int
	)
	return func(window *internal.Window, key glfw.Key, action glfw.Action) {
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
				window.FullScreen(internal.GetPrimaryMonitor())
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

var CursorX, CursorY float64

func cursorCallback(_ *glfw.Window, xPos float64, yPos float64) {
	CursorX, CursorY = xPos, yPos
}
