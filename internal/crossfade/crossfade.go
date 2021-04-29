// Package crossfade
//
// Tutorial #2 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_03_Timer_Related_Issues.shtml
//
// Outras interpolações: https://easings.net/
//
package crossfade

import (
	"image"
	_ "image/png"

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
	at           *image.RGBA
	ww           *image.RGBA
}

func (c Cross) Draw(buffer *image.RGBA) {
	if (k < 0) || (k > 1) {
		step *= -1
	}
	k += step * ebiten.CurrentTPS() / 100.0
	blendLerp(c.ww, c.at, utils.Constrain(k, 0, 1), buffer)
}

func (c *Cross) Setup() (int, int, int) {
	c.at = utils.LoadFileRGBA("at.png")
	c.ww = utils.LoadFileRGBA("ww.png")
	c.screenWidth = c.at.Bounds().Dx()
	c.screenHeight = c.at.Bounds().Dy()
	return c.screenWidth, c.screenHeight, 1
}

func blendLerp(img1, img2 *image.RGBA, k float64, r *image.RGBA) {
	bb := img1.Bounds().Dx() * img1.Bounds().Dy() * 4
	for i := 0; i < bb; i++ {
		f1 := float64(img1.Pix[i])
		f2 := float64(img2.Pix[i])
		r.Pix[i] = uint8(myLerp(f1, f2, k))
	}
}

func myLerp(a, b, k float64) float64 {
	return a + (b-a)*k //easeInOutBack(k)
}

//func easeInOutBack(x float64) float64 {
//	const c1 = 1.70158
//	const c2 = c1 * 1.525
//	if x < 0.5 {
//		return ((2 * x * 2 * x) * ((c2+1)*2*x - c2)) / 2
//	} else {
//		return (((2*x-2)*(2*x-2))*((c2+1)*(x*2-2)+c2) + 2) / 2
//	}
//}
