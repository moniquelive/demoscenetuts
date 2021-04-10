package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 200
)

var stars []*Star

type Game struct {
	doubleBuffer *image.RGBA
}

func (g *Game) Update() error {
	draw.Draw(
		g.doubleBuffer, g.doubleBuffer.Bounds(),
		image.Black, image.Black.Bounds().Min,
		draw.Src)
	for i := 0; i < maxStars; i++ {
		stars[i].Update()
		stars[i].Draw(g.doubleBuffer)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.doubleBuffer.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ebiten Demoscene")
	stars = make([]*Star, maxStars)
	for i := 0; i < maxStars; i++ {
		stars[i] = NewStar(
			float64(theRand.next()%screenWidth),
			float64(theRand.next()%screenHeight),
			float64(theRand.next()%maxPlanes))
	}
	g := &Game{
		doubleBuffer: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
