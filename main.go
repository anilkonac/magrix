// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"image"
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

const (
	ratioLandHeight = 1.0 / 4.0
	playerWidth     = 40.0
	playerHeight    = 100.0
)

var (
	imageEmpty        = ebiten.NewImage(2, 1)
	imageLand         *ebiten.Image
	imagePlayer       *ebiten.Image
	drawOptionsLand   = &ebiten.DrawImageOptions{}
	drawOptionsPlayer = &ebiten.DrawImageOptions{}
)

func init() {
	imageLand = imageEmpty.SubImage(image.Rect(0, 0, 1, 1)).(*ebiten.Image)
	imagePlayer = imageEmpty.SubImage(image.Rect(1, 0, 2, 1)).(*ebiten.Image)
	imageLand.Fill(colornames.Forestgreen)
	imagePlayer.Fill(colornames.Slategray)

	drawOptionsLand.GeoM.Scale(screenWidth, screenHeight*ratioLandHeight)
	drawOptionsLand.GeoM.Translate(0, screenHeight*(1.0-ratioLandHeight))
}

// game implements ebiten.game interface.
type game struct {
	playerX, playerY float64
}

func newGame() *game {
	return &game{
		playerX: screenWidth / 2.0,
		playerY: screenHeight / 2.0,
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	// Write your game's logical update.
	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Lightskyblue)

	// Draw land
	screen.DrawImage(imageLand, drawOptionsLand)

	// Draw prototype player
	drawOptionsPlayer.GeoM.Reset()
	drawOptionsPlayer.GeoM.Scale(playerWidth, playerHeight)
	drawOptionsPlayer.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(imagePlayer, drawOptionsPlayer)

	// Print fps
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// return outsideWidth, outsideHeight
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Magrix")
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
