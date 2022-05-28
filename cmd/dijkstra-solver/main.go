package main

import (
	"encoding/json"
	"fmt"
	"os"
	"terrain/internal"
	algo "terrain/internal/algorithms"
	"terrain/internal/common"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
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
	usedNodes := map[common.Position]struct{}{}
	borderNodes := map[common.Position]struct{}{
		{I: opts.ToI, J: opts.ToJ}: {},
	}
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
	bar.SetTemplateString(`{{ bar . }} {{percent .}} {{ rtime .}} {{ etime . }}`)
	for len(borderNodes) != 0 {
		bar.Increment()

		minDist, minPosition := float32(-1), common.Position{I: 0, J: 0}
		for borderNode := range borderNodes {
			dist := dists.At(borderNode.I, borderNode.J)
			if minDist < -0.5 || (minDist > -0.5 && dist < minDist) {
				minDist, minPosition = dist, borderNode
			}
		}
		usedNodes[minPosition] = struct{}{}
		if minPosition.I == opts.FromI && minPosition.J == opts.FromJ {
			break
		}
		delete(borderNodes, minPosition)
		for i := 0; i < algo.DirectionCount; i++ {
			cost := field.Length(minPosition.I, minPosition.J, algo.Direction(i))
			if cost == nil {
				continue
			}
			iDir, jDir := algo.DirectionToIndexes(minPosition.I, minPosition.J, algo.Direction(i))
			if _, ok := usedNodes[common.Position{I: iDir, J: jDir}]; ok {
				continue
			}
			costDir, newCostDir := dists.At(iDir, jDir), dists.At(minPosition.I, minPosition.J)+*cost
			if costDir < -0.5 || (costDir > -0.5 && newCostDir < costDir) {
				dists.SetAt(iDir, jDir, newCostDir)
			}
			borderNodes[common.Position{I: iDir, J: jDir}] = struct{}{}
		}
	}
	bar.Finish()
	fmt.Printf("Total cost: %0.2f\n", dists.At(opts.FromI, opts.FromJ))
	result := make([]common.Position, 0)
	for i, j := opts.FromI, opts.FromJ; i != opts.ToI || j != opts.ToJ; {
		result = append(result, common.Position{I: i, J: j})
		minDist, minPosition := float32(-1), common.Position{}
		for dir := 0; dir < algo.DirectionCount; dir++ {
			iDir, jDir := algo.DirectionToIndexes(i, j, algo.Direction(dir))
			if !field.IsValidIndex(iDir, jDir) {
				continue
			}
			dist := dists.At(iDir, jDir)
			if minDist < -0.5 || (minDist > -0.5 && dist > -0.5 && dist < minDist) {
				minDist, minPosition = dist, common.Position{I: iDir, J: jDir}
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
