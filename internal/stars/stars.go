package stars

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moniquelive/demoscenetuts/internal/utils"
)

const (
	maxStars  = 500
	maxPlanes = 50
)

var xVel float64

type Star struct {
	screenWidth,
	screenHeight int
	x, y, p float64
}

func NewStar(screenWidth, screenHeight int) *Star {
	return &Star{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		x:            float64(utils.Between(0, screenWidth)),
		y:            float64(utils.Between(0, screenHeight)),
		p:            float64(utils.Between(0, maxPlanes)),
	}
}

func (s *Star) Update() {
	s.updateXVel()
	s.x += (1 + s.p) * xVel
	if s.x >= float64(s.screenWidth) {
		s.x = 0
		s.y = float64(utils.Between(0, s.screenHeight))
	}
}

func (s Star) Draw(buffer *image.RGBA) {
	i := int(float64(s.screenWidth)*s.y + s.x)
	if 4*i+3 > len(buffer.Pix) {
		return
	}

	gray := uint8((256 / maxPlanes) * s.p)
	buffer.Pix[4*i+0] = gray
	buffer.Pix[4*i+1] = gray
	buffer.Pix[4*i+2] = gray
	buffer.Pix[4*i+3] = uint8(0xff)
}

func (s Star) updateXVel() {
	mouseX, mouseY := ebiten.CursorPosition()
	if mouseX > 0 && mouseX < s.screenWidth &&
		mouseY > 0 && mouseY < s.screenHeight {
		xVel = utils.Interpolate(float64(mouseX), 0, float64(s.screenWidth), 0, 0.25)
	}
}

type Stars struct {
	stars        []*Star
	screenWidth  int
	screenHeight int
}

func (s Stars) Draw(buffer *image.RGBA) {
	for i := 0; i < maxStars; i++ {
		s.stars[i].Update()
		s.stars[i].Draw(buffer)
	}
}

func (s *Stars) Setup() (int, int, int) {
	s.screenWidth, s.screenHeight = 320, 200
	s.stars = make([]*Star, maxStars)
	for i := 0; i < maxStars; i++ {
		s.stars[i] = NewStar(s.screenWidth, s.screenHeight)
	}
	return s.screenWidth, s.screenHeight, 2
}
