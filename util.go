package thor

import (
	"fmt"
	"path"
	"strconv"
)

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

// toString try to convert the argument into a string
func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

// toInt32 try to convert the argument into a int32
func toInt32(val interface{}) int32 {
	str := toString(val)
	r, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		r = 0
	}
	return int32(r)
}

// toUint32 try to convert the argument into a uint32
func toUint32(val interface{}) uint32 {
	str := toString(val)
	r, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		r = 0
	}
	return uint32(r)
}

// toFloat32 try to convert the argument into a float32
func toFloat32(val interface{}) float32 {
	str := toString(val)
	r, err := strconv.ParseFloat(str, 32)
	if err != nil {
		r = 0
	}
	return float32(r)
}

// toFloat64 try to convert the argument into a float64
func toFloat64(val interface{}) float64 {
	str := toString(val)
	r, err := strconv.ParseFloat(str, 64)
	if err != nil {
		r = 0
	}
	return r
}
