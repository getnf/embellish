package utils

import (
	"regexp"
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

func StringToInt(version string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	versionCleaned := re.FindAllString(version, -1)
	versionInt, err := strconv.Atoi(strings.Join(versionCleaned[:], ""))
	if err != nil {
		return 0, err
	}
	return versionInt, nil
}
