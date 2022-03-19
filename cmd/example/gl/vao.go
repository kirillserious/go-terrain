package gl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// vao is a link between array buffer and program args location
//
// It is strange thing but commonly used in the OpenGL world
type vao struct {
	Id uint32

	Buf *ArrayBuffer
}

func NewVAO(buf *ArrayBuffer) (link *vao) {
	link = new(vao)
	link.Buf = buf
	gl.GenVertexArrays(1, &link.Id)
	return
}

func (link *vao) Bind() {
	gl.BindVertexArray(link.Id)
}

func (link *vao) LinkWithLocation(loc ProgramArgumentLocation, size, offset int) {
	link.Bind()
	defer UnbindVAO()
	link.Buf.Bind()
	defer UnbindArrayBuffer()
	gl.EnableVertexAttribArray(uint32(loc))
	gl.VertexAttribPointerWithOffset(uint32(loc), int32(size), gl.FLOAT,
		false, int32(link.Buf.Stride*sizeOfFloat), uintptr(offset*sizeOfFloat),
	)
}

func UnbindVAO() {
	gl.BindVertexArray(0)
}
