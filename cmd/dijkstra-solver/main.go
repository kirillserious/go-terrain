package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"image"
	"os"
	"terrain/internal"
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
	field := NewField(heights, rgba)
	usedNodes := map[common.Position]struct{}{}
	borderNodes := map[common.Position]struct{}{
		common.Position{opts.ToI, opts.ToJ}: struct{}{},
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
	for len(borderNodes) != 0 {
		bar.Increment()

		minDist, minPosition := float32(-1), common.Position{0, 0}
		for borderNode := range borderNodes {
			dist := dists.At(borderNode.I, borderNode.J)
			if minDist < -0.5 || dist < minDist {
				minDist, minPosition = dist, borderNode
			}
		}
		usedNodes[minPosition] = struct{}{}
		delete(borderNodes, minPosition)
		for i := 0; i < DirectionCount; i++ {
			cost := field.Length(minPosition.I, minPosition.J, Direction(i))
			if cost == nil {
				continue
			}
			iDir, jDir := DirectionToIndexes(minPosition.I, minPosition.J, Direction(i))
			if _, ok := usedNodes[common.Position{iDir, jDir}]; ok {
				continue
			}
			costDir, newCostDir := dists.At(iDir, jDir), dists.At(minPosition.I, minPosition.J)+*cost
			if costDir < -0.5 || newCostDir < costDir {
				dists.SetAt(iDir, jDir, newCostDir)
			}
			borderNodes[common.Position{iDir, jDir}] = struct{}{}
		}
	}
	bar.Finish()
	fmt.Printf("Total cost: %0.2f\n", dists.At(opts.FromI, opts.FromJ))
	result := make([]common.Position, 0)
	for i, j := opts.FromI, opts.FromJ; i != opts.ToI && j != opts.ToJ; {
		result = append(result, common.Position{i, j})
		minDist, minPosition := float32(-1), common.Position{}
		for dir := 0; dir < DirectionCount; dir++ {
			iDir, jDir := DirectionToIndexes(i, j, Direction(dir))
			if !field.isValidIndex(iDir, jDir) {
				continue
			}
			dist := dists.At(iDir, jDir)
			if minDist < -0.5 || dist < minDist {
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

type Field struct {
	HeightMap *internal.HeightMap
	RGBA      *image.RGBA
}

func NewField(heights *internal.HeightMap, rgba *image.RGBA) (field *Field) {
	iMax, jMax := heights.Bounds()
	if rgba.Rect.Size().X != iMax || rgba.Rect.Size().Y != jMax {
		log.Panic("Incorrect sizes")
	}
	field = new(Field)
	field.HeightMap = heights
	field.RGBA = rgba
	return
}

func (field *Field) Bounds() (iMax, jMax int) {
	return field.HeightMap.Bounds()
}

type Direction int

const DirectionCount int = 8
const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

type Color int

const (
	UnknownColor Color = iota
	BlackColor
)

func (field *Field) color(i, j int) Color {
	color := field.RGBA.RGBAAt(i, j)
	switch {
	case color.R == 0 && color.G == 0 && color.B == 0:
		return BlackColor
	default:
		return UnknownColor
	}
}

func DirectionToIndexes(i, j int, dir Direction) (int, int) {
	switch dir {
	case North:
		return i - 1, j
	case NorthEast:
		return i - 1, j + 1
	case East:
		return i, j + 1
	case SouthEast:
		return i + 1, j + 1
	case South:
		return i + 1, j
	case SouthWest:
		return i + 1, j - 1
	case West:
		return i, j - 1
	case NorthWest:
		return i - 1, j - 1
	default:
		return -1, -1
	}
}

func (field *Field) isValidIndex(i, j int) bool {
	if i < 0 || j < 0 {
		return false
	}
	iMax, jMax := field.HeightMap.Bounds()
	if i >= iMax || j >= jMax {
		return false
	}
	if field.color(i, j) == BlackColor {
		return false
	}
	return true
}

const requireCost float32 = 10

func (field *Field) Length(i, j int, dir Direction) *float32 {
	iTo, jTo := DirectionToIndexes(i, j, dir)
	if !field.isValidIndex(iTo, jTo) {
		return nil
	}
	heightFrom, heightTo := field.HeightMap.At(i, j), field.HeightMap.At(iTo, jTo)
	if heightTo > heightFrom {
		result := (heightTo-heightFrom)*5 + requireCost
		return &result
	} else {
		result := requireCost
		return &result
	}
}
