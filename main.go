// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"image/color"
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
	crosshairRadius      = 10
	crosshairInnerRadius = 3
	WallElasticity       = 1
	wallFriction         = 1
	wallWidth            = 30
	spaceIterations      = 5
)

const (
	wallLeftCenterX        = wallWidth
	wallLeftCenterY        = screenHeight - 2.0*screenHeight/5.0
	wallLeftCenterWidth    = screenWidth / 4.0
	wallRightCenterX       = screenWidth - wallWidth - screenWidth/4.0
	wallRightCenterY       = 2.0 * screenHeight / 5.0
	wallRightCenterWidth   = screenWidth / 4.0
	wallTopCenterX         = screenWidth / 4.0
	wallTopCenterY         = wallWidth
	wallTopCenterHeight    = screenHeight / 4.0
	wallBottomCenterX      = wallRightCenterX
	wallBottomCenterY      = screenHeight - wallWidth - wallBottomCenterHeight
	wallBottomCenterHeight = screenHeight / 4.0
)

var (
	colorBackground = color.RGBA{124, 144, 160, 255} // ~ Light Slate Gray
	colorWall       = color.RGBA{57, 62, 65, 255}    // ~ Onyx
	colorGun        = color.RGBA{242, 129, 35, 255}  // ~ Princeton Orange
	colorPlayer     = color.RGBA{155, 201, 149, 255} // ~ Dark Sea Green
)

var (
	imageWall                   = ebiten.NewImage(1, 1)
	imageCursor                 = ebiten.NewImage(crosshairRadius*2, crosshairRadius*2)
	drawOptionsCursor           ebiten.DrawImageOptions
	drawOptionsWallTop          ebiten.DrawImageOptions
	drawOptionsWallLeft         ebiten.DrawImageOptions
	drawOptionsWallRight        ebiten.DrawImageOptions
	drawOptionsWallBottom       ebiten.DrawImageOptions
	drawOptionsWallLeftCenter   ebiten.DrawImageOptions
	drawOptionsWallRightCenter  ebiten.DrawImageOptions
	drawOptionsWallTopCenter    ebiten.DrawImageOptions
	drawOptionsWallBottomCenter ebiten.DrawImageOptions
)

func init() {
	imageWall.Fill(colorWall)

	drawOptionsWallTop.GeoM.Scale(screenWidth, wallWidth)

	drawOptionsWallLeft.GeoM.Scale(wallWidth, screenHeight)

	drawOptionsWallRight.GeoM = drawOptionsWallLeft.GeoM
	drawOptionsWallRight.GeoM.Translate(screenWidth-wallWidth, 0)

	drawOptionsWallBottom.GeoM = drawOptionsWallTop.GeoM
	drawOptionsWallBottom.GeoM.Translate(0, screenHeight-wallWidth)

	drawOptionsWallLeftCenter.GeoM.Scale(wallLeftCenterWidth, wallWidth)
	drawOptionsWallLeftCenter.GeoM.Translate(wallLeftCenterX, wallLeftCenterY)

	drawOptionsWallRightCenter.GeoM.Scale(wallRightCenterWidth, wallWidth)
	drawOptionsWallRightCenter.GeoM.Translate(wallRightCenterX, wallRightCenterY)

	drawOptionsWallTopCenter.GeoM.Scale(wallWidth, wallTopCenterHeight)
	drawOptionsWallTopCenter.GeoM.Translate(wallTopCenterX, wallTopCenterY)

	drawOptionsWallBottomCenter.GeoM.Scale(wallWidth, wallBottomCenterHeight)
	drawOptionsWallBottomCenter.GeoM.Translate(wallBottomCenterX, wallBottomCenterY)

	initCursorImage()
}

func initCursorImage() {
	ebitenutil.DrawLine(imageCursor, 0, crosshairRadius,
		crosshairRadius-crosshairInnerRadius, crosshairRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, 0,
		crosshairRadius, crosshairRadius-crosshairInnerRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius+crosshairInnerRadius,
		crosshairRadius, 2*crosshairRadius, crosshairRadius, colornames.Red)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, crosshairRadius+crosshairInnerRadius,
		crosshairRadius, 2*crosshairRadius, colornames.Red)
}

// game implements ebiten.game interface.
type game struct {
	player player
	space  *cp.Space
	input  input
}

func newGame() *game {
	game := &game{
		player: *newPlayer(cp.Vector{X: screenWidth / 2.0, Y: screenHeight / 2.0}),
	}

	space := cp.NewSpace()
	space.Iterations = spaceIterations
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	// Add player to the space
	space.AddBody(game.player.body)
	space.AddShape(game.player.shape)
	game.space = space

	// Add walls to the space
	shape := space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: 0, Y: wallWidth}, cp.Vector{X: screenWidth, Y: wallWidth}, 0)) // Top wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: wallWidth, Y: 0}, cp.Vector{X: wallWidth, Y: screenHeight}, 0)) // left wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: screenWidth - wallWidth, Y: 0}, cp.Vector{X: screenWidth - wallWidth, Y: screenHeight}, 0)) // right wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: 0, Y: screenHeight - wallWidth}, cp.Vector{X: screenWidth, Y: screenHeight - wallWidth}, 0)) // bottom wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	const wallRadius = wallWidth / 2.0
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: wallLeftCenterX, Y: wallLeftCenterY + wallRadius},
		cp.Vector{X: wallLeftCenterX + wallLeftCenterWidth - wallRadius, Y: wallLeftCenterY + wallRadius}, wallRadius)) // left center wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: wallRightCenterX + wallRadius, Y: wallRightCenterY + wallRadius},
		cp.Vector{X: wallRightCenterX + wallRightCenterWidth - wallRadius, Y: wallRightCenterY + wallRadius}, wallRadius)) // right center wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: wallTopCenterX + wallRadius, Y: wallTopCenterY + wallRadius},
		cp.Vector{X: wallTopCenterX + wallRadius, Y: wallTopCenterY + wallTopCenterHeight - wallRadius}, wallRadius)) // top center wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)
	shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: wallBottomCenterX + wallRadius, Y: wallBottomCenterY + wallRadius},
		cp.Vector{X: wallBottomCenterX + wallRadius, Y: wallBottomCenterY + wallBottomCenterHeight - wallRadius}, wallRadius)) // bottom center wall
	shape.SetElasticity(WallElasticity)
	shape.SetFriction(wallFriction)

	return game
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	g.space.Step(deltaTime)

	// Update input states(mouse pos and pressed keys)
	g.input.update()

	// Update player and player's gun
	g.player.update(&g.input)

	// Update Crosshair geometry matrix
	drawOptionsCursor.GeoM.Reset()
	drawOptionsCursor.GeoM.Translate(g.input.cursorPos.X-crosshairRadius, g.input.cursorPos.Y-crosshairRadius)
	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// Draw walls
	screen.DrawImage(imageWall, &drawOptionsWallTop)
	screen.DrawImage(imageWall, &drawOptionsWallBottom)
	screen.DrawImage(imageWall, &drawOptionsWallLeft)
	screen.DrawImage(imageWall, &drawOptionsWallRight)
	screen.DrawImage(imageWall, &drawOptionsWallLeftCenter)
	screen.DrawImage(imageWall, &drawOptionsWallRightCenter)
	screen.DrawImage(imageWall, &drawOptionsWallTopCenter)
	screen.DrawImage(imageWall, &drawOptionsWallBottomCenter)

	// Draw player and its gun
	g.player.draw(screen)

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
