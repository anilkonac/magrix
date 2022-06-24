// Copyright 2022 Anıl Konaç

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const (
	cameraWidth  = 320
	cameraHeight = 240
	screenWidth  = 960
	screenHeight = 720
	deltaTime    = 1.0 / 60.0
)

const (
	crosshairRadius      = 6
	crosshairInnerRadius = 2
	rayHitImageWidth     = 4
	wallElasticity       = 0
	wallFriction         = 1
	// spaceIterations      = 10
)

var (
	colorBackground = color.RGBA{38, 38, 38, 255}
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
	tileLength         float64
)

const mapPath = "assets/testLevel.tmx"

var (
	imagePlatforms   *ebiten.Image
	imageComputers   *ebiten.Image
	imageDecorations *ebiten.Image
)

var gamePaused bool

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	initCursorImage()
	initRayHitImages()
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

func initRayHitImages() {
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

// game implements ebiten.game interface.
type game struct {
	player     player
	enemies    []*enemy
	walls      []*cp.Shape
	space      *cp.Space
	input      input
	rayHitInfo cp.SegmentQueryInfo
}

func newGame() *game {
	space := cp.NewSpace()
	// space.Iterations = spaceIterations
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	// Parse map file
	gameMap, err := tiled.LoadFile(mapPath)
	panicErr(err)
	tileLength = float64(gameMap.TileWidth)

	game := &game{
		player: *newPlayer(cp.Vector{X: cameraWidth / 2.0, Y: cameraHeight / 2.0}, space),
		space:  space,
	}

	game.loadMap(gameMap)

	return game
}

func (g *game) loadMap(gameMap *tiled.Map) {
	const (
		objectGroupWalls = 0
		objectGroupEnemy = 1
	)
	// Add enemies
	for _, enemyPos := range gameMap.ObjectGroups[objectGroupEnemy].Objects {
		g.enemies = append(g.enemies, newEnemy(cp.Vector{X: enemyPos.X, Y: enemyPos.Y}, g.space))

	}

	g.addWalls(gameMap.ObjectGroups[objectGroupWalls].Objects)

	// Render layer images
	renderer, err := render.NewRenderer(gameMap)
	panicErr(err)

	err = renderer.RenderLayer(0)
	panicErr(err)
	imageDecorations = ebiten.NewImageFromImage(renderer.Result)

	renderer.Clear()
	err = renderer.RenderLayer(1)
	panicErr(err)
	imageComputers = ebiten.NewImageFromImage(renderer.Result)

	renderer.Clear()
	err = renderer.RenderLayer(2)
	panicErr(err)
	imagePlatforms = ebiten.NewImageFromImage(renderer.Result)
}

func (g *game) addWalls(wallObjects []*tiled.Object) {
	for _, obj := range wallObjects {
		radius := math.Min(obj.Width, obj.Height) / 2.0
		x2 := obj.X + obj.Width - radius
		y2 := obj.Y + obj.Height - radius
		shape := g.space.AddShape(cp.NewSegment(g.space.StaticBody, cp.Vector{X: obj.X + radius, Y: obj.Y + radius}, cp.Vector{X: x2, Y: y2}, radius))
		shape.SetElasticity(wallElasticity)
		shape.SetFriction(wallFriction)

		g.walls = append(g.walls, shape)
	}
}

// Update is called every tick (1/60 [s] by default).
func (g *game) Update() error {
	g.input.update()
	drawOptionsCursor.GeoM.Reset()
	drawOptionsCursor.GeoM.Translate(g.input.cursorPos.X-crosshairRadius, g.input.cursorPos.Y-crosshairRadius)

	g.updateSettings()

	if gamePaused {
		return nil
	}

	g.space.Step(deltaTime)

	g.rayCast()

	// Update player and player's gun
	g.player.update(&g.input, &g.rayHitInfo)

	// Send negation of the player's gun force to the enemy
	var force cp.Vector
	for _, enemy := range g.enemies {
		if g.rayHitInfo.Shape == enemy.shape {
			force = g.player.gunForce.Neg()
			enemy.update(&force)
		} else {
			enemy.update(nil)
		}
	}

	// Update geometry matrices
	const rayHitImageRadius = rayHitImageWidth / 2.0

	drawOptionsRayHit.GeoM.Reset()
	drawOptionsRayHit.GeoM.Translate(g.rayHitInfo.Point.X-rayHitImageRadius, g.rayHitInfo.Point.Y-rayHitImageRadius)

	return nil
}

func (g *game) updateSettings() {
	// Escape from cursor captured mode
	if g.input.escape {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	} else if (ebiten.CursorMode() == ebiten.CursorModeHidden) && g.input.repel {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	}

	if g.input.pausePlay {
		gamePaused = !gamePaused
	}
}

func (g *game) rayCast() {
	gunRay := g.player.gunRay
	var info cp.SegmentQueryInfo
	var success bool
	g.rayHitInfo.Alpha = 1.5

	// Check walls
	for _, shape := range g.walls {
		success = shape.SegmentQuery(gunRay[0], gunRay[1], 0, &info)
		if success && info.Alpha < g.rayHitInfo.Alpha {
			g.rayHitInfo = info
		}
	}

	// Check enemy
	for _, enemy := range g.enemies {
		success = enemy.shape.SegmentQuery(gunRay[0], gunRay[1], 0, &info)
		if success && info.Alpha < g.rayHitInfo.Alpha {
			g.rayHitInfo = info
		}
	}

}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// Draw decorations
	screen.DrawImage(imageDecorations, &ebiten.DrawImageOptions{})
	screen.DrawImage(imageComputers, &ebiten.DrawImageOptions{})

	// Draw player and its gun
	g.player.draw(screen)

	// Draw enemies
	for _, enemy := range g.enemies {
		enemy.draw(screen)
	}

	// Draw walls and platforms
	screen.DrawImage(imagePlatforms, &ebiten.DrawImageOptions{})

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
	return cameraWidth, cameraHeight
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
