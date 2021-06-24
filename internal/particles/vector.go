package particles

type Vector struct {
	v [3]float64
}

func NewVector(x, y, z float64) (r Vector) {
	r.v[0] = x
	r.v[1] = y
	r.v[2] = z
	return
}

//func (v Vector) Add(a Vector) (r Vector) {
//	r.v[0] = v.v[0] + a.v[0]
//	r.v[1] = v.v[1] + a.v[1]
//	r.v[2] = v.v[2] + a.v[2]
//	return
//}
//
//func (v *Vector) Copy(a Vector) {
//	v.v[0] = a.v[0]
//	v.v[1] = a.v[1]
//	v.v[2] = a.v[2]
//}
