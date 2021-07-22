package polygon

import "math"

type Matrix [3]Vector

func NewMatrix(a11, a12, a13, a21, a22, a23, a31, a32, a33 float64) (m Matrix) {
	m[0][0] = a11
	m[0][1] = a12
	m[0][2] = a13
	m[1][0] = a21
	m[1][1] = a22
	m[1][2] = a23
	m[2][0] = a31
	m[2][1] = a32
	m[2][2] = a33
	return
}

func (m Matrix) MulMat(b Matrix) (r Matrix) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r[i][j] = m[i][0]*b[0][j] +
				m[i][1]*b[1][j] +
				m[i][2]*b[2][j]
		}
	}
	return
}

func (m Matrix) MulVec(v Vector) (r Vector) {
	r[0] = v[0]*m[0][0] + v[1]*m[1][0] + v[2]*m[2][0]
	r[1] = v[0]*m[0][1] + v[1]*m[1][1] + v[2]*m[2][1]
	r[2] = v[0]*m[0][2] + v[1]*m[1][2] + v[2]*m[2][2]
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
