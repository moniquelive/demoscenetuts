// Package plane
// Tutorial #14 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_14_Perspective_Correct_Texture_Mapping.shtml
package plane

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed texture.png
var textureBytes []byte

type Plane struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	texture      *image.Paletted
	A, B, C      utils.Vector
}

func (p *Plane) Draw(buffer *image.RGBA) {
	currentTime := float64(p.frameCount) * 2e4

	p.A = utils.NewVector(currentTime/34984.0, -16, currentTime/43512.0)
	p.B = utils.RotY(0.32).MulVec(utils.NewVector(256, 0, 0))
	p.C = utils.RotY(0.32).MulVec(utils.NewVector(0, 0, 256))

	p.drawPlane(buffer)

	p.frameCount++
}

func (p *Plane) drawPlane(buffer *image.RGBA) {
	Cx := p.B[1]*p.C[2] - p.C[1]*p.B[2]
	Cy := p.C[0]*p.B[2] - p.B[0]*p.C[2]
	// the 240 represents the distance of the projection plane
	// change it to modify the field of view
	Cz := (p.B[0]*p.C[1] - p.C[0]*p.B[1]) * 240
	Ax := p.C[1]*p.A[2] - p.A[1]*p.C[2]
	Ay := p.A[0]*p.C[2] - p.C[0]*p.A[2]
	Az := (p.C[0]*p.A[1] - p.A[0]*p.C[1]) * 240
	Bx := p.A[1]*p.B[2] - p.B[1]*p.A[2]
	By := p.B[0]*p.A[2] - p.A[0]*p.B[2]
	Bz := (p.A[0]*p.B[1] - p.B[0]*p.A[1]) * 240
	// only render the lower part of the plane, looks ugly above
	offs := 105 * 320
	for j := float64(105); j < 200; j++ {
		// compute the (U,V) coordinates for the start of the line
		a := Az + Ay*(j-100) + Ax*-161
		b := Bz + By*(j-100) + Bx*-161
		c := Cz + Cy*(j-100) + Cx*-161
		// quick distance check, if it's too far reduce it
		var ic float64
		if math.Abs(c) > 65536 {
			ic = 1 / c
		} else {
			ic = 1 / 65536
		}
		// compute original (U,V)
		u := int(a * 16777216 * ic)
		v := int(b * 16777216 * ic)
		// and the deltas we need to interpolate along this line
		du := int(16777216 * Ax * ic)
		dv := int(16777216 * Bx * ic)
		// start the loop
		for i := 0; i < 320; i++ {
			palIndex := p.texture.Pix[((v>>8)&0xff00)+((u>>16)&0xff)]
			pr, pg, pb, _ := p.texture.Palette[palIndex].RGBA()
			buffer.Pix[offs*4+0] = byte(pr)
			buffer.Pix[offs*4+1] = byte(pg)
			buffer.Pix[offs*4+2] = byte(pb)
			offs++
			// interpolate
			u += du
			v += dv
		}
	}
}

func (p *Plane) Setup() (int, int, int) {
	p.texture = utils.LoadBufferPaletted(textureBytes)
	p.screenWidth = 320
	p.screenHeight = 200

	return p.screenWidth, p.screenHeight, 2
}
