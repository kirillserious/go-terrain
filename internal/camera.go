package internal

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	Position, WorldUp            mgl.Vec3
	Yaw, Pitch, Zoom             float32
	MovementSpeed, RotationSpeed float32
}

func cos(arg float32) float32 {
	return float32(math.Cos(float64(arg)))
}

func sin(arg float32) float32 {
	return float32(math.Sin(float64(arg)))
}

func (camera *Camera) PositionalVectors() (right, front, up mgl.Vec3) {
	yaw, pitch := mgl.DegToRad(camera.Yaw), mgl.DegToRad(camera.Pitch)
	front[0] = cos(yaw) * cos(pitch)
	front[1] = sin(pitch)
	front[2] = sin(yaw) * cos(pitch)
	front = front.Normalize()
	right = front.Cross(camera.WorldUp).Normalize()
	up = right.Cross(front).Normalize()
	return
}

func (camera *Camera) ViewMatrix() mgl.Mat4 {
	_, front, up := camera.PositionalVectors()
	return mgl.LookAtV(camera.Position, camera.Position.Add(front), up)
}

type Direction int

const (
	ForwardDirection Direction = iota
	BackwardDirection
	UpDirection
	DownDirection
	LeftDirection
	RightDirection
	UpRotation
	DownRotation
	LeftRotation
	RightRotation
)

func (camera *Camera) Move(delta float64, directions ...Direction) {
	right, front, up := camera.PositionalVectors()
	positionVelocity, rotationVelocity := camera.MovementSpeed*float32(delta), camera.RotationSpeed*float32(delta)

	for _, direction := range directions {
		switch direction {
		case ForwardDirection:
			camera.Position = camera.Position.Add(front.Mul(positionVelocity))
		case BackwardDirection:
			camera.Position = camera.Position.Sub(front.Mul(positionVelocity))
		case UpDirection:
			camera.Position = camera.Position.Add(up.Mul(positionVelocity))
		case DownDirection:
			camera.Position = camera.Position.Sub(up.Mul(positionVelocity))
		case LeftDirection:
			camera.Position = camera.Position.Sub(right.Mul(positionVelocity))
		case RightDirection:
			camera.Position = camera.Position.Add(right.Mul(positionVelocity))
		case UpRotation:
			camera.Pitch += rotationVelocity
		case DownRotation:
			camera.Pitch -= rotationVelocity
		case LeftRotation:
			camera.Yaw -= rotationVelocity
			if camera.Yaw < -89 {
				camera.Yaw = -89
			}
		case RightRotation:
			camera.Yaw += rotationVelocity
			if camera.Yaw > 89 {
				camera.Yaw = 89
			}
		}
	}
}
