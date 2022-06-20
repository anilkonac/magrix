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
	crosshairRadius      = 14
	crosshairInnerRadius = 4
	rayHitImageWidth     = 16
	// spaceIterations      = 10
)

var (
	colorBackground = color.RGBA{124, 144, 160, 255} // ~ Light Slate Gray
	colorWall       = color.RGBA{57, 62, 65, 255}    // ~ Onyx
	colorGun        = color.RGBA{242, 129, 35, 255}  // ~ Princeton Orange
	colorGunAttract = color.RGBA{216, 17, 89, 255}   // ~ Ruby
	colorGunRepel   = color.RGBA{7, 160, 195, 255}   // ~ Blue Green
	colorPlayer     = color.RGBA{155, 201, 149, 255} // ~ Dark Sea Green
	colorCrosshair  = color.RGBA{255, 251, 255, 255} // ~ Snow
	colorEnemy      = color.RGBA{165, 1, 4, 255}     // ~ Rufous
)

var (
	imageCursor        = ebiten.NewImage(crosshairRadius*2, crosshairRadius*2)
	imageRayHit        = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	imageRayHitAttract = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	imageRayHitRepel   = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	drawOptionsCursor  ebiten.DrawImageOptions
	drawOptionsRayHit  ebiten.DrawImageOptions
)

func init() {
	initCursorImage()

	shader, err := ebiten.NewShader(circle_go)
	if err != nil {
		panic(err)
	}

	// Prepare ray hit images (circle)
	imageRayHit.DrawRectShader(rayHitImageWidth, rayHitImageWidth, shader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius": float32(rayHitImageWidth / 2.0),
			"Color":  []float32{float32(colorGun.R) / 255.0, float32(colorGun.G) / 255.0, float32(colorGun.B) / 255.0, float32(colorGun.A) / 255.0},
		},
	})
	imageRayHitAttract.DrawRectShader(rayHitImageWidth, rayHitImageWidth, shader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius": float32(rayHitImageWidth / 2.0),
			"Color":  []float32{float32(colorGunAttract.R) / 255.0, float32(colorGunAttract.G) / 255.0, float32(colorGunAttract.B) / 255.0, float32(colorGunAttract.A) / 255.0},
		},
	})
	imageRayHitRepel.DrawRectShader(rayHitImageWidth, rayHitImageWidth, shader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius": float32(rayHitImageWidth / 2.0),
			"Color":  []float32{float32(colorGunRepel.R) / 255.0, float32(colorGunRepel.G) / 255.0, float32(colorGunRepel.B) / 255.0, float32(colorGunRepel.A) / 255.0},
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
	enemy      enemy
	walls      []*wall
	space      *cp.Space
	input      input
	rayHitInfo cp.SegmentQueryInfo
}

func newGame() *game {
	space := cp.NewSpace()
	// space.Iterations = spaceIterations
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	game := &game{
		player: *newPlayer(cp.Vector{X: screenWidth / 2.0, Y: screenHeight / 2.0}, space),
		enemy:  *newEnemy(cp.Vector{X: 778, Y: 149}, space),
		space:  space,
	}

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

	// Escape from cursor captured mode
	if g.input.escape {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	} else if (ebiten.CursorMode() == ebiten.CursorModeHidden) && g.input.attract {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	}

	g.rayCast()

	// Update player and player's gun
	g.player.update(&g.input, &g.rayHitInfo)

	// Sen negation of the player's gun force
	var force cp.Vector
	if g.rayHitInfo.Shape == g.enemy.shape {
		force = g.player.gunForce.Neg()
		g.enemy.update(&force)
	} else {
		g.enemy.update(nil)
	}

	// Update geometry matrices
	const rayHitImageRadius = rayHitImageWidth / 2.0
	drawOptionsCursor.GeoM.Reset()
	drawOptionsCursor.GeoM.Translate(g.input.cursorPos.X-crosshairRadius, g.input.cursorPos.Y-crosshairRadius)
	drawOptionsRayHit.GeoM.Reset()
	drawOptionsRayHit.GeoM.Translate(g.rayHitInfo.Point.X-rayHitImageRadius, g.rayHitInfo.Point.Y-rayHitImageRadius)

	return nil
}

func (g *game) rayCast() {
	gunRay := g.player.gunRay
	var info cp.SegmentQueryInfo
	var success bool
	g.rayHitInfo.Alpha = 1.5

	// Check wall
	for _, wall := range g.walls {
		success = wall.shape.SegmentQuery(gunRay[0], gunRay[1], 0, &info)
		if success && info.Alpha < g.rayHitInfo.Alpha {
			g.rayHitInfo = info
		}
	}

	// Check enemy
	success = g.enemy.shape.SegmentQuery(gunRay[0], gunRay[1], 0, &info)
	if success && info.Alpha < g.rayHitInfo.Alpha {
		g.rayHitInfo = info
	}
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

	// Draw enemy
	g.enemy.draw(screen)

	// Draw crosshair
	screen.DrawImage(imageCursor, &drawOptionsCursor)

	// Draw rayhit
	var imageHit *ebiten.Image
	if g.input.attract {
		imageHit = imageRayHitAttract
	} else if g.input.repel {
		imageHit = imageRayHitRepel
	} else {
		imageHit = imageRayHit
	}
	screen.DrawImage(imageHit, &drawOptionsRayHit)

	// Print fps
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f  FPS: %.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %.0f, Y: %.0f", g.input.cursorPos.X, g.input.cursorPos.Y), 0, 15)
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
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
