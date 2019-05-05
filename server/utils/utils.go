package utils

import (
	"bufio"
	"io"
	"os"
	"time"
)

// NoError raises an error if it not nil
func NoError(e error) {
	if e != nil {
		panic(e)
	}
}

func GetValidResult(result interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return result
}

// ParseTime "2018-01-03 15:46:43"
func ParseTime(str string) (time.Time, error) {
	local, err := time.LoadLocation("Asia/Chongqing")
	NoError(err)
	return time.ParseInLocation("2006-01-02 15:04:05", str, local)
}

func IsOutCache(datetime time.Time) bool {
	return time.Now().Sub(datetime) >= 1*time.Hour
}

func FileBufferedReader(filePath string) (io.Reader, error) {
	if file, error := os.Open(filePath); error != nil {
		return nil, error
	} else {
		return bufio.NewReader(file), nil
	}
}
