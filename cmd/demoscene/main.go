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
	"github.com/moniquelive/demoscenetuts/internal"
)

var (
	ScreenWidth  = 320
	ScreenHeight = 200
)

type Game struct {
	doubleBuffer *image.RGBA
}

var currentDemo demos.Demo

var demosMap map[string]demos.Demo

func init() {
	demosMap = demos.FillDemos()
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("end")
	}
	draw.Draw(
		g.doubleBuffer, g.doubleBuffer.Bounds(),
		image.Black, image.Black.Bounds().Min,
		draw.Src)
	currentDemo.Draw(g.doubleBuffer)
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
	if currentDemo, ok = demosMap[os.Args[1]]; !ok {
		log.Printf("Demo não encontrado: %q\n", os.Args[1])
		help()
		return
	}
	var zoom int
	ScreenWidth, ScreenHeight, zoom = currentDemo.Setup()
	ebiten.SetWindowSize(ScreenWidth*zoom, ScreenHeight*zoom)
	ebiten.SetWindowTitle("Ebiten Demoscene")
	g := &Game{
		doubleBuffer: image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight)),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func help() {
	demosList := make([]string, 0, len(demosMap))
	for k := range demosMap {
		demosList = append(demosList, k)
	}
	log.Println("Escolha o demo: [", strings.Join(demosList, ", "), "]")
}
