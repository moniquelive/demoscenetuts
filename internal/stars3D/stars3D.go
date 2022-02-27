package stars3D

import (
	"image"

	"github.com/moniquelive/demoscenetuts/internal/utils"
)

const maxStars = 500

type Star struct {
	screenWidth,
	screenHeight int
	x, y, z int
}

func NewStar(screenWidth, screenHeight int) *Star {
	x := utils.Between(-screenWidth/2, screenWidth/2)
	y := utils.Between(-screenHeight/2, screenHeight/2)
	z := utils.Between(1, screenWidth)
	return &Star{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		x:            x,
		y:            y,
		z:            z,
	}
}

func (s *Star) Update() {
	s.z -= 1
	if s.z < 1 {
		s.x = utils.Between(-s.screenWidth/2, s.screenWidth/2)
		s.y = utils.Between(-s.screenHeight/2, s.screenHeight/2)
		s.z = s.screenWidth
	}
}

func (s Star) Draw(buffer *image.RGBA) {
	const factor = 20
	x := (s.x*factor)/s.z + (s.screenWidth / 2)
	y := (s.y*factor)/s.z + (s.screenHeight / 2)
	if x < 0 || x >= s.screenWidth || y < 0 || y >= s.screenHeight {
		return
	}

	i := s.screenWidth*y + x
	gray := uint8(s.screenWidth * 8 / s.z)
	buffer.Pix[4*i+0] = gray
	buffer.Pix[4*i+1] = gray
	buffer.Pix[4*i+2] = gray
	buffer.Pix[4*i+3] = uint8(0xff)
}

//------------------------------------------------------------------------------

type Stars struct {
	stars []*Star
}

func (s Stars) Draw(buffer *image.RGBA) {
	for i := 0; i < maxStars; i++ {
		s.stars[i].Update()
		s.stars[i].Draw(buffer)
	}
}

func (s *Stars) Setup() (int, int, int) {
	screenWidth, screenHeight := 320, 200
	s.stars = make([]*Star, maxStars)
	for i := 0; i < maxStars; i++ {
		s.stars[i] = NewStar(screenWidth, screenHeight)
	}
	return screenWidth, screenHeight, 2
}
