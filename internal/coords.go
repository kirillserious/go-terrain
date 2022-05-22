package internal

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Coord struct {
	xLine, yLine, zLine *Line
}

func NewCoord() (coord *Coord) {
	coord = new(Coord)
	coord.xLine = NewLine(mgl.Vec3{0, 0, 0}, mgl.Vec3{1, 0, 0})
	coord.xLine.SetColor(mgl.Vec4{1, 0, 0, 1})
	coord.yLine = NewLine(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})
	coord.yLine.SetColor(mgl.Vec4{0, 1, 0, 1})
	coord.zLine = NewLine(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 0, 1})
	coord.zLine.SetColor(mgl.Vec4{0, 0, 1, 1})
	return
}

func (coord *Coord) Draw(view, projection mgl.Mat4) {
	coord.xLine.Draw(view, projection)
	coord.yLine.Draw(view, projection)
	coord.zLine.Draw(view, projection)
}
