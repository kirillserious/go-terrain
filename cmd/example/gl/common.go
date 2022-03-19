package gl

import "github.com/go-gl/gl/v4.1-core/gl"

type Property int

const (
	CurrentProgramProperty Property = gl.CURRENT_PROGRAM
	VersionProperty                 = gl.VERSION
)

func GetIntegerProperty(property Property) (value int) {
	data := new(int32)
	gl.GetIntegerv(uint32(property), data)
	value = int(*data)
	return
}

func GetStringProperty(property Property) (value string) {
	value = gl.GoStr(gl.GetString(uint32(property)))
	return
}

func GetCurrentProgramId() int {
	return GetIntegerProperty(CurrentProgramProperty)
}

func GetVersion() string {
	return GetStringProperty(VersionProperty)
}
