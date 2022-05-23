package algorithms

import (
	log "github.com/sirupsen/logrus"
	"image"
	"terrain/internal"
)

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

func (field *Field) IsValidIndex(i, j int) bool {
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
	if !field.IsValidIndex(iTo, jTo) {
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
