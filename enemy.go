// Copyright 2022 Anıl Konaç

package main

import (
	"math"
	"math/rand"

	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

const (
	enemyMass              = 0.75
	enemyFriction          = 0.75
	enemyMoment            = 50
	enemyWidthTile         = 1
	enemyHeightTile        = 1.5
	enemyEyeRange          = mapWidth / 2.0
	enemyEyeRadius         = mapHeight / 8.0
	enemyAttackCooldownSec = 2.0
)

type enemy struct {
	size              cp.Vector
	drawOptions       ganim8.DrawOptions
	body              *cp.Body
	shape             *cp.Shape
	curAnim           ganim8.Animation
	eyeRay            [2]cp.Vector
	attackCooldownSec float32
	turnedLeft        bool
}

func newEnemy(pos cp.Vector, space *cp.Space, turnedLeft bool) *enemy {
	enemy := &enemy{
		size: cp.Vector{
			X: enemyWidthTile * tileLength,
			Y: enemyHeightTile * tileLength},
		drawOptions: ganim8.DrawOptions{
			OriginX: 0.5,
			OriginY: 0.64,
			ScaleX:  1.0,
			ScaleY:  1.0,
		},
		curAnim:           *animEnemy1Idle,
		attackCooldownSec: 0,
		turnedLeft:        turnedLeft,
	}
	if turnedLeft {
		enemy.drawOptions.ScaleX = -1.0
	}

	enemy.curAnim.GoToFrame(rand.Intn(4)) // Have all enemies start at different frames

	body := cp.NewBody(enemyMass, enemyMoment)
	body.SetPosition(cp.Vector{X: pos.X, Y: pos.Y})
	body.SetVelocityUpdateFunc(enemyUpdateVelocity)
	enemy.body = body

	enemy.shape = cp.NewBox(enemy.body, enemy.size.X, enemy.size.Y, 0)
	enemy.shape.SetElasticity(playerElasticity)
	enemy.shape.SetFriction(enemyFriction)

	space.AddBody(enemy.body)
	space.AddShape(enemy.shape)

	// go enemy.standUpBot()

	return enemy
}

func (e *enemy) update(force *cp.Vector) {
	pos := e.body.Position()

	if force != nil {
		e.body.SetForce(*force)
	}

	// Raycast
	angle := e.body.Angle()
	turnMult := 1.0
	if e.turnedLeft {
		turnMult = -1.0
	}
	e.eyeRay[0] = pos
	e.eyeRay[1] = e.eyeRay[0].Add(
		cp.Vector{
			X: enemyEyeRange * turnMult * math.Cos(angle), Y: enemyEyeRange * math.Sin(angle),
		},
	)

	// Update draw options
	e.curAnim.Update(animDeltaTime)
	e.drawOptions.X = pos.X
	e.drawOptions.Y = pos.Y
	e.drawOptions.Rotate = e.body.Angle()
}

func (e *enemy) draw() {
	e.curAnim.Draw(imageObjects, &e.drawOptions)
}

func enemyUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
