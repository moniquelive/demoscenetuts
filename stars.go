package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	maxStars  = 500
	maxPlanes = 50
)

var xVel float64

func updateXVel() {
	mouseX, mouseY := ebiten.CursorPosition()
	if mouseX > 0 && mouseX < screenWidth &&
		mouseY > 0 && mouseY < screenHeight {
		xVel = interpolate(mouseX, 0, screenWidth, 0, 0.25)
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
	if s.x >= screenWidth {
		s.x = 0
		s.y = float64(theRand.next() % screenHeight)
	}
}

func (s Star) Draw(buffer *image.RGBA) {
	i := int(screenWidth*s.y + s.x)
	if 4*i+3 > len(buffer.Pix) {
		return
	}

	gray := uint8((256 / maxPlanes) * s.p)
	buffer.Pix[4*i+0] = gray
	buffer.Pix[4*i+1] = gray
	buffer.Pix[4*i+2] = gray
	buffer.Pix[4*i+3] = uint8(0xff)
}
