package polygon

import "math"

type Vector [3]float64

func NewVector(x, y, z float64) (r Vector) {
	r[0] = x
	r[1] = y
	r[2] = z
	return
}

func (v Vector) Add(a Vector) (r Vector) {
	r[0] = v[0] + a[0]
	r[1] = v[1] + a[1]
	r[2] = v[2] + a[2]
	return
}

func (v Vector) Sub(a Vector) (r Vector) {
	r[0] = v[0] - a[0]
	r[1] = v[1] - a[1]
	r[2] = v[2] - a[2]
	return
}

func (v Vector) normalize() Vector {
	id := 1.0 / math.Sqrt(v[0]*v[0]+v[1]*v[1]+v[2]*v[2])
	return NewVector(v[0]*id, v[1]*id, v[2]*id)
}

//func (v *Vector) Copy(a Vector) {
//	v[0] = a[0]
//	v[1] = a[1]
//	v[2] = a[2]
//}
