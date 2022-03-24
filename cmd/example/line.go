package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	gl2 "terrain/cmd/example/gl"
	"terrain/cmd/example/programs"
)

type Line struct {
	from, to mgl.Vec3
	model    mgl.Mat4

	color mgl.Vec4
}

var LineVertices = []float32{
	0, 0, 0,
	0, 0, 1,
}

var LineVertexBuffer = func() func() *gl2.ArrayBuffer {
	var buffer *gl2.ArrayBuffer
	return func() *gl2.ArrayBuffer {
		if buffer != nil {
			return buffer
		}
		buffer = gl2.NewArrayBuffer(LineVertices, 3, gl2.StaticDrawBufferUsage)
		return buffer
	}
}()

func NewLine(from, to mgl.Vec3) (line *Line) {
	line = new(Line)
	line.from, line.to, line.color = from, to, mgl.Vec4{0, 0, 0, 1}
	// NB: American way to numerate matrix (vertically), so we need to transpose
	line.model = mgl.Mat4{
		0, 0, to[0] - from[0], from[0],
		0, 0, to[1] - from[1], from[1],
		0, 0, to[2] - from[2], from[2],
		0, 0, 0, 1,
	}.Transpose()
	return
}

func (line *Line) SetColor(color mgl.Vec4) {
	line.color = color
}

func (line *Line) Draw(view, projection mgl.Mat4) {
	programs.ColorOnly().MustDraw(gl2.LinesDrawMod, LineVertexBuffer(), map[string]interface{}{
		"vert":       gl2.BufferBind{Size: 3, Offset: 0},
		"model":      line.model,
		"camera":     view,
		"projection": projection,
		"Color":      line.color,
	})
}
