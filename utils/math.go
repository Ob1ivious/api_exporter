package utils

import (
	"fmt"
	"strconv"
)

func Avg(arr []float64) float64 {
	var sum float64

	for _, item := range arr {
		sum += item
	}

	return sum / float64(len(arr))
}

func Round(value float64, precision int) float64 {
	format := "%." + strconv.Itoa(precision) + "f"
	value, _ = strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return value
}
