// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
)

const (
	screenWidth  = 960
	screenHeight = 720
	deltaTime    = 1.0 / 60.0
)

const (
	ratioLandHeight      = 1.0 / 4.0
	landY                = (1.0 - ratioLandHeight) * screenHeight
	crosshairRadius      = 14
	crosshairInnerRadius = 4
	rayHitImageWidth     = 16
	wallElasticity       = 1
	wallFriction         = 1
	wallWidth            = 30
	wallRadius           = wallWidth / 2.0
	// spaceIterations      = 10
)

var (
	colorBackground = color.RGBA{124, 144, 160, 255} // ~ Light Slate Gray
	colorWall       = color.RGBA{57, 62, 65, 255}    // ~ Onyx
	colorGun        = color.RGBA{242, 129, 35, 255}  // ~ Princeton Orange
	colorPlayer     = color.RGBA{155, 201, 149, 255} // ~ Dark Sea Green
	colorCrosshair  = color.RGBA{216, 17, 89, 255}   // ~ Ruby
	colorRayHit     = color.RGBA{7, 160, 195, 255}   // ~ Blue Green
)

var (
	imageWall         = ebiten.NewImage(1, 1)
	imageCursor       = ebiten.NewImage(crosshairRadius*2, crosshairRadius*2)
	imageRayHit       = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	drawOptionsCursor ebiten.DrawImageOptions
	drawOptionsRayHit ebiten.DrawImageOptions
)

func init() {
	initCursorImage()

	shader, err := ebiten.NewShader(circle_go)
	if err != nil {
		panic(err)
	}

	imageRayHit.DrawRectShader(rayHitImageWidth, rayHitImageWidth, shader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius": float32(rayHitImageWidth / 2.0),
			"Color":  []float32{float32(colorRayHit.R / 0xf), float32(colorRayHit.G / 0xf), float32(colorRayHit.B / 0xf), float32(colorRayHit.A / 0xf)},
		},
	})
}

func initCursorImage() {
	ebitenutil.DrawLine(imageCursor, 0, crosshairRadius,
		crosshairRadius-crosshairInnerRadius, crosshairRadius, colorCrosshair)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, 0,
		crosshairRadius, crosshairRadius-crosshairInnerRadius, colorCrosshair)
	ebitenutil.DrawLine(imageCursor, crosshairRadius+crosshairInnerRadius,
		crosshairRadius, 2*crosshairRadius, crosshairRadius, colorCrosshair)
	ebitenutil.DrawLine(imageCursor, crosshairRadius, crosshairRadius+crosshairInnerRadius,
		crosshairRadius, 2*crosshairRadius, colorCrosshair)
}

// game implements ebiten.game interface.
type game struct {
	player     player
	space      *cp.Space
	walls      []*wall
	input      input
	rayHitInfo cp.SegmentQueryInfo
}

func newGame() *game {
	game := &game{
		player: *newPlayer(cp.Vector{X: screenWidth / 2.0, Y: screenHeight / 2.0}),
	}

	space := cp.NewSpace()
	// space.Iterations = spaceIterations
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	// Add player to the space
	space.AddBody(game.player.body)
	space.AddShape(game.player.shape)
	game.space = space

	addWalls(space, &game.walls)

	return game
}

func addWalls(space *cp.Space, walls *[]*wall) {
	const (
		wallLeftCenterX        = 3 * wallRadius
		wallLeftCenterY        = screenHeight - 2.0*screenHeight/5.0
		wallLeftCenterWidth    = screenWidth / 4.0
		wallRightCenterX       = screenWidth - screenWidth/4.0 - wallWidth
		wallRightCenterY       = 2.0 * screenHeight / 5.0
		wallRightCenterWidth   = screenWidth / 4.0
		wallTopCenterX         = screenWidth/4.0 + wallWidth
		wallTopCenterY         = 3 * wallRadius
		wallTopCenterHeight    = screenHeight / 4.0
		wallBottomCenterX      = wallRightCenterX
		wallBottomCenterHeight = screenHeight / 4.0
		wallBottomCenterY      = screenHeight - wallWidth - wallBottomCenterHeight
	)

	*walls = append(*walls, newWall(wallRadius, wallRadius, screenWidth-wallRadius, wallRadius, wallRadius, space))                                                   // Top wall
	*walls = append(*walls, newWall(wallRadius, screenHeight-wallRadius, screenWidth-wallRadius, screenHeight-wallRadius, wallRadius, space))                         // Bottom wall
	*walls = append(*walls, newWall(wallRadius, 0, wallRadius, screenHeight-wallRadius, wallRadius, space))                                                           // left wall
	*walls = append(*walls, newWall(screenWidth-wallRadius, 0, screenWidth-wallRadius, screenHeight-wallRadius, wallRadius, space))                                   // right wall
	*walls = append(*walls, newWall(wallLeftCenterX, wallLeftCenterY, wallLeftCenterX+wallLeftCenterWidth-wallRadius, wallLeftCenterY, wallRadius, space))            // left center wall
	*walls = append(*walls, newWall(wallRightCenterX, wallRightCenterY, wallRightCenterX+wallRightCenterWidth-wallRadius, wallRightCenterY, wallRadius, space))       // right center wall
	*walls = append(*walls, newWall(wallTopCenterX, wallTopCenterY, wallTopCenterX, wallTopCenterY+wallTopCenterHeight-wallRadius, wallRadius, space))                // top center wall
	*walls = append(*walls, newWall(wallBottomCenterX, wallBottomCenterY, wallBottomCenterX, wallBottomCenterY+wallBottomCenterHeight-wallRadius, wallRadius, space)) // bottom center wall
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	g.space.Step(deltaTime)

	// Update input states(mouse pos and pressed keys)
	g.input.update()

	// Raycast
	gunRay := g.player.gunRay
	var info cp.SegmentQueryInfo
	var success bool
	g.rayHitInfo.Alpha = 1.5
	for _, wall := range g.walls {
		success = wall.shape.SegmentQuery(gunRay[0], gunRay[1], 0, &info)
		if success && info.Alpha < g.rayHitInfo.Alpha {
			g.rayHitInfo = info
		}
	}

	// Update player and player's gun
	g.player.update(&g.input, &g.rayHitInfo)

	// Update Crosshair geometry matrix
	drawOptionsCursor.GeoM.Reset()
	drawOptionsCursor.GeoM.Translate(g.input.cursorPos.X-crosshairRadius, g.input.cursorPos.Y-crosshairRadius)
	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// Draw walls
	for _, wall := range g.walls {
		wall.draw(screen)
	}

	// Draw player and its gun
	g.player.draw(screen)

	// Draw crosshair
	screen.DrawImage(imageCursor, &drawOptionsCursor)

	// Draw rayhit
	const rayHitImageRadius = rayHitImageWidth / 2.0
	drawOptionsRayHit.GeoM.Reset()
	drawOptionsRayHit.GeoM.Translate(g.rayHitInfo.Point.X-rayHitImageRadius, g.rayHitInfo.Point.Y-rayHitImageRadius)
	screen.DrawImage(imageRayHit, &drawOptionsRayHit)

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
