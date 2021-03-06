package main

import (
	"bytes"
	"image/png"
	"math"
	"time"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8/v2"
)

const (
	gridWidth, gridHeight   = 16, 32
	durationEnemy1IdleMs    = 150
	durationPlayerIdleMs    = 200
	durationPlayerWalkingMs = 75
	durationExplosionMs     = 50
	durationElectricMs      = 60
	numFramesExplosion      = 14
)

var (
	animDeltaTime = time.Millisecond * time.Duration(math.Ceil(deltaTimeSec*1000))
)

var (
	//go:embed assets/enemy1_idle.png
	byteseEemy1Idle []byte
	//go:embed assets/player_idle.png
	bytesPlayerIdle []byte
	//go:embed assets/player_walk.png
	bytesPlayerWalk []byte
	//go:embed assets/rocket_anim.png
	bytesRocket []byte
	//go:embed assets/Explosion_duplicateframes.png
	bytesExplosion []byte
	//go:embed assets/anim_electric_blue.png
	bytesElectricBlue []byte
	//go:embed assets/anim_electric_orange.png
	bytesElectricOrange []byte
	//go:embed assets/terminal_blue.png
	bytesTerminalBlue []byte
	//go:embed assets/terminal_orange.png
	bytesTerminalOrange []byte
	//go:embed assets/terminal_green.png
	bytesTerminalGreen []byte
	//go:embed assets/arrows.png
	bytesArrows []byte
	//go:embed assets/theButton.png
	bytesButton []byte
)

var (
	animEnemy1Idle       *ganim8.Animation
	animPlayerIdle       *ganim8.Animation
	animPlayerWalk       *ganim8.Animation
	animRocket           *ganim8.Animation
	animExplosion        *ganim8.Animation
	animElectricBlue     *ganim8.Animation
	animElectricOrange   *ganim8.Animation
	spriteTerminalBlue   *ganim8.Sprite
	spriteTerminalOrange *ganim8.Sprite
	spriteTerminalGreen  *ganim8.Sprite
	spriteArrows         *ganim8.Sprite
	spriteButton         *ganim8.Sprite
)

func init() {
	animPlayerIdle = newAnim("1-4", 1, bytesPlayerIdle, gridWidth, gridHeight, 64, 32, durationPlayerIdleMs)
	animPlayerWalk = newAnim("1-8", 1, bytesPlayerWalk, gridWidth, gridHeight, 128, 32, durationPlayerWalkingMs)
	animEnemy1Idle = newAnim("1-4", 1, byteseEemy1Idle, gridWidth, gridHeight, 64, 32, durationEnemy1IdleMs)
	animRocket = newAnim("1-2", 1, bytesRocket, 16, 16, 32, 16, 50)
	animExplosion = newAnim("1-14", 1, bytesExplosion, 32, 32, 32*numFramesExplosion, 32, durationExplosionMs)
	animElectricBlue = newAnim("1-3", 1, bytesElectricBlue, 16, 16, 48, 16, durationElectricMs)
	animElectricOrange = newAnim("1-3", 1, bytesElectricOrange, 16, 16, 48, 16, durationElectricMs)
	spriteTerminalBlue = newSprite("1-2", 1, bytesTerminalBlue, 16, 32, 32, 32)
	spriteTerminalOrange = newSprite("1-2", 1, bytesTerminalOrange, 16, 32, 32, 32)
	spriteTerminalGreen = newSprite("1-2", 1, bytesTerminalGreen, 16, 32, 32, 32)
	spriteArrows = newSprite("1-3", 1, bytesArrows, tileLength, tileLength, tileLength*3, tileLength)
	spriteButton = newSprite("1-2", 1, bytesButton, tileLength, tileLength, 32, 16)
}

func newAnim(column string, row int, fileBytes []byte, gridWidth, gridHeight, imageWidth, imageHeight, frameDurationMs int) *ganim8.Animation {
	spr := newSprite(column, row, fileBytes, gridWidth, gridHeight, imageWidth, imageHeight)
	return ganim8.NewAnimation(spr, time.Millisecond*time.Duration(frameDurationMs), ganim8.Nop)
}

func newSprite(column string, row int, fileBytes []byte, gridWidth, gridHeight, imageWidth, imageHeight int) *ganim8.Sprite {
	img, err := png.Decode(bytes.NewReader(fileBytes))
	panicErr(err)
	image := ebiten.NewImageFromImage(img)

	grid := ganim8.NewGrid(gridWidth, gridHeight, imageWidth, imageHeight)
	frames := grid.GetFrames(column, row)
	return ganim8.NewSprite(image, frames)
}
