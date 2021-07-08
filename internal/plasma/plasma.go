// Package plasma
// Tutorial #3 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_04_Per_Pixel_Control.shtml
package plasma

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed plasma.png
var bgBytes []byte

type Plasma struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	bg           *image.RGBA
	plasma1      *image.RGBA
	plasma2      *image.RGBA
}

func (c *Plasma) Draw(buffer *image.RGBA) {
	c.update(buffer)
}

func (c *Plasma) Setup() (int, int, int) {
	c.bg = utils.LoadBufferRGBA(bgBytes)
	c.screenWidth = c.bg.Bounds().Dx()
	c.screenHeight = c.bg.Bounds().Dy()

	doubleBounds := image.Rect(0, 0, c.screenWidth*2, c.screenHeight*2)
	c.plasma1 = image.NewRGBA(doubleBounds)
	c.plasma2 = image.NewRGBA(doubleBounds)
	c.plasma()
	return c.screenWidth, c.screenHeight, 2
}

func (c Plasma) f1(i, j int) uint8 {
	// 64 + 63 * ( sin( hypot( 200-j, 320-i )/16 ) ) )
	ii := float64(i)
	jj := float64(j)
	ww := float64(c.screenWidth)
	hh := float64(c.screenHeight)
	return 128 + uint8(127*(math.Sin(math.Hypot(hh-jj, ww-ii)/16.0)))
}

func (c Plasma) f2(i, j int) uint8 {
	// 64 + 63 * sin( i/(37+15*cos(j/74)) ) * cos( j/(31+11*sin(i/57))) )
	ii := float64(i)
	jj := float64(j)
	return 128 + uint8(127.0*math.Sin(ii/(37.0+15.0*math.Cos(jj/74.0)))*math.Cos(jj/(31.0+11.0*math.Sin(ii/57.0))))
}

func (c Plasma) plasma() {
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

func (c *Plasma) update(r *image.RGBA) {
	w := r.Bounds().Dx()
	h := r.Bounds().Dy()
	// move plasma with more sine functions :)
	halfW := float64(w) / 2.0
	halfH := float64(h) / 2.0
	x1 := int(halfW) + int((halfW-1)*math.Cos(float64(c.frameCount)/97.0))
	x2 := int(halfW) + int((halfW-1)*math.Sin(float64(-c.frameCount)/114.0))
	x3 := int(halfW) + int((halfW-1)*math.Sin(float64(-c.frameCount)/137.0))

	y1 := int(halfH) + int((halfH-1)*math.Sin(float64(c.frameCount)/123.0))
	y2 := int(halfH) + int((halfH-1)*math.Cos(float64(-c.frameCount)/75.0))
	y3 := int(halfH) + int((halfH-1)*math.Cos(float64(-c.frameCount)/108.0))

	src1 := 4 * (y1*c.plasma1.Bounds().Dx() + x1)
	src2 := 4 * (y2*c.plasma2.Bounds().Dx() + x2)
	src3 := 4 * (y3*c.plasma1.Bounds().Dx() + x3)

	loc := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			for rgb := 0; rgb < 4; rgb++ {
				r.Pix[loc+rgb] = uint8(
					utils.ConstrainU32(
						uint32(c.bg.Pix[loc+rgb])+
							uint32(c.plasma1.Pix[src1+rgb])>>3+
							uint32(c.plasma2.Pix[src2+rgb])>>3+
							uint32(c.plasma2.Pix[src3+rgb])>>3, 0, 255))
			}
			loc += 4
			src1 += 4
			src2 += 4
			src3 += 4
		}
		src1 += w * 4
		src2 += w * 4
		src3 += w * 4
	}
	c.frameCount++
}
