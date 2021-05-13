package cyber1

import (
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed cyber1.png
var bgBytes []byte

type Lerp struct {
	t      float64
	pixels []Pixel
	bg     *image.RGBA
}

type Pixel struct {
	xDest, yDest int
	xCurr, yCurr int
	color        color.RGBA
}

func (c *Lerp) Draw(buffer *image.RGBA) {
	w := c.bg.Bounds().Dx()

	c.t = utils.Constrain(c.t+0.002, 0, 1)
	for _, p := range c.pixels {
		p.xCurr = int(utils.Lerp(float64(p.xCurr), float64(p.xDest), c.t))
		p.yCurr = int(utils.Lerp(float64(p.yCurr), float64(p.yDest), c.t))

		loc := p.yCurr*w + p.xCurr
		buffer.Pix[4*loc+0] = p.color.R
		buffer.Pix[4*loc+1] = p.color.G
		buffer.Pix[4*loc+2] = p.color.B
		buffer.Pix[4*loc+3] = p.color.A
	}
}

func (c *Lerp) Setup() (int, int, int) {
	c.bg = utils.LoadBufferRGBA(bgBytes)
	c.readPixels()

	return c.bg.Bounds().Dx(), c.bg.Bounds().Dy(), 2
}

func (c *Lerp) readPixels() {
	w := c.bg.Bounds().Dx()
	h := c.bg.Bounds().Dy()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			loc := (y*w + x) * 4
			c.pixels = append(c.pixels, Pixel{
				xDest: x,
				yDest: y,
				xCurr: rand.Intn(w),
				yCurr: rand.Intn(h),
				color: color.RGBA{
					R: c.bg.Pix[loc+0],
					G: c.bg.Pix[loc+1],
					B: c.bg.Pix[loc+2],
					A: c.bg.Pix[loc+3],
				},
			})
		}
	}
}
