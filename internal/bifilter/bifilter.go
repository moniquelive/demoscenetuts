// Package bifilter
// Tutorial #6 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_06_Bitmap_Distortion.shtml
package bifilter

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed bifilter.png
var bgBytes []byte

type Bifilter struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	bg           *image.RGBA
	displaceX    []int8
	displaceY    []int8
	flip         bool
}

func (c *Bifilter) Draw(buffer *image.RGBA) {
	// move distortion buffer
	x1 := 160 + int(159.0*math.Cos(float64(c.frameCount)/205.0))
	x2 := 160 + int(159.0*math.Sin(float64(-c.frameCount)/197.0))
	y1 := 100 + int(99.0*math.Sin(float64(c.frameCount)/231.0))
	y2 := 100 + int(99.0*math.Cos(float64(-c.frameCount)/224.0))
	// draw the effct
	if c.flip {
		c.Distort(buffer, x1, y1, x2, y2)
	} else {
		c.DistortBili(buffer, x1, y1, x2, y2)
	}
	c.frameCount += 2
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.flip = !c.flip
	}
}

func (c *Bifilter) Setup() (int, int, int) {
	c.bg = utils.LoadBufferRGBA(bgBytes)
	bgw := c.bg.Bounds().Dx()
	bgh := c.bg.Bounds().Dy()

	c.displaceX = make([]int8, bgw*2*bgh*2)
	c.displaceY = make([]int8, bgw*2*bgh*2)
	c.screenWidth = bgw
	c.screenHeight = bgh
	c.precalculate()
	return c.screenWidth, c.screenHeight, 2
}

func (c Bifilter) precalculate() {
	dst := 0
	for j := 0; j < 400; j++ {
		for i := 0; i < 640; i++ {
			x := float64(i)
			y := float64(j)
			// notice the values contained in the buffers are signed
			// i.e. can be both positive and negative
			c.displaceX[dst] = (int8)(8 * (2 * (math.Sin(x/20) + math.Sin(x*y/2000) +
				math.Sin((x+y)/100) + math.Sin((y-x)/70) + math.Sin((x+4*y)/70) +
				2*math.Sin(math.Hypot(256-x, 150-y/8)/40))))
			// also notice we multiply by 8 to get 5.3 fixed point distortion
			// coefficients for our bilinear filtering
			c.displaceY[dst] = (int8)(8 * (math.Cos(x/31) + math.Cos(x*y/1783) +
				2*math.Cos((x+y)/137) + math.Cos((y-x)/55) + 2*math.Cos((x+8*y)/57) +
				math.Cos(math.Hypot(384-x, 274-y/9)/51)))
			dst++
		}
	}
}

func (c Bifilter) Distort(r *image.RGBA, x1 int, y1 int, x2 int, y2 int) {
	dst := 0
	src1 := y1*640 + x1
	src2 := y2*640 + x2
	for j := 0; j < 200; j++ {
		for i := 0; i < 320; i++ {
			// get distorted coordinates, use the integer part of the distortion
			// buffers and truncate to closest texel
			dY := j + int(c.displaceY[src1]>>3)
			dX := i + int(c.displaceX[src2]>>3)
			// check the texel is valid
			if (dY >= 0) && (dY < 199) && (dX >= 0) && (dX < 319) {
				// copy it to the screen
				r.Pix[dst+0] = c.bg.Pix[(dY*320+dX)*4+0]
				r.Pix[dst+1] = c.bg.Pix[(dY*320+dX)*4+1]
				r.Pix[dst+2] = c.bg.Pix[(dY*320+dX)*4+2]
				r.Pix[dst+3] = c.bg.Pix[(dY*320+dX)*4+3]
			} else {
				// otherwise, just set it to black
				r.Pix[dst+0] = 0
				r.Pix[dst+1] = 0
				r.Pix[dst+2] = 0
				r.Pix[dst+3] = 0
			}
			// next pixel
			dst += 4
			src1++
			src2++
		}
		// next line
		src1 += 320
		src2 += 320
	}
}

func (c *Bifilter) DistortBili(r *image.RGBA, x1 int, y1 int, x2 int, y2 int) {
	dst := 0
	src1 := y1*640 + x1
	src2 := y2*640 + x2
	for j := 0; j < 200; j++ {
		// for all pixels
		for i := 0; i < 320; i++ {
			// get distorted coordinates, by using the truncated integer part
			// of the distortion coefficients
			dY := j + int(c.displaceY[src1]>>3)
			dX := i + int(c.displaceX[src2]>>3)
			// get the linear interpolation coefficiants by using the fractionnal
			// part of the distortion coefficients
			cY := c.displaceY[src1] & 7
			cX := c.displaceX[src2] & 7
			// check if the texel is valid
			if (dY >= 0) && (dY < 199) && (dX >= 0) && (dX < 319) {
				// load the 4 surrounding texels and multiply them by the
				// right bilinear coefficients, then get rid of the fractionnal
				// part by shifting right by 6
				for rgb := 0; rgb < 4; rgb++ {
					o1 := (dY*320 + dX) * 4
					o2 := (dY*320 + dX + 1) * 4
					o3 := (dY*320 + dX + 320) * 4
					o4 := (dY*320 + dX + 321) * 4
					r.Pix[dst+rgb] = uint8(
						utils.ConstrainU32((
							(uint32(c.bg.Pix[o1+rgb])*uint32((8-cX)*(8-cY)))+
								(uint32(c.bg.Pix[o2+rgb])*uint32(cX*(8-cY)))+
								(uint32(c.bg.Pix[o3+rgb])*uint32((8-cX)*cY))+
								(uint32(c.bg.Pix[o4+rgb])*uint32(cX*cY)))>>6,
							0, 255))
				}
			} else {
				// otherwise, just make it black
				r.Pix[dst+0] = 0
				r.Pix[dst+1] = 0
				r.Pix[dst+2] = 0
				r.Pix[dst+3] = 0
			}
			dst += 4
			src1++
			src2++
		}
		// next line
		src1 += 320
		src2 += 320
	}
}
