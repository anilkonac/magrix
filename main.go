// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 960
	screenHeight = 720
	deltaTime    = 1.0 / 60.0
)

const ratioLandHeight = 1.0 / 4.0

var (
	landImage       = ebiten.NewImage(1, 1)
	landDrawOptions = &ebiten.DrawImageOptions{}
)

func init() {
	landImage.Fill(colornames.Forestgreen)
	landDrawOptions.GeoM.Scale(screenWidth, screenHeight*ratioLandHeight)
	landDrawOptions.GeoM.Translate(0, screenHeight*(1.0-ratioLandHeight))

}

// game implements ebiten.game interface.
type game struct{}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	// Write your game's logical update.
	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Lightskyblue)

	// Draw land
	screen.DrawImage(landImage, landDrawOptions)

	// Print fps
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Magrix")

	game := &game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
