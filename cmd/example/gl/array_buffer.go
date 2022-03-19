package gl

import "github.com/go-gl/gl/v4.1-core/gl"

// ArrayBuffer is a VBO object with linked ArrayBuffer
type ArrayBuffer struct {
	Id     uint32
	Buffer []float32
	Stride int
}

type BufferUsage uint32

const (
	StaticDrawBufferUsage BufferUsage = gl.STATIC_DRAW
)

func NewArrayBuffer(buf []float32, stride int, usage BufferUsage) (arrBuf *ArrayBuffer) {
	arrBuf = new(ArrayBuffer)
	arrBuf.Buffer = buf
	gl.GenBuffers(1, &arrBuf.Id)
	gl.BindBuffer(gl.ARRAY_BUFFER, arrBuf.Id)
	gl.BufferData(gl.ARRAY_BUFFER, len(arrBuf.Buffer)*sizeOfFloat, gl.Ptr(arrBuf.Buffer), uint32(usage))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	arrBuf.Stride = stride
	return
}

// GL Functions

func (buf *ArrayBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, buf.Id)
}

func UnbindArrayBuffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
