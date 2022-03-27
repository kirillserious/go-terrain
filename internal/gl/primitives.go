package gl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

// Bool
func toBool(glBool int32) bool {
	return glBool == gl.TRUE
}

// Float

const sizeOfFloat = 4

// Strings

const nullTerminator string = "\x00"

func glStrIn(s string) (in func() *uint8) {
	str := s + nullTerminator
	in = func() *uint8 {
		return gl.Str(str)
	}
	return
}

func glStrOut(length int) (in func() *uint8, out func() string) {
	str := strings.Repeat(nullTerminator, length+1)
	in = func() *uint8 {
		return gl.Str(str)
	}
	out = func() string {
		return str[:length-1]
	}
	return
}

// glString converts a golang string into the gl one (c-array, that important to be free after work)
func glStrings(goStrings ...string) (glStrings **uint8, free func()) {
	for i, goString := range goStrings {
		goStrings[i] = goString + nullTerminator
	}
	glStrings, free = gl.Strs(goStrings...)
	return
}
