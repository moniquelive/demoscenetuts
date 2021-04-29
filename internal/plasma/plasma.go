// Package plasma
// Tutorial #3 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_04_Per_Pixel_Control.shtml
package plasma

import (
	"image"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

var (
	frameCount = 0
)

type Plasma struct {
	screenWidth  int
	screenHeight int
	bg           *image.RGBA
	plasma1      *image.RGBA
	plasma2      *image.RGBA
}

func (c Plasma) Draw(buffer *image.RGBA) {
	c.update(buffer)
}

func (c *Plasma) Setup() (int, int, int) {
	c.bg = utils.LoadFileRGBA("plasma.png")
	doubleBounds := image.Rect(0, 0, c.bg.Bounds().Dx()*2, c.bg.Bounds().Dy()*2)
	c.plasma1 = image.NewRGBA(doubleBounds)
	c.plasma2 = image.NewRGBA(doubleBounds)
	c.screenWidth = c.bg.Bounds().Dx()
	c.screenHeight = c.bg.Bounds().Dy()
	c.plasma()
	return c.screenWidth, c.screenHeight, 2
}

func (c Plasma) f1(i, j int) uint8 {
	// 64 + 63 * ( sin( hypot( 200-j, 320-i )/16 ) ) )
	ii := float64(i)
	jj := float64(j)
	ww := float64(c.screenWidth)
	hh := float64(c.screenHeight)
	return 64 + uint8(63*(math.Sin(math.Hypot(hh-jj, ww-ii)/16.0)))
}

func (c Plasma) f2(i, j int) uint8 {
	// 64 + 63 * sin( i/(37+15*cos(j/74)) ) * cos( j/(31+11*sin(i/57))) )
	ii := float64(i)
	jj := float64(j)
	return uint8(64.0 +
		63.0*math.Sin(
			ii/(37.0+
				15.0*math.Cos(jj/74.0)))*
			math.Cos(jj/(31.0+11.0*math.Sin(ii/57.0))))
}

func (c *Plasma) plasma() {
	w := c.plasma1.Bounds().Dx()
	h := c.plasma1.Bounds().Dy()
	loc := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c1 := c.f1(x, y)
			c.plasma1.Pix[loc+0] = c1
			c.plasma1.Pix[loc+1] = c1
			c.plasma1.Pix[loc+2] = c1
			c.plasma1.Pix[loc+3] = 255
			c2 := c.f2(x, y)
			c.plasma2.Pix[loc+0] = c2
			c.plasma2.Pix[loc+1] = c2
			c.plasma2.Pix[loc+2] = c2
			c.plasma2.Pix[loc+3] = 255
			loc += 4
		}
	}
}

func (c Plasma) update(r *image.RGBA) {
	// move plasma with more sine functions :)
	halfW := float64(r.Bounds().Dx()) / 2.0
	halfH := float64(r.Bounds().Dy()) / 2.0
	x1 := int(halfW) + int((halfW-1)*math.Cos(float64(frameCount)/97.0))
	x2 := int(halfW) + int((halfW-1)*math.Sin(float64(-frameCount)/114.0))
	x3 := int(halfW) + int((halfW-1)*math.Sin(float64(-frameCount)/137.0))

	y1 := int(halfH) + int((halfH-1)*math.Sin(float64(frameCount)/123.0))
	y2 := int(halfH) + int((halfH-1)*math.Cos(float64(-frameCount)/75.0))
	y3 := int(halfH) + int((halfH-1)*math.Cos(float64(-frameCount)/108.0))

	w1 := c.plasma1.Bounds().Dx()
	w2 := c.plasma2.Bounds().Dx()

	src1 := 4 * (y1*w1 + x1)
	src2 := 4 * (y2*w2 + x2)
	src3 := 4 * (y3*w1 + x3)

	loc := 0
	w := r.Bounds().Dx()
	h := r.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rr := uint32(c.bg.Pix[loc+0])
			gg := uint32(c.bg.Pix[loc+1])
			bb := uint32(c.bg.Pix[loc+2])
			rr1 := uint32(c.plasma1.Pix[src1+0])
			gg2 := uint32(c.plasma2.Pix[src2+1])
			bb2 := uint32(c.plasma2.Pix[src3+2])
			r.Pix[loc+0] = uint8(utils.ConstrainU32(rr+rr1, 0, 255))
			r.Pix[loc+1] = uint8(utils.ConstrainU32(gg+gg2, 0, 255))
			r.Pix[loc+2] = uint8(utils.ConstrainU32(bb+bb2, 0, 255))
			r.Pix[loc+3] = 255
			loc += 4
			src1 += 4
			src2 += 4
			src3 += 4
		}
		src1 += w * 4
		src2 += w * 4
		src3 += w * 4
	}
	frameCount++
}
