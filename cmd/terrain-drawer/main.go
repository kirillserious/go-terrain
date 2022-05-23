package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"terrain/internal"
	gl2 "terrain/internal/gl"
)

var opts = struct {
	HeightMap string `short:"m" long:"map" description:"Height map" required:"yes"`
	Texture   string `short:"t" long:"texture" description:"Texture" required:"yes"`
}{}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	window, terminate := internal.NewWindow(800, 600)
	defer terminate()

	heights := internal.LoadHeightMap(opts.HeightMap)
	texture := gl2.NewTextureFromFile(opts.Texture)

	terrain := internal.NewTerrain2(heights, texture)

	window.Render(func() {
		terrain.Draw(window.Camera(), window.Projection())
	})
}
