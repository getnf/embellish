package utils

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func Filter[T any](items []T, condition func(T) bool) (results []T) {
	for _, item := range items {
		if condition(item) {
			results = append(results, item)
		}
	}
	return
}

func Fold[T, U any](data []T, f func(T) U) []U {

	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func StringToInt(version string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	versionCleaned := re.FindAllString(version, -1)
	versionInt, err := strconv.Atoi(strings.Join(versionCleaned[:], ""))
	if err != nil {
		return 0, err
	}
	return versionInt, nil
}

func FontNameWithoutExtention(name string) string {
	return strings.Split(name, ".")[0]
}

func OsType() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "unsupported"
	}
}
