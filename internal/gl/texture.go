package gl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	log "github.com/sirupsen/logrus"
	"image"
	"terrain/internal/common"
)

// TODO: Add up to 32 textures, we do not need it now
var textureSlots = []int{gl.TEXTURE0, gl.TEXTURE1, gl.TEXTURE2}
var usedSlots = map[int]struct{}{}

type Texture struct {
	slot    int
	texture uint32
}

func NewTexture(rgba *image.RGBA) (texture *Texture) {
	texture = new(Texture)
	for _, slot := range textureSlots {
		if _, ok := usedSlots[slot]; !ok {
			usedSlots[slot] = struct{}{}
			texture.slot = slot
			break
		}
	}
	if texture.slot == 0 {
		log.Panic("limit of textures")
	}

	gl.GenTextures(1, &texture.texture)
	gl.ActiveTexture(uint32(texture.slot))
	gl.BindTexture(gl.TEXTURE_2D, texture.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	return
}

func NewTextureFromFile(path string) (texture *Texture) {
	rgba := common.LoadRGBA(path)
	return NewTexture(rgba)
}

func (texture *Texture) Bind() {
	gl.ActiveTexture(uint32(texture.slot))
	gl.BindTexture(gl.TEXTURE_2D, texture.texture)
}
