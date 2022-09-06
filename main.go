// Copyright 2022 Anıl Konaç

package main

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"time"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/yohamta/ganim8/v2"
)

const (
	screenWidth    = 960
	screenHeight   = 720
	deltaTimeSec   = 1.0 / 60.0
	mapWidth       = 960
	mapHeight      = 960
	tileLength     = 16
	restartTimeSec = 3
)

const (
	crosshairRadius      = 6
	crosshairInnerRadius = 2
	rayHitImageWidth     = 4
	wallElasticity       = 0
	wallFriction         = 1
	// spaceIterations      = 10
)

const mapPath = "gameMap.tmx"

const (
	zoomMultiplier        = 0.1
	uiArrowDistance       = screenHeight/2.0 - 50
	interactionRadiusTile = 1.25
)

var (
	colorBackground = color.RGBA{38, 38, 38, 255}
	colorOrange     = color.RGBA{253, 147, 89, 255}
	colorBlue       = color.RGBA{111, 215, 231, 255}
	colorGunAttract = color.RGBA{216, 17, 89, 255} // ~ Ruby
	colorGunRepel   = color.RGBA{80, 142, 237, 255}
	colorPlayer     = color.RGBA{155, 201, 149, 255} // ~ Dark Sea Green
	colorCrosshair  = color.RGBA{255, 251, 255, 255} // ~ Snow
	colorEnemy      = color.RGBA{165, 1, 4, 255}     // ~ Rufous
	colorGreen      = color.RGBA{0, 231, 86, 255}
)

var (
	//go:embed circle.go
	bytesCircleShader      []byte
	imageCrosshair         = ebiten.NewImage(crosshairRadius*2, crosshairRadius*2)
	imageRayHit            = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	imageRayHitAttract     = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	imageRayHitRepel       = ebiten.NewImage(rayHitImageWidth, rayHitImageWidth)
	imageArrowBlue         = ebiten.NewImage(tileLength, tileLength)
	imageArrowOrange       = ebiten.NewImage(tileLength, tileLength)
	drawOptionsCursor      ebiten.DrawImageOptions
	drawOptionsRayHit      ebiten.DrawImageOptions
	drawOptionsZero        ebiten.DrawImageOptions
	drawOptionsArrowBlue   ganim8.DrawOptions
	drawOptionsArrowOrange ganim8.DrawOptions
	drawOptionsArrowGreen  ganim8.DrawOptions
	//go:embed assets/heart.png
	bytesHeart []byte
)

var (
	imagePlatforms   *ebiten.Image
	imageDecorations *ebiten.Image
	imageObjects     *ebiten.Image = ebiten.NewImage(mapWidth, mapHeight)
)

var (
	cam              = camera.NewCamera(screenWidth, screenHeight, 0, 0, 0, 1)
	cursorX, cursorY float64
	zoom             = 3.5
	zoomMin          = 1.0
	zoomMax          = 6.0
)

var (
	gamePaused      bool
	gameOver        bool
	showArrowBlue   bool
	showArrowOrange bool
)

//go:embed gameMap.tmx
var bytesMap []byte

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	initCursorImage()
	initRayHitImages()
	imageArrowBlue.Fill(colorBlue)
	imageArrowOrange.Fill(colorOrange)
}

func initCursorImage() {
	ebitenutil.DrawLine(imageCrosshair, 0, crosshairRadius,
		crosshairRadius-crosshairInnerRadius, crosshairRadius, colorCrosshair)
	ebitenutil.DrawLine(imageCrosshair, crosshairRadius, 0,
		crosshairRadius, crosshairRadius-crosshairInnerRadius, colorCrosshair)
	ebitenutil.DrawLine(imageCrosshair, crosshairRadius+crosshairInnerRadius,
		crosshairRadius, 2*crosshairRadius, crosshairRadius, colorCrosshair)
	ebitenutil.DrawLine(imageCrosshair, crosshairRadius, crosshairRadius+crosshairInnerRadius,
		crosshairRadius, 2*crosshairRadius, colorCrosshair)
}

