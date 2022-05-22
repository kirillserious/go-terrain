package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"terrain/internal"
)

var opts = struct {
	HeightMap string  `short:"m" long:"map" description:"Height map"`
	Texture   *string `short:"t" long:"texture" description:"Texture"`
}{}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	heights := internal.LoadHeightMap(opts.HeightMap)
	terrain := internal.NewTerrain(heights)

	window, terminate := internal.NewWindow(800, 600)
	defer terminate()

	window.Render(func() {
		terrain.Draw(window.Camera(), window.Projection())
	})
}
