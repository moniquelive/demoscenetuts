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

const (
	ScreenWidth  = 320
	ScreenHeight = 200
)

func updateXVel() {
	mouseX, mouseY := ebiten.CursorPosition()
	if mouseX > 0 && mouseX < ScreenWidth &&
		mouseY > 0 && mouseY < ScreenHeight {
		xVel = utils.Interpolate(mouseX, 0, ScreenWidth, 0, 0.25)
	}
}

type Star struct {
	x, y, p float64
}

func NewStar(x, y, p float64) *Star {
	return &Star{x: x, y: y, p: p}
}

func (s *Star) Update() {
	updateXVel()
	s.x += (1 + s.p) * xVel
	if s.x >= ScreenWidth {
		s.x = 0
		s.y = float64(utils.TheRand.Next() % ScreenHeight)
	}
}

func (s Star) Draw(buffer *image.RGBA) {
	i := int(ScreenWidth*s.y + s.x)
	if 4*i+3 > len(buffer.Pix) {
		return
	}

	gray := uint8((256 / maxPlanes) * s.p)
	buffer.Pix[4*i+0] = gray
	buffer.Pix[4*i+1] = gray
	buffer.Pix[4*i+2] = gray
	buffer.Pix[4*i+3] = uint8(0xff)
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

func (s *Stars) Setup(screenWidth, screenHeight int) {
	s.screenWidth = screenWidth
	s.screenHeight = screenHeight
	s.stars = make([]*Star, maxStars)
	for i := 0; i < maxStars; i++ {
		s.stars[i] = NewStar(
			float64(utils.TheRand.Next()%ScreenWidth),
			float64(utils.TheRand.Next()%ScreenHeight),
			float64(utils.TheRand.Next()%maxPlanes))
	}
}