func initRayHitImages() {
	shader, err := ebiten.NewShader(bytesCircleShader)
	if err != nil {
		panic(err)
	}

	// Prepare ray hit images (circle)
	imageRayHit.DrawRectShader(rayHitImageWidth, rayHitImageWidth, shader, &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"Radius": float32(rayHitImageWidth / 2.0),
			"Color":  []float32{float32(colorOrange.R) / 255.0, float32(colorOrange.G) / 255.0, float32(colorOrange.B) / 255.0, float32(colorOrange.A) / 255.0},
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
	player         player
	enemies        []*enemy
	walls          []*cp.Shape
	space          *cp.Space
	input          input
	rayHitInfo     cp.SegmentQueryInfo
	rocketManager  rocketManager
	terminalBlue   *terminal
	terminalOrange *terminal
	terminalIntro  *terminal
	eWallBlue      *electricWall
	eWallOrange    *electricWall
	button         *button
	gameOverTimer  float32
}

func newGame() *game {
	space := cp.NewSpace()
	// space.Iterations = spaceIterations
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	// Parse map file
	// gameMap, err := tiled.LoadFile(mapPath)
	gameMap, err := tiled.LoadReader("", bytes.NewReader(bytesMap))
	panicErr(err)
	// tileLength = float64(gameMap.TileWidth)

	game := &game{
		space: space,
		rocketManager: rocketManager{
			space: space,
		},
	}
	cam.Zoom(zoom)

	game.loadMap(gameMap)

	drawOptionsArrowBlue.ScaleX = 2.0
	drawOptionsArrowBlue.ScaleY = 2.0
	drawOptionsArrowOrange.ScaleX = 2.0
	drawOptionsArrowOrange.ScaleY = 2.0
	drawOptionsArrowGreen.ScaleX = 2.0
	drawOptionsArrowGreen.ScaleY = 2.0

	return game
}

func (g *game) restart() {
	space := cp.NewSpace()
	// space.Iterations = spaceIterations
	space.SetGravity(cp.Vector{X: 0, Y: gravity})

	// Parse map file
	// gameMap, err := tiled.LoadFile(mapPath)
	gameMap, err := tiled.LoadReader("", bytes.NewReader(bytesMap))
	panicErr(err)

	*g = game{
		space: space,
		rocketManager: rocketManager{
			space: space,
		},
	}
	cam.Zoom(zoom)

	g.loadMap(gameMap)

	gameOver = false
	showArrowBlue = false
	showArrowOrange = false
}

func (g *game) loadMap(gameMap *tiled.Map) {
	const (
		objectGroupWalls         = 0
		objectGroupPlayer        = 1
		objectGroupEnemy         = 2
		objectGroupWallsElectric = 3
		objectGroupTerminals     = 4
		objectGroupButton        = 5
	)

	g.addWalls(gameMap.ObjectGroups[objectGroupWalls].Objects)

	// Add Electric Walls (Manual assignment for now)
	g.eWallBlue = newElectricWall(gameMap.ObjectGroups[objectGroupWallsElectric].Objects[0], g.space)
	g.eWallOrange = newElectricWall(gameMap.ObjectGroups[objectGroupWallsElectric].Objects[1], g.space)

	// Add terminals
	g.terminalIntro = newTerminal(gameMap.ObjectGroups[objectGroupTerminals].Objects[2], g.space)
	g.terminalIntro.spr = spriteTerminalGreen
	g.terminalBlue = newTerminal(gameMap.ObjectGroups[objectGroupTerminals].Objects[0], g.space)
	g.terminalOrange = newTerminal(gameMap.ObjectGroups[objectGroupTerminals].Objects[1], g.space)

	var playerStartLoc cp.Vector
	playerStartLoc.X = gameMap.ObjectGroups[objectGroupPlayer].Objects[0].X
	playerStartLoc.Y = gameMap.ObjectGroups[objectGroupPlayer].Objects[0].Y
	g.player = *newPlayer(playerStartLoc, g.space)

	// Add enemies
	for _, enemyPos := range gameMap.ObjectGroups[objectGroupEnemy].Objects {
		g.enemies = append(g.enemies, newEnemy(cp.Vector{X: enemyPos.X, Y: enemyPos.Y}, g.space, enemyPos.Properties.GetBool("turnedLeft")))

	}

	// Add the button
	g.button = newButton(gameMap.ObjectGroups[objectGroupButton].Objects[0], g.space)

	const (
		layerPlatform    = 1
		layerDecorations = 0
	)

	// Render layer images
	renderer, err := render.NewRenderer(gameMap)
	panicErr(err)

	err = renderer.RenderLayer(layerPlatform)
	panicErr(err)
	imagePlatforms = ebiten.NewImageFromImage(renderer.Result)

	renderer.Clear()
	err = renderer.RenderLayer(layerDecorations)
	panicErr(err)
	imageDecorations = ebiten.NewImageFromImage(renderer.Result)

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
	if g.input.wheelDy > 0 {
		zoom += zoomMultiplier
	} else if g.input.wheelDy < 0 {
		zoom -= zoomMultiplier
	}
	zoom = cp.Clamp(zoom, zoomMin, zoomMax)
	cam.SetZoom(zoom)
	cursorX, cursorY = cam.GetCursorCoords()
	drawOptionsCursor.GeoM.Reset()
	cam.GetTranslation(&drawOptionsCursor, cursorX-crosshairRadius, cursorY-crosshairRadius)

	g.updateSettings()

	if gamePaused {
		return nil
	}

	g.space.Step(deltaTimeSec)

	g.rayCast()
	hitBodies := g.rocketManager.update()
	for _, hitBody := range hitBodies {
		if hitBody == g.player.body {
			g.player.hit()
			if g.player.numLives <= 0 {
				gameOver = true
			}
		} else {
			for _, enemy := range g.enemies {
				if hitBody == enemy.body && enemy.isAlive {
					enemy.isAlive = false
					g.killEnemy(enemy)
				}
			}
		}
	}

	if gameOver {
		g.gameOverTimer += deltaTimeSec
		if g.gameOverTimer >= restartTimeSec {
			g.restart()
			gameOver = false
			return nil
		}
	}

	// Update player and player's gun
	g.player.update(&g.input, &g.rayHitInfo)
	cam.SetPosition(g.player.pos.X, g.player.pos.Y)

	// Send the negative of the player's gun force to the enemy
	var force cp.Vector
	var enemyFell bool
	for _, enemy := range g.enemies {
		if g.rayHitInfo.Shape == enemy.shape {
			force = g.player.gunForce.Neg()
			enemyFell = enemy.update(&force)
		} else {
			enemyFell = enemy.update(nil)
		}

		if enemyFell {
			g.killEnemy(enemy)
		}
	}
	// Send the negative of the player's gun force to the rocket
	if g.input.attract || g.input.repel {
		for _, rocket := range g.rocketManager.rockets {
			if g.rayHitInfo.Shape == rocket.shape {
				force = g.player.gunForce.Neg()
				rocket.body.SetForce(force)
			}
		}
	}

	g.checkPlayerInteraction()

	// Update ewall animations
	if g.eWallBlue != nil {
		g.eWallBlue.update()
	}
	if g.eWallOrange != nil {
		g.eWallOrange.update()
	}

	// Update draw options
	// -------------------
	const rayHitImageRadius = rayHitImageWidth / 2.0

	drawOptionsZero.GeoM.Reset()
	cam.GetTranslation(&drawOptionsZero, 0, 0)
	drawOptionsRayHit.GeoM.Reset()
	cam.GetTranslation(&drawOptionsRayHit, g.rayHitInfo.Point.X-rayHitImageRadius, g.rayHitInfo.Point.Y-rayHitImageRadius)

	// Blue arrow
	direction := g.terminalBlue.pos.Sub(g.player.pos)
	distance := direction.LengthSq()
	if distance < 10*screenWidth {
		showArrowBlue = false
	} else if g.terminalIntro.triggered && !g.terminalBlue.triggered {
		showArrowBlue = true
		dirAngle := math.Atan2(direction.Y, direction.X)
		drawOptionsArrowBlue.Rotate = dirAngle
		drawOptionsArrowBlue.X = screenWidth/2.0 + uiArrowDistance*math.Cos(dirAngle)
		drawOptionsArrowBlue.Y = screenHeight/2.0 + uiArrowDistance*math.Sin(dirAngle)
	}

	// Orange arrow
	direction = g.terminalOrange.pos.Sub(g.player.pos)
	distance = direction.LengthSq()
	if distance < 10*screenWidth {
		showArrowOrange = false
	} else if g.terminalIntro.triggered && !g.terminalOrange.triggered {
		showArrowOrange = true
		dirAngle := math.Atan2(direction.Y, direction.X)
		drawOptionsArrowOrange.Rotate = dirAngle
		drawOptionsArrowOrange.X = screenWidth/2.0 + uiArrowDistance*math.Cos(dirAngle)
		drawOptionsArrowOrange.Y = screenHeight/2.0 + uiArrowDistance*math.Sin(dirAngle)
	}

	return nil
}

// Goroutine
func (g *game) killEnemy(e *enemy) {
	e.isAlive = false

	go func() {
		ticker := time.NewTicker(time.Second * 2.0)
		<-ticker.C
		g.rocketManager.explosions = append(g.rocketManager.explosions, newExplosion(e.body.Position()))
		<-ticker.C

		// Delete the enemy
		// ----------------
		g.space.RemoveShape(e.shape)
		g.space.RemoveBody(e.body)

		// copy(g.enemies[iEnemy:], g.enemies[iEnemy+1:])
		// g.enemies[len(g.enemies)-1] = nil
		// g.enemies = g.enemies[:len(g.enemies)-1]
		// --

		e.drawActive = false

		ticker.Stop()
	}()
}

func (g *game) checkPlayerInteraction() {
	if !g.input.activate {
		return
	}

	interactionRadius := float64(interactionRadiusTile * tileLength)
	// Check if near intro terminal
	if g.terminalIntro.pos.Distance(g.player.pos) < interactionRadius {
		showTextIntro = true

		go func() {
			duration := time.Duration(1000 * durationTextIntroSec)
			timer := time.NewTimer(time.Millisecond * duration)
			<-timer.C
			showTextIntro = false
			showArrowBlue = true
			showArrowOrange = true
			g.terminalIntro.trigger()
		}()
	}

	// Check if near blue terminal
	if !g.terminalBlue.triggered && (g.terminalBlue.pos.Distance(g.player.pos) < interactionRadius) {
		g.terminalBlue.trigger()

		showTextTerminalBlue = true
		go func() {
			timer := time.NewTimer(time.Second * durationTextTerminals)
			<-timer.C
			showTextTerminalBlue = false
			g.player.numLives++
			g.player.prepareLivesIndicator()

			// Remove wall
			g.space.RemoveShape(g.eWallBlue.shape)
			g.space.RemoveBody(g.eWallBlue.shape.Body())
			g.eWallBlue = nil

		}()

	}

	// Check if near orange terminal
	if !g.terminalOrange.triggered && (g.terminalOrange.pos.Distance(g.player.pos) < interactionRadius) {
		g.terminalOrange.trigger()

		showTextTerminalOrange = true
		go func() {
			timer := time.NewTimer(time.Second * durationTextTerminals)
			<-timer.C
			showTextTerminalOrange = false
			g.player.numLives++
			g.player.prepareLivesIndicator()

			// Remove wall
			g.space.RemoveShape(g.eWallOrange.shape)
			g.space.RemoveBody(g.eWallOrange.shape.Body())
			g.eWallOrange = nil

		}()

	}

	// Check if near the button
	if !g.button.triggered && (g.button.pos.Distance(g.player.pos) < interactionRadius) {
		g.button.trigger()
		showTextButton = true
	}
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
		if gamePaused && (musicState == musicOn) {
			playerMusic.Pause()
			musicState = musicPaused
		} else if musicState == musicPaused {
			playerMusic.Play()
			musicState = musicOn
		}
	}

	if g.input.musicToggle {
		if musicState == musicOn {
			musicState = musicMuted
			playerMusic.Pause()
		} else if (musicState == musicMuted) && !gamePaused {
			musicState = musicOn
			playerMusic.Play()
		}
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

	// Check rockets
	for _, rocket := range g.rocketManager.rockets {
		success = rocket.shape.SegmentQuery(gunRay[0], gunRay[1], 0, &info)
		if success && info.Alpha < g.rayHitInfo.Alpha {
			g.rayHitInfo = info
		}
	}

	// Check player
	// for enemy to detect player
	for _, enemy := range g.enemies {
		if !enemy.isAlive {
			continue
		}
		success = g.player.shape.SegmentQuery(enemy.eyeRay[0], enemy.eyeRay[1], enemyEyeRadius, &info)
		if success && enemy.attackCooldownSec <= 0 {
			enemyPos := enemy.body.Position()
			enemyAngle := enemy.body.Angle()
			var rocketSpawnPos cp.Vector
			var rocketAngle float64
			if !enemy.turnedLeft {
				rocketSpawnPosHypot := math.Hypot(rocketSpawnPosRelative.X, rocketSpawnPosRelative.Y)
				angleSpawnPos := math.Atan2(rocketSpawnPosRelative.Y, rocketSpawnPosRelative.X)
				newSpawnPosAngle := angleSpawnPos + enemyAngle
				rocketSpawnPos = enemyPos.Add(cp.Vector{
					X: rocketSpawnPosHypot * math.Cos(newSpawnPosAngle), Y: rocketSpawnPosHypot * math.Sin(newSpawnPosAngle),
				})
				rocketAngle = enemyAngle
			} else {
				rocketSpawnPosRelativeLeft := cp.Vector{X: -rocketSpawnPosRelative.X, Y: rocketSpawnPosRelative.Y}

				rocketSpawnPosHypot := math.Hypot(rocketSpawnPosRelativeLeft.X, rocketSpawnPosRelativeLeft.Y)
				angleSpawnPos := math.Atan2(rocketSpawnPosRelativeLeft.Y, rocketSpawnPosRelativeLeft.X)
				newSpawnPosAngle := angleSpawnPos + enemyAngle
				rocketSpawnPos = enemyPos.Add(cp.Vector{
					X: rocketSpawnPosHypot * math.Cos(newSpawnPosAngle), Y: rocketSpawnPosHypot * math.Sin(newSpawnPosAngle),
				})
				rocketAngle = enemyAngle - math.Pi
			}
			g.rocketManager.rockets = append(g.rocketManager.rockets, newRocket(
				rocketSpawnPos, rocketAngle, g.space))
			enemy.attackCooldownSec = enemyAttackCooldownSec
		} else {
			enemy.attackCooldownSec -= deltaTimeSec
		}
		// fmt.Printf("success: %v\n", success)
	}

}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *game) Draw(screen *ebiten.Image) {
	// screen.Fill(colorBackground)
	imageObjects.Clear()
	cam.Surface.Fill(colorBackground)

	// Draw decorations
	cam.Surface.DrawImage(imageDecorations, &drawOptionsZero)

	// Draw terminals
	g.terminalIntro.draw()
	g.terminalBlue.draw()
	g.terminalOrange.draw()

	// Draw enemies
	for _, enemy := range g.enemies {
		enemy.draw()
	}

	// Draw rockets
	g.rocketManager.draw()

	// Draw electric walls
	if g.eWallBlue != nil {
		g.eWallBlue.draw()
	}
	if g.eWallOrange != nil {
		g.eWallOrange.draw()
	}

	// Draw the button
	g.button.draw()

	cam.Surface.DrawImage(imageObjects, &drawOptionsZero)

	// Draw player and its gun
	g.player.draw()

	// Draw walls and platforms
	cam.Surface.DrawImage(imagePlatforms, &drawOptionsZero)

	// Draw crosshair
	cam.Surface.DrawImage(imageCrosshair, &drawOptionsCursor)

	// Draw rayhit
	var imageHit *ebiten.Image
	if g.input.attract {
		imageHit = imageRayHitAttract
	} else if g.input.repel {
		imageHit = imageRayHitRepel
	} else {
		imageHit = imageRayHit
	}
	cam.Surface.DrawImage(imageHit, &drawOptionsRayHit)

	cam.Blit(screen)

	if showTextIntro {
		screen.DrawImage(imageTextIntro, &drawOptionsTextIntro)
	}

	if showTextTerminalBlue {
		screen.DrawImage(imageTextTerminalBlue, &drawOptionsTextTerminalBlue)
	}

	if showTextTerminalOrange {
		screen.DrawImage(imageTextTerminalOrange, &drawOptionsTextTerminalOrange)
	}

	if showTextButton {
		screen.DrawImage(imageTextButton, &drawOptionsTextButton)
	}

	if gameOver {
		screen.DrawImage(imageTextFail, &drawOptionsTextFail)
	}

	if showArrowBlue {
		spriteArrows.Draw(screen, 1, &drawOptionsArrowBlue)

	}

	if showArrowOrange {
		spriteArrows.Draw(screen, 0, &drawOptionsArrowOrange)
	}

	// Draw hearts
	screen.DrawImage(imageLives, &drawOptionsLives)

	// Print fps
	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f  FPS: %.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()), screenWidth-140, 0)
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
