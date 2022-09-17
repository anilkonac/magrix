package main

import (
	"bytes"
	"image/png"
	"math"
	"time"

	_ "embed"

	"github.com/anilkonac/magrix/asset"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8/v2"
)

const (
	durationEnemy1IdleMs    = 150
	durationPlayerIdleMs    = 200
	durationPlayerWalkingMs = 75
	durationExplosionMs     = 50
	durationElectricMs      = 60
	numFramesExplosion      = 14

	gridWidthGun, gridHeightGun       = 8, 3
	gridWidthPlayer, gridHeightPlayer = 16, 32
)

var (
	animDeltaTime = time.Millisecond * time.Duration(math.Ceil(deltaTimeSec*1000))

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
	spriteButton         *ganim8.Sprite
	spriteGun            ganim8.Sprite
)

func init() {
	animPlayerIdle = newAnim("1-4", 1, asset.Bytes(asset.AnimPlayerIdle), gridWidthPlayer, gridHeightPlayer, 64, 32, durationPlayerIdleMs)
	animPlayerWalk = newAnim("1-8", 1, asset.Bytes(asset.AnimPlayerWalk), gridWidthPlayer, gridHeightPlayer, 128, 32, durationPlayerWalkingMs)
	animEnemy1Idle = newAnim("1-4", 1, asset.Bytes(asset.AnimEnemy1Idle), gridWidthPlayer, gridHeightPlayer, 64, 32, durationEnemy1IdleMs)
	animRocket = newAnim("1-2", 1, asset.Bytes(asset.AnimRocket), 16, 16, 32, 16, 50)
	animExplosion = newAnim("1-14", 1, asset.Bytes(asset.AnimExplosion), 32, 32, 32*numFramesExplosion, 32, durationExplosionMs)
	animElectricBlue = newAnim("1-3", 1, asset.Bytes(asset.AnimElectricBlue), 16, 16, 48, 16, durationElectricMs)
	animElectricOrange = newAnim("1-3", 1, asset.Bytes(asset.AnimElectricOrange), 16, 16, 48, 16, durationElectricMs)
	spriteTerminalBlue = newSprite("1-2", 1, asset.Bytes(asset.SpriteTerminalBlue), 16, 32, 32, 32)
	spriteTerminalOrange = newSprite("1-2", 1, asset.Bytes(asset.SpriteTerminalOrange), 16, 32, 32, 32)
	spriteTerminalGreen = newSprite("1-2", 1, asset.Bytes(asset.SpriteTerminalGreen), 16, 32, 32, 32)
	spriteButton = newSprite("1-2", 1, asset.Bytes(asset.SpriteButton), tileLength, tileLength, 32, 16)
	spriteGun = *newSprite("1-3", 1, asset.Bytes(asset.SpriteGun), gridWidthGun, gridHeightGun, 3*gridWidthGun, gridHeightGun)
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
