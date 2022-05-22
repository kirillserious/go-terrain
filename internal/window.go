package internal

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	log "github.com/sirupsen/logrus"
	mygl "terrain/internal/gl"
)

type Window struct {
	GLWindow  *mygl.Window
	Near, Far float32

	camera     *Camera
	projection mgl32.Mat4

	timer          *Timer
	current, delta float64
}

func NewWindow(width, height int) (window *Window, terminate func()) {
	terminate = GLFWMustInit()
	var err error

	window = new(Window)
	window.GLWindow, err = mygl.CreateWindow(width, height)
	if err != nil {
		log.WithError(err).Panic("failed to create window")
	}
	window.GLWindow.MakeContextCurrent()
	window.GLWindow.SetKeyCallback(keyCallback())
	window.GLWindow.SetCursorPosCallback(cursorCallback)

	if err := gl.Init(); err != nil {
		log.WithError(err).Panic("failed to init GL")
		panic(err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//gl.ClearColor(1.0, 0.9, 0.8, 0.7)

	window.camera = new(Camera)
	window.camera.Position = mgl32.Vec3{10, 10, 10}
	window.camera.WorldUp = mgl32.Vec3{0, 1, 0}
	window.camera.Yaw, window.camera.Pitch, window.camera.Zoom = 0, 0, 0
	window.camera.MovementSpeed, window.camera.RotationSpeed = 30, 100

	window.Near, window.Far = 0.1, 100
	window.projection = projectionMatrix(width, height, window.Near, window.Far)

	window.timer = NewTimer()
	return
}

func (window *Window) IsOpen() bool {
	return !window.GLWindow.ShouldClose()
}

func (window *Window) GetSize() (int, int) {
	return window.GLWindow.GetSize()
}

func (window *Window) Camera() *Camera {
	return window.camera
}

func projectionMatrix(windowWidth, windowHeight int, near, far float32) mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(45.0),
		float32(windowWidth)/float32(windowHeight),
		near,
		far,
	)
}

func (window *Window) Projection() mgl32.Mat4 {
	return window.projection
}

func (window *Window) CurrentTime() float64 {
	return window.current
}

func (window *Window) DeltaTime() float64 {
	return window.delta
}

func (window *Window) Render(drawer func()) {
	for !window.GLWindow.ShouldClose() {
		window.GLWindow.SwapBuffers()
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		width, height := window.GetSize()
		gl.Viewport(0, 0, int32(width), int32(height))
		window.current, window.delta = window.timer.Update()

		window.moveCamera(window.delta)
		window.projection = projectionMatrix(width, height, window.Near, window.Far)

		drawer()
	}
}

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

func (window *Window) moveCamera(delta float64) {
	for key, direction := range movementMap {
		if _, ok := PressedKeys[key]; ok {
			window.camera.Move(delta, direction)
		}
	}
}

var PressedKeys = map[glfw.Key]struct{}{}

func keyCallback() mygl.KeyCallback {
	var (
		fullScreen                                bool
		prevXPos, prevYPos, prevHeight, prevWidth int
	)
	return func(window *mygl.Window, key glfw.Key, action glfw.Action) {
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
				window.FullScreen(mygl.GetPrimaryMonitor())
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
