// Package crossfade
//
// Tutorial #2 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_03_Timer_Related_Issues.shtml
//
// Outras interpolações: https://easings.net/
//
package crossfade

import (
	"bytes"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moniquelive/demoscenetuts/internal/utils"
)

var (
	step = 0.025
	k    = 0.0
)

type Cross struct {
	screenWidth  int
	screenHeight int
	at           image.Image
	ww           image.Image
}

func (c Cross) Draw(buffer *image.RGBA) {
	if (k < 0) || (k > 1) {
		step *= -1
	}
	k += step * (ebiten.CurrentTPS() / 100.0)
	blendLerp(c.ww, c.at, utils.Constrain(k, 0, 1), buffer)
}

func (c *Cross) Setup() (int, int, int) {
	c.at = loadFile("at.png")
	c.ww = loadFile("ww.png")
	c.screenWidth = c.at.Bounds().Dx()
	c.screenHeight = c.at.Bounds().Dy()
	return c.screenWidth, c.screenHeight, 1
}

func loadFile(filename string) image.Image {
	var err error

	var f []byte
	if f, err = ioutil.ReadFile(filename); err != nil {
		log.Fatal(err)
	}
	var img image.Image
	if img, _, err = image.Decode(bytes.NewReader(f)); err != nil {
		log.Fatal(err)
	}
	return img
}

func blendLerp(img1, img2 image.Image, k float64, r *image.RGBA) {
	w := img1.Bounds().Dx()
	h := img1.Bounds().Dy()
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			r.Pix[i+0] = uint8(myLerp(float64(r1>>8), float64(r2>>8), k))
			r.Pix[i+1] = uint8(myLerp(float64(g1>>8), float64(g2>>8), k))
			r.Pix[i+2] = uint8(myLerp(float64(b1>>8), float64(b2>>8), k))
			r.Pix[i+3] = uint8(myLerp(float64(a1>>8), float64(a2>>8), k))
			i += 4
		}
	}
}

func myLerp(a, b, k float64) float64 {
	return a + (b-a)*k //easeInOutBack(k)
}

//func easeInOutBack(x float64) float64 {
//	const c1 = 1.70158
//	const c2 = c1 * 1.525
//	if x < 0.5 {
//		return (math.Pow(2*x, 2) * ((c2+1)*2*x - c2)) / 2
//	} else {
//		return (math.Pow(2*x-2, 2)*((c2+1)*(x*2-2)+c2) + 2) / 2
//	}
//}
