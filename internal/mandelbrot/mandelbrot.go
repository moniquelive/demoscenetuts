// Package mandelbrot
// Tutorial #8 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_08_Fractal_Zooming.shtml
package mandelbrot

import (
	_ "embed"
	"image"
	"image/draw"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

type Mandelbrot struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	cm           *image.RGBA
	d            complex128
	p            complex128
	s            complex128
	z            complex128
	offs         int
	zoom         bool
}

const o = complex(-0.577816-9.31323E-10-1.16415E-10, -0.631121-2.38419E-07+1.49012E-08)
const resx = 320.0
const resy = 200.0

func (b *Mandelbrot) Draw(buffer *image.RGBA) {

	if b.frameCount%3 == 0 {
		b.startFrac(o-b.z, o+b.z)
		for i := 0; i < resy; i++ {
			b.computeFrac(b.cm)
		}
	}

	if b.frameCount%250 == 0 {
		b.zoom = !b.zoom
	}
	if b.zoom {
		b.z *= 0.9
	} else {
		b.z *= 1.1
	}
	b.frameCount++
	draw.Draw(buffer, buffer.Bounds(), b.cm, image.Point{}, draw.Src)
}

func (b *Mandelbrot) Setup() (int, int, int) {
	b.cm = image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 320, Y: 200},
	})
	b.screenWidth = 320
	b.screenHeight = 200
	b.z = complex(4, 4)
	return b.screenWidth, b.screenHeight, 2
}

func (b *Mandelbrot) startFrac(_s, e complex128) {
	// compute deltas for interpolation in complex plane
	b.d = complex(
		(real(e)-real(_s))/resx,
		(imag(e)-imag(_s))/resy)
	// remember start values
	b.p = _s
	b.s = _s
	b.offs = 0
}

func (b Mandelbrot) pal(i byte) [3]byte {
	fi := float64(i)
	ffc := float64(b.frameCount)
	return [3]byte{
		byte(utils.Constrain(127.0-128.0*math.Cos(fi*math.Pi/128.0+ffc*0.0041), 0, 255)),
		byte(utils.Constrain(127.0-128.0*math.Cos(fi*math.Pi/128.0+ffc*0.00141), 0, 255)),
		byte(utils.Constrain(127.0-128.0*math.Cos(fi*math.Pi/64.0+ffc*0.00136), 0, 255)),
	}
}

func (b *Mandelbrot) computeFrac(buffer *image.RGBA) {
	b.p = complex(
		real(b.s),
		imag(b.p))
	for i := 0; i < resx; i++ {
		c := byte(0)
		v := b.p
		var nv complex128
		// loop until distance is above 2, or counter hits limit
		for (real(v)*real(v)+imag(v)*imag(v) < 4) && (c < 255) {
			// Z(0) = C
			// Z(n+1) = Z(n)^2 + C      (1)

			// compute Z(n+1) given Z(n)
			nv = complex(
				real(v)*real(v)-imag(v)*imag(v)+real(b.p),
				2*real(v)*imag(v)+imag(b.p))

			v = nv // that becomes Z(n)
			c++    // increment counter
		}
		// store colour
		p := b.pal(c)
		buffer.Pix[b.offs*4+0] = p[0]
		buffer.Pix[b.offs*4+1] = p[1]
		buffer.Pix[b.offs*4+2] = p[2]
		buffer.Pix[b.offs*4+3] = 255
		b.offs++
		// interpolate X
		b.p = complex(
			real(b.p)+real(b.d),
			imag(b.p))
	}
	// interpolate Y
	b.p = complex(
		real(b.p),
		imag(b.p)+imag(b.d))
}
