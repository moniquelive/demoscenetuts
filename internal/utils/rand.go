package utils

import "math"

type rand struct {
	x, y, z, w uint32
}

var TheRand = &rand{12345678, 4185243, 776511, 45411}

func (r *rand) Next() uint32 {
	// math/rand is too slow to keep 60 FPS on web browsers.
	// Use Xorshift instead: http://en.wikipedia.org/wiki/Xorshift
	t := r.x ^ (r.x << 11)
	r.x, r.y, r.z = r.y, r.z, r.w
	r.w = (r.w ^ (r.w >> 19)) ^ (t ^ (t >> 8))
	return r.w
}

//func NextUint32() uint32 {
//	return TheRand.Next()
//}

func Random() float64 {
	return float64(TheRand.Next()) / math.MaxUint32
}

func Between(min, max int) int {
	return int(Random()*float64(max-min)) + min
}
