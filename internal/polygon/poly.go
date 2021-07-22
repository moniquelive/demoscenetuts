package polygon

import "image"

// one entry of the edge table
type edge_data struct {
	x, px, py, tx, ty, z int
}

/*
 * clears all entries in the edge table
 */
func (p *Polygon) initEdgeTable() {
	for i := 0; i < 200; i++ {
		p.edge_table[i][0].x = -1
		p.edge_table[i][1].x = -1
	}
	p.poly_minY = 200
	p.poly_maxY = -1
}

/*
 * scan along one edge of the poly, i.e. interpolate all values and store
 * in the edge table
 */
func (p *Polygon) scanEdge(p1 Vector, tx1, ty1, px1, py1 int,
	p2 Vector, tx2, ty2, px2, py2 int) {
	// we can't handle this case, so we recall the proc with reversed params
	// saves having to swap all the vars, but it's not good practice
	if p2[1] < p1[1] {
		p.scanEdge(p2, tx2, ty2, px2, py2, p1, tx1, ty1, px1, py1)
		return
	}
	// convert to fixed point
	x1 := int(p1[0] * 65536)
	y1 := int(p1[1])
	z1 := int(p1[2] * 16)
	x2 := int(p2[0] * 65536)
	y2 := int(p2[1])
	z2 := int(p2[2] * 16)

	// update the min and max of the current polygon
	if y1 < p.poly_minY {
		p.poly_minY = y1
	}
	if y2 > p.poly_maxY {
		p.poly_maxY = y2
	}
	// compute deltas for interpolation
	dy := y2 - y1
	if dy == 0 {
		return
	}
	dx := (x2 - x1) / dy
	dtx := (tx2 - tx1) / dy
	dty := (ty2 - ty1) / dy
	dpx := (px2 - px1) / dy
	dpy := (py2 - py1) / dy
	dz := (z2 - z1) / dy
	// interpolate along the edge
	for y := y1; y < y2; y++ {
		// don't go out of the screen
		if y > 199 {
			return
		}
		// only store if inside the screen, we should really clip
		if y >= 0 {
			// is first slot free?
			if p.edge_table[y][0].x == -1 {
				p.edge_table[y][0].x = x1
				p.edge_table[y][0].tx = tx1
				p.edge_table[y][0].ty = ty1
				p.edge_table[y][0].px = px1
				p.edge_table[y][0].py = py1
				p.edge_table[y][0].z = z1
			} else { // otherwise use the other
				p.edge_table[y][1].x = x1
				p.edge_table[y][1].tx = tx1
				p.edge_table[y][1].ty = ty1
				p.edge_table[y][1].px = px1
				p.edge_table[y][1].py = py1
				p.edge_table[y][1].z = z1
			}
		}
		// interpolate our values
		x1 += dx
		px1 += dpx
		py1 += dpy
		tx1 += dtx
		ty1 += dty
		z1 += dz
	}
}

/*
 * draw a horizontal double textured span
 */
func (p *Polygon) drawSpan(buffer *image.RGBA, y int, p1, p2 *edge_data) {
	// quick check, if facing back then draw span in the other direction,
	// avoids having to swap all the vars... not a very elegant
	if p1.x > p2.x {
		p.drawSpan(buffer, y, p2, p1)
		return
	}
	// load starting points
	z1 := p1.z
	px1 := p1.px
	py1 := p1.py
	tx1 := p1.tx
	ty1 := p1.ty
	x1 := p1.x >> 16
	x2 := p2.x >> 16

	// check if it's inside the screen
	if (x1 > 319) || (x2 < 0) {
		return
	}

	// compute deltas for interpolation
	dx := x2 - x1
	if dx == 0 {
		return
	}
	dtx := (p2.tx - p1.tx) / dx
	dty := (p2.ty - p1.ty) / dx
	dpx := (p2.px - p1.px) / dx
	dpy := (p2.py - p1.py) / dx
	dz := (p2.z - p1.z) / dx
	// get destination offset in buffer
	offs := y*320 + x1
	// loop for all pixels concerned
	for i := x1; i < x2; i++ {
		if i > 319 {
			return
		}
		// check z buffer
		if i >= 0 {
			if z1 < int(p.zbuffer[offs]) {
				// if visible load the texel from the translated texture
				texel := int(p.texture.Pix[((ty1>>8)&0xff00)+((tx1>>16)&0xff)])
				lumel := int(p.light[((py1>>8)&0xff00)+((px1>>16)&0xff)])
				palIndex := p.lut[(texel<<8)+lumel]
				pr, pg, pb, _ := p.texture.Palette[palIndex].RGBA()
				buffer.Pix[offs*4+0] = byte(pr)
				buffer.Pix[offs*4+1] = byte(pg)
				buffer.Pix[offs*4+2] = byte(pb)
				//buffer.Pix[offs*4+3] = byte(pa)
				// mix them together, and store
				// and update the zbuffer
				p.zbuffer[offs] = uint16(z1)
			}
		}
		// interpolate our values
		px1 += dpx
		py1 += dpy
		tx1 += dtx
		ty1 += dty
		z1 += dz
		offs++
	}
}
