// Copyright 2022 Anıl Konaç

package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
)

const (
	playerWidth  = 40.0
	playerHeight = 100.0
	gunWidth     = playerHeight / 2.0
	gunHeight    = playerWidth / 3.0
)

var (
	imagePlayer = ebiten.NewImage(1, 1)
	imageGun    = ebiten.NewImage(1, 1)
)

func init() {
	imagePlayer.Fill(colornames.Slategray)
	imageGun.Fill(colornames.Orange)
}

type player struct {
	pos               cp.Vector
	posGun            cp.Vector
	angleGun          float64
	shape             *cp.Shape
	drawOptionsPlayer ebiten.DrawImageOptions
	drawOptionsGun    ebiten.DrawImageOptions
}

func newPlayer(pos cp.Vector) *player {
	player := &player{
		pos: pos,
	}

	body := cp.NewBody(1, cp.INFINITY)
	body.SetPosition(pos)
	player.shape = cp.NewBox(body, playerWidth, playerHeight, playerHeight/2.0)
	// TODO: Set elasticity and friction

	return player
}

func (p *player) update(cursorPos *cp.Vector) {
	// Update position
	p.pos = p.shape.Body().Position()

	// Update gun position
	p.posGun = p.pos.Add(cp.Vector{X: playerWidth / 2.0, Y: playerHeight / 2.0})

	// Update gun angle
	distX := cursorPos.X - p.posGun.X
	distY := cursorPos.Y - p.posGun.Y
	p.angleGun = math.Atan2(distY, distX)

	p.updateGeometryMatrices()
}

func (p *player) updateGeometryMatrices() {
	// Player
	p.drawOptionsPlayer.GeoM.Reset()
	p.drawOptionsPlayer.GeoM.Scale(playerWidth, playerHeight)
	p.drawOptionsPlayer.GeoM.Translate(p.pos.X, p.pos.Y)

	// Gun
	p.drawOptionsGun.GeoM.Reset()
	p.drawOptionsGun.GeoM.Scale(gunWidth, gunHeight)
	p.drawOptionsGun.GeoM.Translate(0, -gunHeight/2.0)
	p.drawOptionsGun.GeoM.Rotate(p.angleGun)
	p.drawOptionsGun.GeoM.Translate(p.posGun.X, p.posGun.Y)
}

func (p *player) draw(dst *ebiten.Image) {
	// Draw prototype player
	dst.DrawImage(imagePlayer, &p.drawOptionsPlayer)

	// Draw prototype gun
	dst.DrawImage(imageGun, &p.drawOptionsGun)
}
