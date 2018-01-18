package utils

import "time"

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
	return time.Parse("2006-01-02 15:04:05", str)
}

func IsOutCache(datetime time.Time) bool {
	return time.Now().Sub(datetime) >= 24*time.Hour
}
