// Copyright 2022 Anıl Konaç

package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const (
	enemyMass       = 0.75
	enemyFriction   = 0.5
	enemyMoment     = 50
	enemyWidthTile  = 1
	enemyHeightTile = 2
)

var (
	enemyImage = ebiten.NewImage(1, 1)
)

func init() {
	enemyImage.Fill(colorEnemy)
}

type enemy struct {
	pos         cp.Vector
	size        cp.Vector
	drawOptions ebiten.DrawImageOptions
	body        *cp.Body
	shape       *cp.Shape
}

func newEnemy(pos cp.Vector, space *cp.Space) *enemy {
	enemy := &enemy{
		pos: pos,
		size: cp.Vector{
			X: enemyWidthTile * tileLength,
			Y: enemyHeightTile * tileLength},
	}

	body := cp.NewBody(enemyMass, enemyMoment)
	body.SetPosition(cp.Vector{X: pos.X, Y: pos.Y})
	body.SetVelocityUpdateFunc(enemyUpdateVelocity)
	enemy.body = body
	enemy.shape = cp.NewBox(enemy.body, enemy.size.X, enemy.size.Y, 0)
	enemy.shape.SetElasticity(playerElasticity)
	enemy.shape.SetFriction(enemyFriction)

	space.AddBody(enemy.body)
	space.AddShape(enemy.shape)

	go enemy.standUpBot()

	return enemy
}

func (e *enemy) update(force *cp.Vector) {
	e.pos = e.body.Position()

	if force != nil {
		e.body.SetForce(*force)
	}

	angle := e.body.Angle()
	e.drawOptions.GeoM.Reset()
	e.drawOptions.GeoM.Scale(e.size.X, e.size.Y)
	e.drawOptions.GeoM.Translate(-e.size.X/2.0, -e.size.Y/2.0)
	e.drawOptions.GeoM.Rotate(angle)
	e.drawOptions.GeoM.Translate(e.pos.X, e.pos.Y)
}

func (e *enemy) draw(dst *ebiten.Image) {
	dst.DrawImage(enemyImage, &e.drawOptions)
}

func enemyUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}

// Goroutine
func (e *enemy) standUpBot() {
	const standUpForceY = -8000
	const standUpAngularVelocity = -4
	const checkIntervalSec = 3.0
	const checkEpsilon = 1.0

	vForce := cp.Vector{X: 0, Y: standUpForceY}

	ticker := time.NewTicker(time.Second * checkIntervalSec)
	for range ticker.C {
		if gamePaused {
			continue
		}
		angleDeg := e.body.Angle() * cp.DegreeConst
		angleDegMod := math.Mod(angleDeg, 180)
		// isNegative := math.Signbit(angleDeg)
		// fmt.Printf("angleDeg: %.2f\tisNegative: %v\n", angleDeg, isNegative)

		if math.Abs(angleDegMod-90) < checkEpsilon {
			e.body.SetAngularVelocity(standUpAngularVelocity)
			e.body.SetForce(vForce)
		} else if math.Abs(angleDegMod+90) < checkEpsilon {
			e.body.SetAngularVelocity(-standUpAngularVelocity)
			e.body.SetForce(vForce)
		}
	}
}
