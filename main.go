// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"log"
	"math"

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
	ratioLandHeight      = 1.0 / 4.0
	playerWidth          = 40.0
	playerHeight         = 100.0
	gunWidth             = playerHeight / 2.0
	gunHeight            = playerWidth / 3.0
	crosshairRadius      = 10
	crosshairInnerRadius = 3
)

var (
	imageLand         = ebiten.NewImage(1, 1)
	imagePlayer       = ebiten.NewImage(1, 1)
	imageGun          = ebiten.NewImage(1, 1)
	imageCursor       = ebiten.NewImage(crosshairRadius*2, crosshairRadius*2)
	drawOptionsLand   ebiten.DrawImageOptions
	drawOptionsPlayer ebiten.DrawImageOptions
	drawOptionsGun    ebiten.DrawImageOptions
	drawOptionsCursor ebiten.DrawImageOptions
)

func init() {
	imageLand.Fill(colornames.Forestgreen)
	imagePlayer.Fill(colornames.Slategray)
	imageGun.Fill(colornames.Orange)

	drawOptionsLand.GeoM.Scale(screenWidth, screenHeight*ratioLandHeight)
	drawOptionsLand.GeoM.Translate(0, screenHeight*(1.0-ratioLandHeight))

	// Prepare cursor image
	ebitenutil.DrawLine(imageCursor, 0, crosshairRadius, crosshairRadius-crosshairInnerRadius, crosshairRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, 0, crosshairRadius, crosshairRadius-crosshairInnerRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius+crosshairInnerRadius, crosshairRadius, 2*crosshairRadius, crosshairRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, crosshairRadius+crosshairInnerRadius, crosshairRadius, 2*crosshairRadius, colornames.Red)
}

// game implements ebiten.game interface.
type game struct {
	playerX, playerY   float64
	fCursorX, fCursorY float64
	gunAngle           float64
	gunX, gunY         float64
}

func newGame() *game {
	return &game{
		playerX: screenWidth / 2.0,
		playerY: screenHeight / 2.0,
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	x, y := ebiten.CursorPosition()
	g.fCursorX, g.fCursorY = float64(x), float64(y)

	// Update gun position
	g.gunX = g.playerX + playerWidth/2.0
	g.gunY = g.playerY + playerHeight/2.0

	// Calculate gun rotation angle
	distX := g.fCursorX - g.gunX
	distY := g.fCursorY - g.gunY
	g.gunAngle = math.Atan2(distY, distX)

	g.updateGeometryMatrices()
	return nil
}

func (g *game) updateGeometryMatrices() {
	// Player
	drawOptionsPlayer.GeoM.Reset()
	drawOptionsPlayer.GeoM.Scale(playerWidth, playerHeight)
	drawOptionsPlayer.GeoM.Translate(g.playerX, g.playerY)

	// Gun
	drawOptionsGun.GeoM.Reset()
	drawOptionsGun.GeoM.Scale(gunWidth, gunHeight)
	drawOptionsGun.GeoM.Translate(0, -gunHeight/2.0)
	drawOptionsGun.GeoM.Rotate(g.gunAngle)
	drawOptionsGun.GeoM.Translate(g.gunX, g.gunY)

	// Crosshair
	drawOptionsCursor.GeoM.Reset()
	drawOptionsCursor.GeoM.Translate(g.fCursorX-crosshairRadius, g.fCursorY-crosshairRadius)
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Lightskyblue)

	// Draw land
	screen.DrawImage(imageLand, &drawOptionsLand)

	// Draw prototype player
	screen.DrawImage(imagePlayer, &drawOptionsPlayer)

	// Draw prototype gun
	screen.DrawImage(imageGun, &drawOptionsGun)

	// Draw crosshair
	screen.DrawImage(imageCursor, &drawOptionsCursor)

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
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
