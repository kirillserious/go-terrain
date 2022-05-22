package internal

import "github.com/go-gl/glfw/v3.3/glfw"

type Timer struct {
	current float64
}

func NewTimer() (timer *Timer) {
	timer = new(Timer)

	timer.current = glfw.GetTime()
	return
}

func (timer *Timer) Update() (current, delta float64) {
	current = glfw.GetTime()
	delta = current - timer.current
	timer.current = current
	return
}
