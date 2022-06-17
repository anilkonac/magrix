// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 960
	screenHeight = 720
	deltaTime    = 1.0 / 60.0
)

const (
	ratioLandHeight      = 1.0 / 4.0
	landY                = (1.0 - ratioLandHeight) * screenHeight
	playerWidth          = 40.0
	playerHeight         = 100.0
	gunWidth             = playerHeight / 2.0
	gunHeight            = playerWidth / 3.0
	crosshairRadius      = 10
	crosshairInnerRadius = 3
)

const (
	gravity = 500.0
)

var (
	imageLand         = ebiten.NewImage(1, 1)
	imagePlayer       = ebiten.NewImage(1, 1)
	imageGun          = ebiten.NewImage(1, 1)
	imageCursor       = ebiten.NewImage(crosshairRadius*2, crosshairRadius*2)
	drawOptionsLand   ebiten.DrawImageOptions
	drawOptionsCursor ebiten.DrawImageOptions
)

func init() {
	imageLand.Fill(colornames.Forestgreen)
	imagePlayer.Fill(colornames.Slategray)
	imageGun.Fill(colornames.Orange)

	drawOptionsLand.GeoM.Scale(screenWidth, screenHeight*ratioLandHeight)
	drawOptionsLand.GeoM.Translate(0, landY)

	initCursorImage()
}

func initCursorImage() {
	ebitenutil.DrawLine(imageCursor, 0, crosshairRadius, crosshairRadius-crosshairInnerRadius, crosshairRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, 0, crosshairRadius, crosshairRadius-crosshairInnerRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius+crosshairInnerRadius, crosshairRadius, 2*crosshairRadius, crosshairRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, crosshairRadius+crosshairInnerRadius, crosshairRadius, 2*crosshairRadius, colornames.Red)
}

// game implements ebiten.game interface.
type game struct {
	player     player
	posCursor  cp.Vector
	gunAngle   float64
	gunX, gunY float64
	space      *cp.Space
	// playerBody         *cp.Body
}

func newGame() *game {
	game := &game{
		player: *newPlayer(cp.Vector{X: screenWidth / 2.0, Y: screenHeight / 2.0}),
	}

	space := cp.NewSpace()
	space.Iterations = 1
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	// Player

	space.AddBody(game.player.shape.Body())
	space.AddShape(game.player.shape)
	game.space = space

	// Land
	space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: 0, Y: landY}, cp.Vector{X: screenWidth, Y: landY}, 0))

	return game
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	// g.space.Step(1.0 / float64(ebiten.MaxTPS()))
	g.space.Step(deltaTime)
	x, y := ebiten.CursorPosition()
	g.posCursor = cp.Vector{X: float64(x), Y: float64(y)}

	g.player.update(&g.posCursor)

	g.updateGeometryMatrices()
	return nil
}

func (g *game) updateGeometryMatrices() {
	// Crosshair
	drawOptionsCursor.GeoM.Reset()
	drawOptionsCursor.GeoM.Translate(g.posCursor.X-crosshairRadius, g.posCursor.Y-crosshairRadius)
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f  FPS: %.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %.0f, Y: %.0f", g.fCursorX, g.fCursorY), 0, 15)
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
