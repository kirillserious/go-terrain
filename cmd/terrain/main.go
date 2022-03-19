package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize GLFW")
	}
	defer glfw.Terminate()
	log.Info("GLFW successfully initialized")

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(800, 600, "Terrain", nil, nil)
	if err != nil {
		log.WithError(err).Fatal("Failed to create window")
	}
	window.MakeContextCurrent()
	log.Info("Window successfully created")

	if err := gl.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize OpenGL")
	}
	log.Info("OpenGL successfully initialized")

	window.SetKeyCallback(keyCallback())

	triangleProgram, _, triangleVBO := InitTriangle()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.2, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindBuffer(gl.ARRAY_BUFFER, triangleVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
		gl.UseProgram(triangleProgram.Id)
		vertAttrib := uint32(gl.GetAttribLocation(triangleProgram.Id, gl.Str("position\x00")))
		gl.EnableVertexAttribArray(vertAttrib)
		gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 3*4, 0)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func keyCallback() glfw.KeyCallback {
	var (
		fullScreen                                bool
		prevXPos, prevYPos, prevHeight, prevWidth int
	)
	return func(window *glfw.Window, key glfw.Key, _ int,
		action glfw.Action, _ glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
			return
		}

		if key == glfw.KeyF11 && action == glfw.Press {
			if fullScreen {
				window.SetMonitor(nil, prevXPos, prevYPos, prevWidth, prevHeight, 0)
			} else {
				prevXPos, prevYPos = window.GetPos()
				prevWidth, prevHeight = window.GetSize()
				monitor := glfw.GetPrimaryMonitor()
				mode := monitor.GetVideoMode()
				window.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
			}
			fullScreen = !fullScreen
		}
	}
}
