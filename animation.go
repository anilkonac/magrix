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
	durationEnemy1IdleMs    = 200
	durationPlayerWalkingMs = 75
)

var (
	animDeltaTime = time.Millisecond * time.Duration(math.Ceil(deltaTime*1000))
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
)

var (
	animEnemy1Idle *ganim8.Animation
	animPlayerIdle *ganim8.Animation
	animPlayerWalk *ganim8.Animation
	animRocket     *ganim8.Animation
)

func init() {
	animPlayerIdle = newAnim("1-4", 1, bytesPlayerIdle, gridWidth, gridHeight, 64, 32, durationEnemy1IdleMs)
	animPlayerWalk = newAnim("1-8", 1, bytesPlayerWalk, gridWidth, gridHeight, 128, 32, durationPlayerWalkingMs)
	animEnemy1Idle = newAnim("1-4", 1, byteseEemy1Idle, gridWidth, gridHeight, 64, 32, durationEnemy1IdleMs)
	animRocket = newAnim("1-2", 1, bytesRocket, 16, 16, 32, 16, 50)
}

func newAnim(column string, row int, fileBytes []byte, gridWidth, gridHeight, imageWidth, imageHeight, frameDurationMs int) *ganim8.Animation {
	img, err := png.Decode(bytes.NewReader(fileBytes))
	panicErr(err)
	image := ebiten.NewImageFromImage(img)

	grid := ganim8.NewGrid(gridWidth, gridHeight, imageWidth, imageHeight)
	frames := grid.GetFrames(column, row)
	spr := ganim8.NewSprite(image, frames)
	return ganim8.NewAnimation(spr, time.Millisecond*time.Duration(frameDurationMs), ganim8.Nop)
}
