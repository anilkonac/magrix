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
	gridWidth, gridHeight = 16, 32
	enemy1IdleDurationMs  = 200
)

var (
	animDeltaTime = time.Duration(math.Ceil(deltaTime * 1000))
)

var (
	//go:embed assets/enemy1_idle.png
	enemy1IdleBytes []byte
	//go:embed assets/player_idle.png
	playerIdleBytes []byte
	//go:embed assets/player_walk.png
	playerWalkBytes []byte
)

var (
	animEnemy1Idle *ganim8.Animation
	animPlayerIdle *ganim8.Animation
	animPlayerWalk *ganim8.Animation
)

func init() {
	animPlayerIdle = newAnim("1-4", 1, playerIdleBytes, 64, 32, enemy1IdleDurationMs)
	animPlayerWalk = newAnim("1-8", 1, playerWalkBytes, 64, 32, enemy1IdleDurationMs)
	animEnemy1Idle = newAnim("1-4", 1, enemy1IdleBytes, 64, 32, enemy1IdleDurationMs)
}

func newAnim(column string, row int, fileBytes []byte, imageWidth, imageHeight, frameDurationMs int) *ganim8.Animation {
	img, err := png.Decode(bytes.NewReader(fileBytes))
	panicErr(err)
	image := ebiten.NewImageFromImage(img)

	grid := ganim8.NewGrid(gridWidth, gridHeight, imageWidth, imageHeight)
	frames := grid.GetFrames(column, row)
	spr := ganim8.NewSprite(image, frames)
	return ganim8.NewAnimation(spr, time.Millisecond*time.Duration(frameDurationMs), ganim8.Nop)
}
