package internal

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func GLFWInit() (terminate func(), err error) {
	if err = glfw.Init(); err != nil {
		terminate = func() {}
		return
	}
	terminate = glfw.Terminate

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	return
}

func GLFWMustInit() (terminate func()) {
	terminate, err := GLFWInit()
	if err != nil {
		log.WithError(err).Panicf("failed to init glfw")
	}
	return
}
