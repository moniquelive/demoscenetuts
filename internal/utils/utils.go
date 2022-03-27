package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"

	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

func Max[T Numeric](x, y T) T {
	return T(math.Max(float64(x), float64(y)))
}

func Min[T Numeric](x, y T) T {
	return T(math.Min(float64(x), float64(y)))
}

func Constrain[T Numeric](i, low, high T) T {
	return Max(Min(i, high), low)
}

func Interpolate[T Numeric](i, start1, stop1, start2, stop2 T) T {
	return (i-start1)/(stop1-start1)*(stop2-start2) + start2
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

func Lerp[T Numeric](a, b, k T) T {
	return a + (b-a)*k //easeInOutBack(k)
}

func Memset(a *[64000]uint16, v uint16) {
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}
