// Package bump
// Tutorial #7 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_07_Bump_Mapping.shtml
package bump

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"
	"math/rand"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed map.png
var mapBytes []byte

//go:embed bump.png
var bumpBytes []byte

type Bump struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	cm           *image.RGBA
	bm           *image.Gray
	light        [65536]byte
}

func (b *Bump) Draw(buffer *image.RGBA) {
	x1 := int(128.0*math.Cos(float64(b.frameCount)/64.0)) - 20
	y1 := int(128.0*math.Sin(float64(-b.frameCount)/45.0)) + 20
	x2 := int(128.0*math.Cos(float64(-b.frameCount)/51.0)) - 20
	y2 := int(128.0*math.Sin(float64(b.frameCount)/71.0)) + 20
	z := 192 + int(127.0*math.Sin(float64(b.frameCount)/112.0))
	b.updateBump(buffer, x1, y1, x2, y2, z)
	b.frameCount++
}

func (b *Bump) Setup() (int, int, int) {
	b.cm = utils.LoadBufferRGBA(mapBytes)
	b.bm = utils.LoadBufferGray(bumpBytes)
	b.computeLight()
	bgw := b.cm.Bounds().Dx()
	bgh := b.cm.Bounds().Dy()

	b.screenWidth = bgw
	b.screenHeight = bgh
	return b.screenWidth, b.screenHeight, 2
}

func (b *Bump) computeLight() {
	const LightSize = 2.4
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			// get the distance from the centre
			dist := float64((128-x)*(128-x) + (128-y)*(128-y))
			if math.Abs(dist) > 1 {
				dist = math.Sqrt(dist)
			}
			// then fade if according to the distance, and a random coefficient
			c := int(LightSize*dist) + (rand.Int() & 7) - 3
			// clip it
			c = utils.ConstrainI(c, 0, 255)
			// and store it
			b.light[(y<<8)+x] = byte(255 - c)
		}
	}
}

func (b *Bump) updateBump(r *image.RGBA, lx1, ly1, lx2, ly2, zoom int) {
	// we skip the first line since there are no pixels above
	// to calculate the slope with
	offs := 320 * 4
	// loop for all the other lines
	for j := 1; j < 200; j++ {
		// likewise, skip first pixel since there are no pixels on the left
		r.Pix[offs+0] = 0
		r.Pix[offs+1] = 0
		r.Pix[offs+2] = 0
		r.Pix[offs+3] = 0
		offs += 4
		for i := 1; i < 320; i++ {
			// calculate coordinates of the pixel we need in light map
			// given the slope at this point, and the zoom coefficient
			px := (i * zoom >> 8) + int(b.bm.Pix[offs>>2-1]) - int(b.bm.Pix[offs>>2])
			py := (j * zoom >> 8) + int(b.bm.Pix[offs>>2-320]) - int(b.bm.Pix[offs>>2])
			// add the movement of the first light
			x := px + lx1
			y := py + ly1
			// check if the coordinates are inside the light buffer
			c := 0
			if (y >= 0) && (y < 256) && (x >= 0) && (x < 256) {
				// if so get the pixel
				c = int(b.light[(y<<8)+x])
			}
			// now do the same for the second light
			x = px + lx2
			y = py + ly2
			// this time we add the light's intensity to the first value
			if (y >= 0) && (y < 256) && (x >= 0) && (x < 256) {
				c += int(b.light[(y<<8)+x])
			}
			// make sure it's not too big
			c = utils.ConstrainI(c, 0, 255)
			// look up the colour multiplied by the light coeficient
			r.Pix[offs+0] = byte(float64(b.cm.Pix[offs+0]) * float64(c) / 256.0)
			r.Pix[offs+1] = byte(float64(b.cm.Pix[offs+1]) * float64(c) / 256.0)
			r.Pix[offs+2] = byte(float64(b.cm.Pix[offs+2]) * float64(c) / 256.0)
			r.Pix[offs+3] = byte(float64(b.cm.Pix[offs+3]) * float64(c) / 256.0)
			offs += 4
		}
	}
}
