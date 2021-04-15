package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/moniquelive/demoscenetuts/internal/stars"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 200
)

var d = &stars.Stars{}

type Game struct {
	doubleBuffer *image.RGBA
}

func (g *Game) Update() error {
	draw.Draw(
		g.doubleBuffer, g.doubleBuffer.Bounds(),
		image.Black, image.Black.Bounds().Min,
		draw.Src)
	d.Draw(g.doubleBuffer)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.doubleBuffer.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("Ebiten Demoscene")
	d.Setup(ScreenWidth, ScreenHeight)
	g := &Game{
		doubleBuffer: image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight)),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
