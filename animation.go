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

var (
	animDeltaTime = time.Duration(math.Ceil(deltaTime * 1000))
)

var (
	//go:embed assets/enemy1_idle.png
	enemy1IdleBytes []byte
	//go:embed assets/player_idle.png
	playerIdleBytes []byte
)

var (
	enemy1IdleAnim *ganim8.Animation
	playerIdleAnim *ganim8.Animation
)

func init() {
	initPlayerIdle()
	initEnemy1Idle()
}

func initPlayerIdle() {
	// Prepare player idle anim
	const column string = "1-4"
	const row int = 1
	img, err := png.Decode(bytes.NewReader(playerIdleBytes))
	panicErr(err)
	playerIdleImage := ebiten.NewImageFromImage(img)

	grid := ganim8.NewGrid(gridWidth, gridHeight, 64, 32)
	frames := grid.GetFrames(column, row)
	spr := ganim8.NewSprite(playerIdleImage, frames)
	playerIdleAnim = ganim8.NewAnimation(spr, time.Millisecond*enemy1IdleDurationMs, ganim8.Nop)
}

func initEnemy1Idle() {
	// Prepare enemy 1 idle anim
	const column string = "1-4"
	const row int = 1
	img, err := png.Decode(bytes.NewReader(enemy1IdleBytes))
	panicErr(err)
	enemy1IdleImage := ebiten.NewImageFromImage(img)

	grid := ganim8.NewGrid(gridWidth, gridHeight, 64, 32)
	frames := grid.GetFrames(column, row)
	spr := ganim8.NewSprite(enemy1IdleImage, frames)
	enemy1IdleAnim = ganim8.NewAnimation(spr, time.Millisecond*enemy1IdleDurationMs, ganim8.Nop)
}
