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

	dr, di, pr, pi, sr, si float64
	offs                   int
	zx, zy                 float64
	zoom                   bool
}

const or = -0.577816 - 9.31323E-10 - 1.16415E-10
const oi = -0.631121 - 2.38419E-07 + 1.49012E-08
const resx = 320.0
const resy = 200.0

func (b *Mandelbrot) Draw(buffer *image.RGBA) {

	if b.frameCount%3 == 0 {
		b.startFrac(or-b.zx, oi-b.zy, or+b.zx, oi+b.zy)
		for i := 0; i < resy; i++ {
			b.computeFrac(b.cm)
		}
	}

	if b.frameCount%250 == 0 {
		b.zoom = !b.zoom
	}
	if b.zoom {
		b.zx *= 0.9
		b.zy *= 0.9
	} else {
		b.zx *= 1.1
		b.zy *= 1.1
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
	b.zx = 4.0
	b.zy = 4.0
	return b.screenWidth, b.screenHeight, 2
}

func (b *Mandelbrot) startFrac(_sr, _si, er, ei float64) {
	// compute deltas for interpolation in complex plane
	b.dr = (er - _sr) / resx
	b.di = (ei - _si) / resy
	// remember start values
	b.pr = _sr
	b.pi = _si
	b.sr = _sr
	b.si = _si
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
	b.pr = b.sr
	for i := 0; i < resx; i++ {
		c := byte(0)
		vi := b.pi
		vr := b.pr
		nvi := 0.0
		nvr := 0.0
		// loop until distance is above 2, or counter hits limit
		for (vr*vr+vi*vi < 4) && (c < 255) {
			// Z(0) = C
			// Z(n+1) = Z(n)^2 + C      (1)
			// compute Z(n+1) given Z(n)
			nvr = vr*vr - vi*vi + b.pr
			nvi = 2*vi*vr + b.pi

			// that becomes Z(n)
			vi = nvi
			vr = nvr

			// increment counter
			c++
		}
		// store colour
		p := b.pal(c)
		buffer.Pix[b.offs*4+0] = p[0]
		buffer.Pix[b.offs*4+1] = p[1]
		buffer.Pix[b.offs*4+2] = p[2]
		buffer.Pix[b.offs*4+3] = 255
		b.offs++
		// interpolate X
		b.pr += b.dr
	}
	// interpolate Y
	b.pi += b.di
}
