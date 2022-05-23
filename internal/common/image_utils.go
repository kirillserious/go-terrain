package common

import (
	log "github.com/sirupsen/logrus"
	"image"
	"image/draw"
	_ "image/png"
	"os"
)

func LoadRGBA(filePath string) *image.RGBA {
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.WithError(err).Panic("texture %q not found on disk", filePath)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.WithError(err).Panic("failed to decode image")
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba
}
