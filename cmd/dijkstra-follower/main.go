package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"terrain/internal"
	algo "terrain/internal/algorithms"
	"terrain/internal/common"
)

const (
	HeightsPath string = "/mnt/hm.json"
	TexturePath string = "/mnt/texture.png"
)

var opts = struct {
	FromI int `long:"from-i" required:"yes"`
	FromJ int `long:"from-j" required:"yes"`
	ToI   int `long:"to-i" required:"yes"`
	ToJ   int `long:"to-j" required:"yes"`
}{}

func handler() func(http.ResponseWriter, *http.Request) {
	heights := internal.LoadHeightMap(HeightsPath)
	rgba := common.LoadRGBA(TexturePath)

	// Prepare types
	field := algo.NewField(heights, rgba)
	usedNodes := map[common.Position]struct{}{}
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func main() {
	heights := internal.LoadHeightMap(HeightsPath)
	rgba := common.LoadRGBA(TexturePath)

	// Prepare types
	field := algo.NewField(heights, rgba)
	usedNodes := map[common.Position]struct{}{}

	numCPU := runtime.NumCPU()
	borderNodes := make([]map[common.Position]struct{}, numCPU)
	for i := range borderNodes {
		borderNodes[i] = make(map[common.Position]struct{})
	}
	borderNodes[0][common.Position{opts.ToI, opts.ToJ}] = struct{}{}
	borderLen := func() (count int) {
		for _, batch := range borderNodes {
			count += len(batch)
		}
		return
	}
	borderAdd := func(pos common.Position) {
		count, idx := len(borderNodes[0]), 0
		for i, batch := range borderNodes {
			if _, ok := batch[pos]; ok {
				return
			}
			if len(batch) < count {
				count, idx = len(batch), i
			}
		}
		borderNodes[idx][pos] = struct{}{}
	}
	borderRemove := func(pos common.Position, batchIdx int) {
		delete(borderNodes[batchIdx], pos)
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

	for borderLen() != 0 {
		bar.Increment()

		batchResult := make([]common.Position, numCPU)
		internal.ParallelFor(0, numCPU, 1, func(batchNum int) {
			minDist, minPosition := float32(-1), common.Position{-1, -1}
			for borderNode := range borderNodes[batchNum] {
				dist := dists.At(borderNode.I, borderNode.J)
				if minDist < -0.5 || (minDist > -0.5 && dist < minDist) {
					minDist, minPosition = dist, borderNode
				}
			}
			batchResult[batchNum] = minPosition
		})

		minDist, minPosition, minBatch := float32(-1), common.Position{0, 0}, 0
		for idx, pos := range batchResult {
			if pos.I == -1 {
				continue
			}
			dist := dists.At(pos.I, pos.J)
			if minDist < -0.5 || (minDist > -0.5 && dist < minDist) {
				minDist, minPosition, minBatch = dist, pos, idx
			}
		}
		usedNodes[minPosition] = struct{}{}
		if minPosition.I == opts.FromI && minPosition.J == opts.FromJ {
			break
		}
		borderRemove(minPosition, minBatch)
		for i := 0; i < algo.DirectionCount; i++ {
			cost := field.Length(minPosition.I, minPosition.J, algo.Direction(i))
			if cost == nil {
				continue
			}
			iDir, jDir := algo.DirectionToIndexes(minPosition.I, minPosition.J, algo.Direction(i))
			if _, ok := usedNodes[common.Position{iDir, jDir}]; ok {
				continue
			}
			costDir, newCostDir := dists.At(iDir, jDir), dists.At(minPosition.I, minPosition.J)+*cost
			if costDir < -0.5 || (costDir > -0.5 && newCostDir < costDir) {
				dists.SetAt(iDir, jDir, newCostDir)
			}
			borderAdd(common.Position{iDir, jDir})
		}
	}
	bar.Finish()
	fmt.Printf("Total cost: %0.2f\n", dists.At(opts.FromI, opts.FromJ))
	result := make([]common.Position, 0)
	for i, j := opts.FromI, opts.FromJ; i != opts.ToI && j != opts.ToJ; {
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
