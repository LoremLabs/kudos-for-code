package common

import (
	"math"
	"regexp"
)

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func ExtractPackageName(id string) string {
	re := regexp.MustCompile(`\::(.*?)\:`)
	match := re.FindStringSubmatch(id)

	if len(match) > 0 {
		return match[1]
	}

	return ""
}
