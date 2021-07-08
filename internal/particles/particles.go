// Package particles
// Tutorial #11 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_11_Particle_Systems.shtml
package particles

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

type Particles struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	scaleX       [320]int
	scaleY       [200]int
	baseDist     float64
	pts          [4096]Vector
	page1        [64000]byte
	page2        [64000]byte
	colorTable   [256][3]byte
}

func (r *Particles) Draw(buffer *image.RGBA) {
	currentTime := float64(r.frameCount)
	sx := int(160.0 - 80.0*math.Sin(currentTime/559.0))
	sy := int(100.0 + 50.0*math.Sin(currentTime/611.0))
	for i := 0; i < 320; i++ {
		r.scaleX[i] = int(float64(sx) + float64(i-sx)*0.85)
	}
	for i := 0; i < 200; i++ {
		r.scaleY[i] = int(float64(sy) + float64(i-sy)*0.85)
	}
	// rescale the image
	r.rescale()
	// blur it
	r.blur()
	// setup the position of the object
	r.baseDist = 256.0 + 64.0*math.Sin(currentTime/327.0)
	rx := rotX(2.0 * math.Pi * math.Sin(currentTime/289.0))
	ry := rotY(2.0 * math.Pi * math.Cos(currentTime/307.0))
	rz := rotZ(-2.0 * math.Pi * math.Sin(currentTime/251.0))
	obj := rx.MulMat(ry).MulMat(rz)
	// draw the particles
	for i := 0; i < len(r.pts); i++ {
		r.drawSingle(obj.MulVec(r.pts[i]))
	}
	// calculate a new colour
	k := float64(len(r.pts))
	cr := int(128.0 + 127.0*math.Cos(k*math.Pi/256.0+currentTime/74.0))
	cg := int(128.0 + 127.0*math.Cos(k*math.Pi/256.0+currentTime/63.0))
	cb := int(128.0 + 127.0*math.Cos(k*math.Pi/256.0+currentTime/81.0))
	// fade palette from black (0) to colour (128) to white (255)
	for i := 0; i < 128; i++ {
		r.colorTable[i][0] = byte(cr * i / 128)
		r.colorTable[i][1] = byte(cg * i / 128)
		r.colorTable[i][2] = byte(cb * i / 128)
		r.colorTable[128+i][0] = byte(cr + (255-cr)*i/128)
		r.colorTable[128+i][1] = byte(cg + (255-cg)*i/128)
		r.colorTable[128+i][2] = byte(cb + (255-cb)*i/128)
	}
	for i, index := range r.page2 {
		buffer.Pix[i*4+0] = r.colorTable[index][0]
		buffer.Pix[i*4+1] = r.colorTable[index][1]
		buffer.Pix[i*4+2] = r.colorTable[index][2]
		//buffer.Pix[i*4+3] = 255
	}
	r.frameCount++
}

func (r *Particles) Setup() (int, int, int) {
	r.screenWidth = 320
	r.screenHeight = 200

	for i := 0; i < 256; i++ {
		r.colorTable[i][0] = byte(i)
		r.colorTable[i][1] = byte(i)
		r.colorTable[i][2] = byte(i)
	}

	for i := 0; i < len(r.pts); i++ {
		rx := rotX(2.0 * math.Pi * math.Sin(float64(i)/203.0))
		ry := rotY(2.0 * math.Pi * math.Cos(float64(i)/157.0))
		rz := rotZ(-2.0 * math.Pi * math.Cos(float64(i)/181.0))
		v := NewVector(64.0+16.0*math.Sin(float64(i)/191.0), 0, 0)
		r.pts[i] = rx.MulMat(ry).MulMat(rz).MulVec(v)
	}
	return r.screenWidth, r.screenHeight, 2
}

func (r *Particles) rescale() {
	offs := 0
	for j := 0; j < 200; j++ {
		for i := 0; i < 320; i++ {
			// get value from pixel in scaled image, and store
			index := r.scaleY[j]*320 + r.scaleX[i]
			r.page1[offs] = r.page2[index]
			offs++
		}
	}
}

func (r *Particles) blur() {
	offs := 320
	for j := 1; j < 199; j++ {
		// set first pixel of the line to 0
		r.page2[offs] = 0
		offs++
		// calculate the filter for all the other pixels
		for i := 1; i < 319; i++ {
			// calculate the average
			b := (
				int(r.page1[offs-321]) + int(r.page1[offs-320]) + int(r.page1[offs-319]) +
					int(r.page1[offs-1]) + int(r.page1[offs+1]) +
					int(r.page1[offs+319]) + int(r.page1[offs+320]) + int(r.page1[offs+321])) >> 3
			if b > 16 {
				b -= 16
			} else {
				b = 0
			}
			// store the pixel
			r.page2[offs] = byte(b)
			offs++
		}
		// set last pixel of the line to 0
		r.page2[offs] = 0
		offs++
	}
}

func (r *Particles) drawSingle(v Vector) {
	// calculate the screen coordinates of the particle
	iz := 1.0 / (v.v[2] + r.baseDist)
	x := 160.0 + 160.0*v.v[0]*iz
	y := 100.0 + 160.0*v.v[1]*iz
	// clipping
	if (x < 0) || (x > 319) || (y < 0) || (y > 199) {
		return
	}
	// convert to fixed point to help antialiasing
	sx := int(x * 8.0)
	sy := int(y * 8.0)
	// compute offset
	offs := (sy>>3)*320 + (sx >> 3)
	sx = sx & 0x7
	sy = sy & 0x7
	// add antialias particle to buffer, check for overflow
	r.page2[offs] = byte(utils.ConstrainI(int(r.page2[offs])+(7-sx)*(7-sy), 0, 255))
	r.page2[offs+1] = byte(utils.ConstrainI(int(r.page2[offs+1])+(7-sy)*sx, 0, 255))
	r.page2[offs+320] = byte(utils.ConstrainI(int(r.page2[offs+320])+sy*(7-sx), 0, 255))
	r.page2[offs+321] = byte(utils.ConstrainI(int(r.page2[offs+321])+sy*sx, 0, 255))
}
