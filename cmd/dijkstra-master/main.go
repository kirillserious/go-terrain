package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"terrain/internal"
	algo "terrain/internal/algorithms"
	"terrain/internal/common"
)

const (
	FromIEnv       string = "FROM_I"
	FromJEnv       string = "FROM_J"
	ToIEnv         string = "TO_I"
	ToJEnv         string = "TO_J"
	FollowerIPsEnv string = "FOLLOWER_IPS"
)

const (
	HeightsPath string = "/mnt/hm.json"
	TexturePath string = "/mnt/texture.png"
	OutputPath  string = "/mnt/path.json"
)

var FollowersIPs []string = strings.Split(os.Getenv(FollowerIPsEnv), " ")
var FromI, FromJ, ToI, ToJ int

func init() {
	var err error
	FromI, err = strconv.Atoi(os.Getenv(FromIEnv))
	if err != nil {
		log.WithError(err).Panic("failed to parse environment")
	}
	FromJ, err = strconv.Atoi(os.Getenv(FromJEnv))
	if err != nil {
		log.WithError(err).Panic("failed to parse environment")
	}
	ToI, err = strconv.Atoi(os.Getenv(ToIEnv))
	if err != nil {
		log.WithError(err).Panic("failed to parse environment")
	}
	ToJ, err = strconv.Atoi(os.Getenv(ToJEnv))
	if err != nil {
		log.WithError(err).Panic("failed to parse environment")
	}
}

type request struct {
	Use common.Position
	Add map[common.Position]float32
}

type response struct {
	MinPosition common.Position
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
	borderNodes[0][common.Position{ToI, ToJ}] = struct{}{}
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
	dists.SetAt(ToI, ToJ, 0)

	bar := pb.StartNew(iMax * jMax)
	bar.SetTemplateString(`{{ bar . }} {{percent .}} {{ rtime .}} {{ etime . }}`)

	message := request{Use: common.Position{I: -1, J: -1}}
	for borderLen() != 0 {
		bar.Increment()

		batchResult := make([]common.Position, len(FollowersIPs))
		internal.ParallelFor(0, len(FollowersIPs), 1, func(batchNum int) {
			msgData, err := json.Marshal(message)
			if err != nil {
				log.WithError(err).Panic("failed to marshal message")
			}
			req, err := http.NewRequest(http.MethodPost,
				FollowersIPs[batchNum],
				bytes.NewBuffer(msgData),
			)
			if err != nil {
				log.WithError(err).Panic("failed to create request")
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.WithError(err).Panic("failed to do request")
			}
			defer resp.Body.Close()

			respData, err := io.ReadAll(resp.Body)
			if err != nil {
				log.WithError(err).Panic("failed to read response body")
			}
			respMsg := new(response)
			if err := json.Unmarshal(respData, respMsg); err != nil {
				log.WithError(err).Panic("failed to unmarshal response")
			}

			batchResult[batchNum] = respMsg.MinPosition
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
		usedNodes[minPosition], message.Use = struct{}{}, minPosition
		if minPosition.I == FromI && minPosition.J == FromJ {
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
			message.Add = append(message.Add, common.Position{iDir, jDir})
		}
	}
	bar.Finish()
	fmt.Printf("Total cost: %0.2f\n", dists.At(FromI, FromJ))
	result := make([]common.Position, 0)
	for i, j := FromI, FromJ; i != ToI && j != ToJ; {
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

	file, err := os.Create(OutputPath)
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
