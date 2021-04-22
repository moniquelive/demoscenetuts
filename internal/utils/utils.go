package utils

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func Constrain(i, low, high float64) float64 {
	return math.Max(math.Min(i, high), low)
}

func ConstrainI(i, low, high int) int {
	return int(math.Max(math.Min(float64(i), float64(high)), float64(low)))
}

func Interpolate(i, start1, stop1, start2, stop2 float64) float64 {
	//return i * (maxTo - minTo) / (maxFrom - minFrom)
	newVal := (i-start1)/(stop1-start1)*(stop2-start2) + start2
	return newVal
	//if start2 < stop2 {
	//	return Constrain(newVal, start2, stop2)
	//} else {
	//	return Constrain(newVal, stop2, start2)
	//}
}

func Fill(buffer *image.RGBA, rgba color.RGBA) {
	bg := image.NewUniform(rgba)
	draw.Draw(buffer, buffer.Bounds(), bg, image.Point{}, draw.Src)
}
