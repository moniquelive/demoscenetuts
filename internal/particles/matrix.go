package particles

import "math"

type Matrix struct {
	m [3]Vector
}

func NewMatrix(a11, a12, a13, a21, a22, a23, a31, a32, a33 float64) (m Matrix) {
	m.m[0].v[0] = a11
	m.m[0].v[1] = a12
	m.m[0].v[2] = a13
	m.m[1].v[0] = a21
	m.m[1].v[1] = a22
	m.m[1].v[2] = a23
	m.m[2].v[0] = a31
	m.m[2].v[1] = a32
	m.m[2].v[2] = a33
	return
}

func (m Matrix) MulMat(b Matrix) (r Matrix) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r.m[i].v[j] = m.m[i].v[0]*b.m[0].v[j] +
				m.m[i].v[1]*b.m[1].v[j] +
				m.m[i].v[2]*b.m[2].v[j]
		}
	}
	return
}

func (m Matrix) MulVec(v Vector) (r Vector) {
	r.v[0] = v.v[0]*m.m[0].v[0] + v.v[1]*m.m[1].v[0] + v.v[2]*m.m[2].v[0]
	r.v[1] = v.v[0]*m.m[0].v[1] + v.v[1]*m.m[1].v[1] + v.v[2]*m.m[2].v[1]
	r.v[2] = v.v[0]*m.m[0].v[2] + v.v[1]*m.m[1].v[2] + v.v[2]*m.m[2].v[2]
	return
}

func rotX(theta float64) Matrix {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return NewMatrix(
		1, 0, 0,
		0, c, s,
		0, -s, c)
}

func rotY(theta float64) Matrix {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return NewMatrix(
		c, 0, -s,
		0, 1, 0,
		s, 0, c)
}

func rotZ(theta float64) Matrix {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return NewMatrix(
		c, s, 0,
		-s, c, 0,
		0, 0, 1)
}
