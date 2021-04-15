package utils

func Interpolate(i, minFrom, maxFrom int, minTo, maxTo float64) float64 {
	return float64(i) * (maxTo - minTo) / float64(maxFrom-minFrom)
}
