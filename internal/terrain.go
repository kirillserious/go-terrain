package internal

import (
	"github.com/go-gl/mathgl/mgl32"
	gl2 "terrain/internal/gl"
	"terrain/internal/programs"
)

type Terrain struct {
	buf *gl2.ArrayBuffer
}

func NewTerrain(heights *HeightMap) (terrain *Terrain) {
	terrain = new(Terrain)
	iMax, jMax := len(heights.Heights)/heights.Stride, heights.Stride
	vertices := make([]float32, (iMax-1)*(jMax-1)*6*3)
	for i := 0; i < iMax-1; i++ {
		for j := 0; j < jMax-1; j++ {
			h1 := heights.Heights[i*jMax+j]
			h2 := heights.Heights[i*jMax+j+1]
			h3 := heights.Heights[(i+1)*jMax+j+1]
			h4 := heights.Heights[(i+1)*jMax+j]

			l := (i*(jMax-1) + j) * 6 * 3
			vertices[l], vertices[l+1], vertices[l+2] = float32(i), h1, float32(j)
			vertices[l+3], vertices[l+4], vertices[l+5] = float32(i), h2, float32(j+1)
			vertices[l+6], vertices[l+7], vertices[l+8] = float32(i+1), h3, float32(j+1)
			vertices[l+9], vertices[l+10], vertices[l+11] = float32(i+1), h3, float32(j+1)
			vertices[l+12], vertices[l+13], vertices[l+14] = float32(i+1), h4, float32(j)
			vertices[l+15], vertices[l+16], vertices[l+17] = float32(i), h1, float32(j)
		}
	}

	terrain.buf = gl2.NewArrayBuffer(AddNormals(vertices), 6, gl2.StaticDrawBufferUsage)
	return
}

func (terrain *Terrain) Draw(camera *Camera, projection mgl32.Mat4) {
	programs.Phong().MustDraw(gl2.TrianglesDrawMode,
		terrain.buf,
		nil,
		map[string]interface{}{
			"vert":       gl2.BufferBind{Size: 3, Offset: 0},
			"normal":     gl2.BufferBind{Size: 3, Offset: 3},
			"model":      mgl32.Scale3D(0.025, 0.025, 0.025),
			"camera":     camera.ViewMatrix(),
			"projection": projection,
			"viewPos":    camera.Position,
			"Color":      mgl32.Vec4{0.4, 0.6, 0.0, 1.0},
		},
	)
}
