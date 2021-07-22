// Package polygon
// Tutorial #13 - Efeitos Demoscene
// https://www.flipcode.com/archives/The_Art_of_Demomaking-Issue_13_Polygon_Engines.shtml
package polygon

import (
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"math"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

const (
	SLICES     = 32
	SPANS      = 16
	EXT_RADIUS = 64
	INT_RADIUS = 24
)

//go:embed texture.png
var textureBytes []byte

type Poly struct {
	p              [4]int
	tx             [4]int
	ty             [4]int
	normal, centre Vector
}

type Polygon struct {
	screenWidth  int
	screenHeight int
	frameCount   int
	texture      *image.Paletted
	light        [256 * 256]byte
	zbuffer      [64000]uint16
	lut          [65536]byte
	org, cur     struct {
		normals  [SLICES * SPANS]Vector
		vertices [SLICES * SPANS]Vector
	}
	polies               [SPANS * SLICES]Poly
	objrot               Matrix
	objpos               Vector
	edge_table           [200][2]edge_data
	poly_minY, poly_maxY int
}

func (p *Polygon) Draw(buffer *image.RGBA) {
	currentTime := float64(p.frameCount) * 2e4
	// clear the zbuffer
	utils.Memset(&p.zbuffer, 0xffff)
	// clear the background
	// set the torus' rotation
	p.objrot =
		rotX(currentTime/1124548.0 + math.Pi*math.Cos(currentTime/2234117.0)).
			MulMat(rotY(currentTime/1345179.0 + math.Pi*math.Sin(currentTime/2614789.0))).
			MulMat(rotZ(currentTime/1515713.0 - math.Pi*math.Cos(currentTime/2421376.0)))
	// and it's position
	p.objpos = NewVector(
		48*math.Cos(currentTime/1266398.0),
		48*math.Sin(currentTime/1424781.0),
		200+80*math.Sin(currentTime/1912378.0))
	// rotate and project our points
	p.transformPts()
	// and draw the polygons
	p.drawPolies(buffer)
	p.frameCount++
}

func (p *Polygon) Setup() (int, int, int) {
	p.texture = utils.LoadBufferPaletted(textureBytes)
	p.screenWidth = 320
	p.screenHeight = 200

	for j := 0; j < 256; j++ {
		for i := 0; i < 256; i++ {
			// calculate distance from the centre
			c := ((128-float64(i))*(128-float64(i)) + (128-float64(j))*(128-float64(j))) / 80
			// check for overflow
			if c > 255 {
				c = 255
			}
			// store lumel
			p.light[(j<<8)+i] = byte(255 - c)
		}
	}
	p.computeLUT(p.texture.Palette)
	p.initObject()

	return p.screenWidth, p.screenHeight, 2
}

func (p *Polygon) drawPolies(buffer *image.RGBA) {
	for n := 0; n < SPANS * SLICES; n++ {
		// rotate the centre and normal of the poly to check if it is actually visible.
		ncent := p.objrot.MulVec(p.polies[n].centre)
		nnorm := p.objrot.MulVec(p.polies[n].normal)

		// calculate the dot product, and check it's sign
		if (ncent[0]+p.objpos[0])*nnorm[0]+
			(ncent[1]+p.objpos[1])*nnorm[1]+
			(ncent[2]+p.objpos[2])*nnorm[2] < 0 {
			// the polygon is visible, so setup the edge table
			p.initEdgeTable()
			// process all our edges
			for i := 0; i < 4; i++ {
				p.scanEdge(
					// the vertex in screen space
					p.cur.vertices[p.polies[n].p[i]],
					// the static texture coordinates
					p.polies[n].tx[i], p.polies[n].ty[i],
					// the dynamic text coords computed with the normals
					(int)(65536*(128+127*p.cur.normals[p.polies[n].p[i]][0])),
					(int)(65536*(128+127*p.cur.normals[p.polies[n].p[i]][1])),
					// second vertex in screen space
					p.cur.vertices[p.polies[n].p[(i+1)&3]],
					// static text coords
					p.polies[n].tx[(i+1)&3], p.polies[n].ty[(i+1)&3],
					// dynamic texture coords
					(int)(65536*(128+127*p.cur.normals[p.polies[n].p[(i+1)&3]][0])),
					(int)(65536*(128+127*p.cur.normals[p.polies[n].p[(i+1)&3]][1])),
				)
			}
			// quick clipping
			if p.poly_minY < 0 {
				p.poly_minY = 0
			}
			if p.poly_maxY > 200 {
				p.poly_maxY = 200
			}
			// do we have to draw anything?
			if (p.poly_minY < p.poly_maxY) && (p.poly_maxY > 0) && (p.poly_minY < 200) {
				// if so just draw relevant lines
				for i := p.poly_minY; i < p.poly_maxY; i++ {
					p.drawSpan(buffer, i, &p.edge_table[i][0], &p.edge_table[i][1])
				}
			}
		}
	}
}

func (p *Polygon) initObject() {
	k := 0
	for i := 0; i < SLICES; i++ {
		// find angular position
		ext_angle := float64(i) * math.Pi * 2.0 / SLICES
		ca := math.Cos(ext_angle)
		sa := math.Sin(ext_angle)
		// now loop round C2
		for j := 0; j < SPANS; j++ {
			int_angle := float64(j) * math.Pi * 2.0 / SPANS
			int_rad := EXT_RADIUS + INT_RADIUS*math.Cos(int_angle)
			// compute position of vertex by rotating it round C1
			p.org.vertices[k] = NewVector(
				int_rad*ca,
				INT_RADIUS*math.Sin(int_angle),
				int_rad*sa)
			// then find the normal, i.e. the normalised vector representing the
			// distance to the correpsonding point on C1
			p.org.normals[k] = (p.org.vertices[k].
				Sub(NewVector(EXT_RADIUS*ca, 0, EXT_RADIUS*sa))).normalize()
			k++
		}
	}

	// now initialize the polygons, there are as many quads as vertices
	// perform the same loop
	for i := 0; i < SLICES; i++ {
		for j := 0; j < SPANS; j++ {
			P := &p.polies[i*SPANS+j]

			// setup the pointers to the 4 concerned vertices
			P.p[0] = i*SPANS + j
			P.p[1] = i*SPANS + ((j + 1) % SPANS)
			P.p[3] = ((i+1)%SLICES)*SPANS + j
			P.p[2] = ((i+1)%SLICES)*SPANS + ((j + 1) % SPANS)

			// now compute the static texture refs (X)
			P.tx[0] = (i * 512 / SLICES) << 16
			P.tx[1] = (i * 512 / SLICES) << 16
			P.tx[3] = ((i + 1) * 512 / SLICES) << 16
			P.tx[2] = ((i + 1) * 512 / SLICES) << 16

			// now compute the static texture refs (Y)
			P.ty[0] = (j * 512 / SPANS) << 16
			P.ty[1] = ((j + 1) * 512 / SPANS) << 16
			P.ty[3] = (j * 512 / SPANS) << 16
			P.ty[2] = ((j + 1) * 512 / SPANS) << 16

			// get the normalized diagonals
			d1 := p.org.vertices[P.p[2]].Sub(p.org.vertices[P.p[0]]).normalize()
			d2 := p.org.vertices[P.p[3]].Sub(p.org.vertices[P.p[1]]).normalize()
			// and their dot product
			temp := NewVector(d1[1]*d2[2]-d1[2]*d2[1],
				d1[2]*d2[0]-d1[0]*d2[2],
				d1[0]*d2[1]-d1[1]*d2[0])
			// normalize that and we get the face's normal
			P.normal = temp.normalize()

			// the centre of the face is just the average of the 4 corners
			// we could use this for depth sorting
			temp = p.org.vertices[P.p[0]].
				Add(p.org.vertices[P.p[1]]).
				Add(p.org.vertices[P.p[2]]).
				Add(p.org.vertices[P.p[3]])
			P.centre = NewVector(temp[0]*0.25, temp[1]*0.25, temp[2]*0.25)
		}
	}
}

func (p *Polygon) transformPts() {
	for i := 0; i < SLICES*SPANS; i++ {
		// perform rotation
		p.cur.normals[i] = p.objrot.MulVec(p.org.normals[i])
		p.cur.vertices[i] = p.objrot.MulVec(p.org.vertices[i])
		// now project onto the screen
		p.cur.vertices[i][2] += p.objpos[2]
		p.cur.vertices[i][0] = 200*(p.cur.vertices[i][0]+p.objpos[0])/p.cur.vertices[i][2] + 160
		p.cur.vertices[i][1] = 200*(p.cur.vertices[i][1]+p.objpos[1])/p.cur.vertices[i][2] + 100
		// optionaly draw
		//     vga->PutPixel( (int)(cur.vertices[i][0]), (int)(cur.vertices[i][1]), 255 );
	}
}

func (p *Polygon) closestColour(pal color.Palette, r, g, b uint32) byte {
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

func (p *Polygon) computeLUT(pal color.Palette) {
	// for each colour
	for j := uint32(0); j < 256; j++ {
		r, g, b, _ := pal[j].RGBA()
		r >>= 8
		g >>= 8
		b >>= 8
		// calculate shades from a third of the colour to the colour
		for i := uint32(0); i < 224; i++ {
			// get the index of the best colour
			p.lut[(j<<8)+i] = p.closestColour(pal, r, g, b)
		}
		// calc shades half way from the colour to white
		for i := uint32(0); i < 32; i++ {
			p.lut[(j<<8)+224+i] = p.closestColour(
				pal, r+(255-r)*i/32, g+(255-g)*i/32, b+(255-b)*i/32)
		}
	}
}
