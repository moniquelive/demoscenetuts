// Package span
// Tutorial #12 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_12_Span_Based_Rendering.shtml
package span

import (
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

//go:embed texture.png
var textureBytes []byte

const (
	SPANS    = 32
	SPANMASK = 31
)

type XZ struct {
	x, z float64
}

type Span struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	texture      *image.Paletted
	texcoord     [128000]byte
	light        [256]byte
	lut          [65536]byte
	zbuffer      [64000]uint16
	pts          [200][SPANS]XZ
	nrm          [200][SPANS]XZ
}

func memset(a *[64000]uint16, v uint16) {
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func (s *Span) Draw(buffer *image.RGBA) {
	// clear the zbuffer
	memset(&s.zbuffer, 0xffff)
	// draw all slices
	fc := float64(s.frameCount) * 1e4
	for i := 0; i < 200; i++ {
		var xc [SPANS]int
		var zc [SPANS]int
		var nc [SPANS]int
		// set a different rotation angle for each slice
		angle := fc/1948613 +
			2*math.Pi*(math.Cos(fc/3179640)*
				math.Cos(-float64(i)/193)*math.Sin(fc/2714147)+1)
		// compute the sine and cosine in 8.8 fxd point
		ca := 256 * math.Cos(angle)
		sa := 256 * math.Sin(angle)
		// and determine the horizontal position of this slice
		xoffs := int(40960 + 10240*math.Cos(fc/294517/*4*/)*
			math.Sin(float64(i)/59))
		for j := 0; j < SPANS; j++ {
			// rotate the point and get X coordinate
			xc[j] = xoffs + int(s.pts[i][j].x*ca+s.pts[i][j].z*-sa)
			// rotate the point and get Z coordinate
			zc[j] = int(s.pts[i][j].x*sa + s.pts[i][j].z*ca)
			// now get the coordinate in the lightmap by computing
			// the X component of the normal
			tmp := int(s.nrm[i][j].x*ca + s.nrm[i][j].z*-sa)
			tmp = 256 - tmp
			nc[j] = (128 << 16) + (tmp << 15)
		}
		// now just draw all spans
		for j := 0; j < SPANS; j++ {
			s.drawSpan(buffer, i,
				xc[j]>>8, j<<20, nc[j], 49152+zc[j],
				xc[(j+1)&SPANMASK]>>8, (j+1)<<20, nc[(j+1)&SPANMASK], 49152+zc[(j+1)&SPANMASK])
		}
	}
	s.frameCount++
}

func (s *Span) drawSpan(buffer *image.RGBA, y, x1, tx1, px1, z1, x2, tx2, px2, z2 int) {
	textOffs := s.frameCount / 31456
	// quick check, if facing back then return
	if x1 >= x2 {
		return
	}
	// compute deltas for interpolation
	dx := x2 - x1
	dtx := (tx2 - tx1) / dx // assume 16.16 fixed point
	dpx := (px2 - px1) / dx // 16.16
	dz := (z2 - z1) / dx    // doesn't matter, whatever parameters are is fine
	// get destination offset in buffer
	offs := y*320 + x1
	// loop for all pixels concerned
	for i := 0; i < dx; i++ {
		// check z buffer
		if z1 < int(s.zbuffer[offs]) {
			// if visible load the texel from the translated texture
			texel := int(s.texture.Pix[(((y+textOffs)<<8)&0xff00)+((tx1>>16)&0xff)])
			// and the texel from the light map
			lumel := int(s.light[(px1>>16)&0xff])
			// mix them together
			palIndex := s.lut[(texel<<8)+lumel]
			pr, pg, pb, _ := s.texture.Palette[palIndex].RGBA()
			buffer.Pix[offs*4+0] = byte(pr)
			buffer.Pix[offs*4+1] = byte(pg)
			buffer.Pix[offs*4+2] = byte(pb)
			//buffer.Pix[offs*4+3] = byte(pa)
			// and update the zbuffer
			s.zbuffer[offs] = uint16(z1)
		}
		// interpolate our values
		px1 += dpx
		tx1 += dtx
		z1 += dz
		// and find next pixel
		offs++
	}
}
func (s *Span) Setup() (int, int, int) {
	s.texture = utils.LoadBufferPaletted(textureBytes)
	s.screenWidth = 320
	s.screenHeight = 200
	// prepare the lighting
	for i := 0; i < 256; i++ {
		// calculate distance from the centre
		c := math.Abs(float64(255 - i*2))
		// check for overflow
		if c > 255 {
			c = 255
		}
		// store lumel
		s.light[i] = byte(255 - c)
	}
	s.computeLUT(s.texture.Palette)
	// prepare 3D data
	// store points as new data type, the type vector would take up more
	// space than necessary
	// initialise each point of each slice
	for i := 0; i < SPANS; i++ {
		angle := float64(i) * math.Pi * 2 / SPANS
		ca := math.Cos(angle)
		sa := math.Sin(angle)
		for j := 0; j < 200; j++ {
			// use these equations to get a weird shaped flubber
			//        float radx = 48 + 32 * sin( ((float)i+47*cos((float)j*0.123)) * 0.019 )
			//                        * cos( ((float)j+39*sin((float)i*0.137)) * 0.023 ),
			//              rady = 48 + 32 * cos( ((float)i+31*sin((float)j*0.147)) * 0.027 )
			//                         * sin( ((float)j+43*cos((float)i*0.111)) * 0.014 );
			// this generates a flat cylinder
			radx := 32.0
			rady := 96.0
			// store the points
			s.pts[j][i].x = radx * ca
			s.pts[j][i].z = rady * sa
		}
	}
	// now calculate the normals
	for i := 0; i < SPANS; i++ {
		for j := 0; j < 200; j++ {
			// get coords of tangeant vector
			nx := s.pts[j][(i+1)&SPANMASK].x - s.pts[j][(i-1)&SPANMASK].x
			nz := s.pts[j][(i+1)&SPANMASK].z - s.pts[j][(i-1)&SPANMASK].z
			// and compute it's length
			inorm := 1 / math.Sqrt(nx*nx+nz*nz)
			// now make sure the vector has length one
			nx *= inorm
			nz *= inorm
			// and store it
			s.nrm[j][i].x = nx
			s.nrm[j][i].z = nz
		}
	}
	return s.screenWidth, s.screenHeight, 2
}

func (s *Span) computeLUT(pal color.Palette) {
	// for each colour
	for j := uint32(0); j < 256; j++ {
		r, g, b, _ := pal[j].RGBA()
		r >>= 8
		g >>= 8
		b >>= 8
		// calculate shades from a third of the colour to the colour
		for i := uint32(0); i < 240; i++ {
			// get the index of the best colour
			s.lut[(j<<8)+i] = s.closestColour(
				pal, r*(120+i)/360, g*(120+i)/360, b*(120+i)/360)
		}
		// calc shades half way from the colour to white
		for i := uint32(0); i < 16; i++ {
			s.lut[(j<<8)+240+i] = s.closestColour(
				pal, r+(255-r)*i/31, g+(255-g)*i/31, b+(255-b)*i/31)
		}
	}
}

func (s *Span) closestColour(pal color.Palette, r, g, b uint32) byte {
	var index byte

	dist := uint32(1 << 30)
	// for all colours loop
	for i := uint32(0); i < 256; i++ {
		// calculate how far the current colour is from the required colour
		pr, pg, pb, _ := pal[i].RGBA()
		pr >>= 8
		pg >>= 8
		pb >>= 8
		newDist := (r-pr)*(r-pr) + (g-pg)*(g-pg) + (b-pb)*(b-pb)
		// if it's the same then we've found a match, so return it
		if newDist == 0 {
			return byte(i)
		}
		// if the distance is closer, then remember this colour's index
		if newDist < dist {
			index = byte(i)
			dist = newDist
		}
	}
	// return the one we found closest
	return index
}
