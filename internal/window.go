package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Monitor struct {
	*glfw.Monitor
}

func GetPrimaryMonitor() (m *Monitor) {
	m = new(Monitor)
	m.Monitor = glfw.GetPrimaryMonitor()
	return
}

// Window embeds glfw.Window
type Window struct {
	*glfw.Window
}

func CreateWindow(width, height int) (w *Window, err error) {
	w = new(Window)
	w.Window, err = glfw.CreateWindow(width, height, "Cube", nil, nil)
	return
}

type KeyCallback func(*Window, glfw.Key, glfw.Action)

func (w *Window) SetKeyCallback(callback KeyCallback) {
	gCallback := func(window *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		callback(w, key, action)
	}
	w.Window.SetKeyCallback(gCallback)
}

func (w *Window) SetMonitor(monitor *Monitor, xPos, yPos, width, height, refreshRate int) {
	w.Window.SetMonitor(monitor.Monitor, xPos, yPos, width, height, refreshRate)
}

func (w *Window) FullScreen(monitor *Monitor) {
	mode := monitor.GetVideoMode()
	w.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
}

func (w *Window) UnFullScreen(xPos, yPos, width, height int) {
	w.Window.SetMonitor(nil, xPos, yPos, width, height, 0)
}
