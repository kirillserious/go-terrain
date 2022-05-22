package internal

import (
	"encoding/json"
	"github.com/aquilax/go-perlin"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type HeightMap struct {
	Heights []float32
	Stride  int
}

func GenerateHeightMap(iMax, jMax int) (heights *HeightMap) {
	heights = new(HeightMap)
	heights.Stride = jMax
	heights.Heights = make([]float32, iMax*jMax)
	p := perlin.NewPerlin(4., 2., 7, time.Now().UnixMilli())
	noise := func(x, y float64) float32 {
		return float32(p.Noise2D(x, y))
	}
	for i := 0; i < iMax; i++ {
		ParallelFor(0, jMax, 1, func(j int) {
			ni, nj := float64(i)/float64(iMax)-0.5, float64(j)/float64(jMax)-0.5
			heights.Heights[i*jMax+j] = 15*noise(ni, nj) +
				7*noise(2*ni, 2*nj) +
				3*noise(4*ni, 2*nj) +
				3*noise(2*ni, 4*nj) +
				noise(7*ni, 7*nj) +
				3*noise(10*ni, 8*nj)
			heights.Heights[i*jMax+j] *= float32(jMax) / 50
		})
	}
	return
}

func LoadHeightMap(filePath string) (heights *HeightMap) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.WithError(err).Panic("failed to load height map from file")
	}
	heights = new(HeightMap)
	if err = json.Unmarshal([]byte(file), heights); err != nil {
		log.WithError(err).Panic("failed to load height map from file")
	}
	return
}

// Flush writes the height map data into the provided writer
func (hm *HeightMap) Flush(writer io.Writer) (err error) {
	l := log.WithField("fcn", "(*HeightMap)Flush")

	data, err := json.Marshal(hm)
	if err != nil {
		l.WithError(err).Error("Failed to marshal the heights data")
		return
	}
	_, err = writer.Write(data)
	return
}

func (hm *HeightMap) FlushToFile(path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		log.WithError(err).Error("Failed to open the destination file")
		return
	}
	defer file.Close()
	err = hm.Flush(file)
	return
}

func ParallelFor(from, to, step int, fcn func(i int)) {
	var wg sync.WaitGroup
	wg.Add((to - from) / step)
	for i := from; i < to; i += step {
		go func(i int) {
			defer wg.Done()
			fcn(i)
		}(i)
	}
	wg.Wait()
}
