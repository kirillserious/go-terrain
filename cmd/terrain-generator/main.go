package main

import (
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"terrain/internal"
)

var opts struct {
	Destination *string `long:"dst" default-mask:"STDOUT" description:"File path to store the result"`
	Height      int     `short:"h" long:"height" default:"1000" description:"The height map height"`
	Width       int     `short:"w" long:"width" default:"1000" description:"The height map width"`
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	var writer io.Writer = os.Stdout
	if opts.Destination != nil {
		file, err := os.Create(*opts.Destination)
		if err != nil {
			log.WithError(err).Error("Failed to open the destination file")
			os.Exit(2)
		}
		defer file.Close()
		writer = file
	}

	heights := internal.GenerateHeightMap(opts.Height, opts.Width)
	heights.Flush(writer)
}
