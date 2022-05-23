package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"terrain/internal"
	algo "terrain/internal/algorithms"
	"terrain/internal/common"
)

var opts = struct {
	HeightMap string `short:"m" long:"height_map" required:"yes"`
	Texture   string `short:"t" long:"texture" required:"yes"`
	FromI     int    `long:"from-i" required:"yes"`
	FromJ     int    `long:"from-j" required:"yes"`
	ToI       int    `long:"to-i" required:"yes"`
	ToJ       int    `long:"to-j" required:"yes"`
	Out       string `short:"o" long:"out" required:"yes"`
}{}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	heights := internal.LoadHeightMap(opts.HeightMap)
	rgba := common.LoadRGBA(opts.Texture)

	// Prepare types
	field := algo.NewField(heights, rgba)
	iMax, jMax := field.Bounds()
	// Fill inf as -1
	dists := internal.EmptyHeightMap(iMax, jMax)
	for i := 0; i < iMax; i++ {
		for j := 0; j < jMax; j++ {
			dists.SetAt(i, j, float32(-1))
		}
	}
	dists.SetAt(opts.ToI, opts.ToJ, 0)

	bar := pb.StartNew(iMax * jMax)

	stop := iMax*jMax - 2
	for k := 0; k < stop; k++ {
		bar.Increment()
		for i := 0; i < iMax; i++ {
			for j := 0; j < jMax; j++ {
				dist := dists.At(i, j)
				if dist < -0.5 {
					continue
				}
				for dir := 0; dir < algo.DirectionCount; dir++ {
					iDir, jDir := algo.DirectionToIndexes(i, j, algo.Direction(dir))
					cost := field.Length(i, j, algo.Direction(dir))
					if cost == nil {
						continue
					}
					distDir := dists.At(iDir, jDir)
					newDist := dist + *cost
					if distDir < -0.5 || (distDir > -0.5 && newDist < distDir) {
						dists.SetAt(iDir, jDir, newDist)
					}
				}
			}
		}
	}
	bar.Finish()
	fmt.Printf("Total cost: %0.2f\n", dists.At(opts.FromI, opts.FromJ))
	//dists.FlushToFile(opts.Out)

	result := make([]common.Position, 0)
	count := 0
	for i, j := opts.FromI, opts.FromJ; i != opts.ToI && j != opts.ToJ; {
		count++
		result = append(result, common.Position{i, j})
		minDist, minPosition := float32(-1), common.Position{}
		for dir := 0; dir < algo.DirectionCount; dir++ {
			iDir, jDir := algo.DirectionToIndexes(i, j, algo.Direction(dir))
			if !field.IsValidIndex(iDir, jDir) {
				continue
			}
			dist := dists.At(iDir, jDir)
			if minDist < -0.5 || (minDist > -0.5 && dist > -0.5 && dist < minDist) {
				minDist, minPosition = dist, common.Position{iDir, jDir}
			}
		}
		i, j = minPosition.I, minPosition.J
	}
	file, err := os.Create(opts.Out)
	if err != nil {
		log.WithError(err).Panic("Failed to open the destination file")
		return
	}
	defer file.Close()
	data, err := json.Marshal(result)
	if err != nil {
		log.WithError(err).Panic("failed to marshal the result")
		return
	}
	_, err = file.Write(data)

}
