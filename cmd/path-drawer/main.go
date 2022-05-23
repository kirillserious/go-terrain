package main

import (
	"encoding/json"
	"github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"image/color"
	"image/png"
	_ "image/png"
	"io"
	"os"
	"terrain/internal/common"
)

var opts = struct {
	Texture string `short:"t" long:"texture" required:"yes"`
	Path    string `short:"p" long:"path" required:"yes"`
	Out     string `short:"o" long:"out" required:"yes"`
}{}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}
	rgba := common.LoadRGBA(opts.Texture)
	path := make([]common.Position, 0)
	file, err := os.Open(opts.Path)
	if err != nil {
		log.WithError(err).Panic("failed to open path file")
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.WithError(err).Panic("failed to read path file")
	}
	json.Unmarshal(data, &path)

	maxI, maxJ := rgba.Rect.Size().X, rgba.Rect.Size().Y
	maxK := 1
	if maxI > 500 {
		maxK = 3
	}

	bar := pb.StartNew(len(path))
	for _, pos := range path {
		bar.Increment()
		rgba.SetRGBA(pos.I, pos.J, color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff})
		for l := 0; l < maxK; l++ {
			for k := 0; k < maxK; k++ {
				fcn := func(i, j int) {
					if i < 0 || j < 0 || i >= maxI || j >= maxJ {
						return
					}
					rgba.SetRGBA(i, j, color.RGBA{R: 0xff, G: 0, B: 0, A: 0xff})
				}
				i, j := pos.I+l, pos.J+k
				fcn(i, j)
				i, j = pos.I+l, pos.J-k
				fcn(i, j)
				i, j = pos.I-l, pos.J+k
				fcn(i, j)
				i, j = pos.I-l, pos.J-k
				fcn(i, j)
			}
		}
	}
	bar.Finish()

	f, _ := os.Create(opts.Out)
	png.Encode(f, rgba)
}
