package main

import (
	"github.com/go-gl/mathgl/mgl32"
	gl2 "terrain/internal/gl"
	"terrain/internal/programs"
)

func AddNormals(in []float32) (out []float32) {
	if len(in)%3 != 0 {
		panic("incorrect in length")
	}
	out = make([]float32, len(in)*2)
	ParallelFor(0, len(in), 9, func(i int) {
		p1 := mgl32.Vec3{in[i+0], in[i+1], in[i+2]}
		p2 := mgl32.Vec3{in[i+3], in[i+4], in[i+5]}
		p3 := mgl32.Vec3{in[i+6], in[i+7], in[i+8]}
		u, v := p2.Sub(p1), p3.Sub(p1)
		normal := u.Cross(v)

		j := i * 2
		out[j], out[j+1], out[j+2] = in[i], in[i+1], in[i+2]
		out[j+3], out[j+4], out[j+5] = normal[0], normal[1], normal[2]
		out[j+6], out[j+7], out[j+8] = in[i+3], in[i+4], in[i+5]
		out[j+9], out[j+10], out[j+11] = normal[0], normal[1], normal[2]
		out[j+12], out[j+13], out[j+14] = in[i+6], in[i+7], in[i+8]
		out[j+15], out[j+16], out[j+17] = normal[0], normal[1], normal[2]
	})
	return
}

var CubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	// Top
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,

	// Left
	-1.0, -1.0, 1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	// Right
	1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
}

var LasyCubeBuffer = func() func() *gl2.ArrayBuffer {
	withNormals := AddNormals(CubeVertices)
	var buffer *gl2.ArrayBuffer
	return func() *gl2.ArrayBuffer {
		if buffer != nil {
			return buffer
		}
		buffer = gl2.NewArrayBuffer(withNormals, 6, gl2.StaticDrawBufferUsage)
		return buffer
	}
}()

type Cube struct{}

func (cube *Cube) Draw(camera Camera, projection mgl32.Mat4) {
	programs.Phong().MustDraw(gl2.TrianglesDrawMode, LasyCubeBuffer(), map[string]interface{}{
		"vert":       gl2.BufferBind{Size: 3, Offset: 0},
		"normal":     gl2.BufferBind{Size: 3, Offset: 3},
		"model":      mgl32.Ident4(),
		"camera":     camera.ViewMatrix(),
		"projection": projection,
		"viewPos":    camera.Position,
		"Color":      mgl32.Vec4{1.0, 0.5, 0.31, 1.0},
	})
}
