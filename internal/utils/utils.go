package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
)

func Constrain(i, low, high float64) float64 {
	return math.Max(math.Min(i, high), low)
}

func ConstrainI(i, low, high int) int {
	return int(math.Max(math.Min(float64(i), float64(high)), float64(low)))
}

func ConstrainU8(i, low, high uint8) uint8 {
	return uint8(math.Max(math.Min(float64(i), float64(high)), float64(low)))
}

func ConstrainU32(i, low, high uint32) uint32 {
	return uint32(math.Max(math.Min(float64(i), float64(high)), float64(low)))
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

func LoadFileRGBA(filename string) *image.RGBA {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba
}

func LoadBufferRGBA(b []byte) *image.RGBA {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba
}

func LoadBufferGray(b []byte) *image.Gray {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, image.Point{}, draw.Src)
	return gray
}

func LoadBufferPaletted(b []byte) *image.Paletted {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	pal := image.NewPaletted(img.Bounds(), img.ColorModel().(color.Palette))
	draw.Draw(pal, pal.Bounds(), img, image.Point{}, draw.Src)
	return pal
}

func Lerp(a, b, k float64) float64 {
	return a + (b-a)*k //easeInOutBack(k)
}
