package gl

import "github.com/go-gl/gl/v4.1-core/gl"

// ArrayBuffer is a VBO object with linked ArrayBuffer
type ArrayBuffer struct {
	Id     uint32
	Buffer []float32
	Stride int
	Usage  BufferUsage
}

type BufferUsage uint32

const (
	StaticDrawBufferUsage BufferUsage = gl.STATIC_DRAW
)

func NewArrayBuffer(buf []float32, stride int, usage BufferUsage) (arrBuf *ArrayBuffer) {
	arrBuf = new(ArrayBuffer)
	arrBuf.Buffer = buf
	arrBuf.Usage = usage
	arrBuf.Stride = stride
	return
}

// GL Functions

func (buf *ArrayBuffer) Bind() {
	if buf.Id == 0 {
		gl.GenBuffers(1, &buf.Id)
		gl.BindBuffer(gl.ARRAY_BUFFER, buf.Id)
		gl.BufferData(gl.ARRAY_BUFFER,
			len(buf.Buffer)*sizeOfFloat,
			gl.Ptr(buf.Buffer),
			uint32(buf.Usage))
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, buf.Id)
}

func UnbindArrayBuffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
