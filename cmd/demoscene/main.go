package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/moniquelive/demoscenetuts/internal/stars"
	"github.com/moniquelive/demoscenetuts/internal/stars3D"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 200
)

type demo interface {
	Draw(*image.RGBA)
	Setup(int, int)
}

var d demo

type Game struct {
	doubleBuffer *image.RGBA
}

var demos map[string]demo

func init() {
	demos = make(map[string]demo)
	demos["stars"] = &stars.Stars{}
	demos["3d"] = &stars3D.Stars{}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("end")
	}
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
	if len(os.Args) != 2 {
		help()
		return
	}
	var ok bool
	if d, ok = demos[os.Args[1]]; !ok {
		log.Printf("Demo não encontrado: %q\n", os.Args[1])
		help()
		return
	}
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

func help() {
	demosList := make([]string, 0, len(demos))
	for k := range demos {
		demosList = append(demosList, k)
	}
	log.Println("Escolha o demo: [", strings.Join(demosList, ", "), "]")
	return
}
