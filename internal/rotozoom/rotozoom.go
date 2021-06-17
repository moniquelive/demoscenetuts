// Package rotozoom
// Tutorial #10 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_10_Roto-Zooming.shtml
package rotozoom

import (
	_ "embed"
	"image"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed texture.png
var textureBytes []byte

type Rotozoom struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	texture      *image.Paletted
}

func (r *Rotozoom) Draw(buffer *image.RGBA) {
	const f = 4e4
	t := float64(r.frameCount)
	xRot := t / 21547814.0 * f
	yRot := t / 18158347.0 * f
	zoomF := t / 4210470.0 * f
	r.DoRotoZoom(
		buffer,
		2048*math.Sin(xRot),         // X centre coord
		2048*math.Cos(yRot),         // Y centre coord
		256.0+192.0*math.Cos(zoomF), // zoom coef
		t/731287.0*f,                // angle
	)
	r.frameCount++
}

func (r *Rotozoom) Setup() (int, int, int) {
	r.texture = utils.LoadBufferPaletted(textureBytes)
	r.screenWidth = 320
	r.screenHeight = 200
	return r.screenWidth, r.screenHeight, 2
}

func (r *Rotozoom) DoRotoZoom(buffer *image.RGBA, cx, cy, radius, angle float64) {
	x1 := (int)(65536.0 * (cx + radius*math.Cos(angle)))
	y1 := (int)(65536.0 * (cy + radius*math.Sin(angle)))
	x2 := (int)(65536.0 * (cx + radius*math.Cos(angle+2.02458)))
	y2 := (int)(65536.0 * (cy + radius*math.Sin(angle+2.02458)))
	x3 := (int)(65536.0 * (cx + radius*math.Cos(angle-1.11701)))
	y3 := (int)(65536.0 * (cy + radius*math.Sin(angle-1.11701)))
	r.BlockTextureScreen(buffer, x1, y1, x2, y2, x3, y3)
}

func (r *Rotozoom) BlockTextureScreen(buffer *image.RGBA, Ax, Ay, Bx, By, Cx, Cy int) {
	// compute global deltas, 40 blocks along X, 25 blocks along Y
	dxdx := (Bx - Ax) / 40
	dydx := (By - Ay) / 40
	dxdy := (Cx - Ax) / 25
	dydy := (Cy - Ay) / 25
	// compute internal block deltas
	dxbdx := (Bx - Ax) / 320
	dybdx := (By - Ay) / 320
	dxbdy := (Cx - Ax) / 200
	dybdy := (Cy - Ay) / 200
	var baseoffs, offs int
	// for all blocks along Y
	for j := 0; j < 25; j++ {
		Cx = Ax
		Cy = Ay
		// for all blocks along X
		for i := 0; i < 40; i++ {
			offs = baseoffs
			// set original position
			Ex := Cx
			Ey := Cy
			// for each line of 8 pixels in the block
			for y := 0; y < 8; y++ {
				// set original position
				Fx := Ex
				Fy := Ey
				// for each pixel in the block
				for x := 0; x < 8; x++ {
					palIndex := r.texture.Pix[((Fy>>8)&0xff00)+((Fx>>16)&0xff)]
					cr, cg, cb, ca := r.texture.Palette[palIndex].RGBA()
					buffer.Pix[offs*4+0] = byte(cr)
					buffer.Pix[offs*4+1] = byte(cg)
					buffer.Pix[offs*4+2] = byte(cb)
					buffer.Pix[offs*4+3] = byte(ca)
					// interpolate to find next texel
					Fx += dxbdx
					Fy += dybdx
					offs++
				}
				// next line in the block
				Ex += dxbdy
				Ey += dybdy
				offs += 312
			}

			// next block
			baseoffs += 8
			Cx += dxdx
			Cy += dydx
		}

		// next line of blocks
		baseoffs += 320 * 7
		Ax += dxdy
		Ay += dydy
	}
}
