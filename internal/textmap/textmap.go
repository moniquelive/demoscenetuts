// Package textmap
// Tutorial #9 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_09_Static_Texture_Mapping.shtml
package textmap

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed texture.png
var textureBytes []byte

type Textmap struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	texture      *image.Paletted
	texcoord     [128000]byte
}

func (b *Textmap) Draw(buffer *image.RGBA) {
	const vel = 1
	du := byte(b.frameCount >> vel)
	dv := byte(b.frameCount >> (vel - 1))
	doffs := 0
	soffs := 0
	for i := 0; i < 320*200; i++ {
		u := uint(b.texcoord[soffs] + du)
		v := uint(b.texcoord[soffs+1] + dv)
		palIndex := b.texture.Pix[((v << 8) + u)]
		cr, cg, cb, ca := b.texture.Palette[palIndex].RGBA()
		buffer.Pix[doffs*4+0] = byte(cr)
		buffer.Pix[doffs*4+1] = byte(cg)
		buffer.Pix[doffs*4+2] = byte(cb)
		buffer.Pix[doffs*4+3] = byte(ca)
		doffs++
		soffs += 2
	}
	b.frameCount++
}

func (b *Textmap) Setup() (int, int, int) {
	b.texture = utils.LoadBufferPaletted(textureBytes)
	b.screenWidth = 320
	b.screenHeight = 200
	b.initHole()
	return b.screenWidth, b.screenHeight, 2
}

func (b *Textmap) initHole() {
	// alloc memory to store 320*200 times u, v
	offs := 0
	// precalc the (u,v) coordinates
	for j := -100; j < 100; j++ {
		for i := -160; i < 160; i++ {
			// get coordinates of ray that projects through this pixel
			dx := float64(i) / 200.0
			dy := float64(-j) / 200.0
			dz := float64(1)
			// normalize them
			d := 20.0 / math.Sqrt(dx*dx+dy*dy+1)
			dx *= d
			dy *= d
			dz *= d
			// start interpolation at origin
			var x float64
			var y float64
			var z float64
			// set original precision
			d = 16
			// interpolate along ray
			for d > 0 {
				// continue until we hit a wall
				for ((x-getXPos(z))*(x-getXPos(z))+(y-getYPos(z))*(y-getYPos(z)) < getRadius(z)) && (z < 1024) {
					x += dx
					y += dy
					z += dz
				}
				// reduce precision and reverse direction
				x -= dx
				y -= dy
				z -= dz
				dx /= 2.0
				dy /= 2.0
				dz /= 2.0
				d -= 1.0
			}
			// calculate the texture coordinates
			x -= getXPos(z)
			y -= getYPos(z)
			ang := math.Atan2(y, x) * 256.0 / math.Pi
			u := byte(ang)
			v := byte(z)
			// store texture coordinates
			b.texcoord[offs] = u
			b.texcoord[offs+1] = v
			offs += 2
		}
	}
}

func getXPos(f float64) float64 {
	return -16.0 * math.Sin(f*math.Pi/256.0)
}

func getYPos(f float64) float64 {
	return -16.0 * math.Sin(f*math.Pi/256.0)
}

func getRadius(_ float64) float64 {
	return 128
}
