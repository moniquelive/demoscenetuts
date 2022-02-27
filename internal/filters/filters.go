// Package filters
// Tutorial #5 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_05_Filters.shtml
package filters

import (
	_ "embed"
	"image"
	"image/draw"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed filter.png
var bgBytes []byte

type Filter struct {
	screenWidth  int
	screenHeight int
	bg           *image.RGBA
	filter       [][]int
	fps          int
}

func (c *Filter) Draw(screen *image.RGBA) {
	lf := len(c.filter)
	hlf := lf / 2

	defer draw.Draw(screen,
		image.Rect(hlf, hlf, screen.Rect.Dx()-hlf, screen.Rect.Dy()-hlf),
		c.bg, image.Point{X: hlf, Y: hlf}, draw.Src)
	c.fps++
	if c.fps%2 != 0 {
		return
	}
	width := screen.Rect.Dx()
	height := screen.Rect.Dy()

	for x := 0; x < width/2; x++ {
		i := ((height-hlf-1)*width + utils.Between(hlf, width-hlf)) * 4
		c.bg.Pix[i+0] = 255
		c.bg.Pix[i+1] = 255
		c.bg.Pix[i+2] = 255
	}

	for y := hlf; y < height-hlf; y++ {
		for x := hlf; x < width-hlf; x++ {
			acc := []int{0, 0, 0}
			for v := 0; v < lf; v++ {
				for u := 0; u < lf; u++ {
					dx := x + u - hlf
					dy := y + v - hlf
					i := (dy*width + dx) * 4
					acc[0] += int(c.bg.Pix[i+0]) * c.filter[v][u]
					acc[1] += int(c.bg.Pix[i+1]) * c.filter[v][u]
					acc[2] += int(c.bg.Pix[i+2]) * c.filter[v][u]
				}
			}
			i := (y*width + x) * 4
			c.bg.Pix[i+0] = uint8(utils.ConstrainI(acc[0]>>3, 0, 255))
			c.bg.Pix[i+1] = uint8(utils.ConstrainI(acc[1]>>3, 0, 255))
			c.bg.Pix[i+2] = uint8(utils.ConstrainI(acc[2]>>3, 0, 255))
		}
	}
}

func (c *Filter) Setup() (int, int, int) {
	c.bg = utils.LoadBufferRGBA(bgBytes)
	c.screenWidth = c.bg.Bounds().Dx()
	c.screenHeight = c.bg.Bounds().Dy()
	c.createFilters()
	return c.screenWidth, c.screenHeight, 2
}

func (c *Filter) createFilters() {
	c.filter = [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 1, 1, 1, 0},
		{0, 1, 1, 1, 0},
	}
}
