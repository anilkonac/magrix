// Copyright 2022 Anıl Konaç

package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

const (
	enemyMass              = 0.75
	enemyFriction          = 0.75
	enemyMoment            = 50
	enemyWidthTile         = 1
	enemyHeightTile        = 1.5
	enemyEyeRange          = cameraWidth
	enemyEyeRadius         = cameraHeight / 4.0
	enemyAttackCooldownSec = 2.0
)

var imageEnemy = ebiten.NewImage(16, 32)

type enemy struct {
	size              cp.Vector
	drawOptions       ebiten.DrawImageOptions
	drawOptionsAnim   ganim8.DrawOptions
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
		drawOptionsAnim: ganim8.DrawOptions{
			OriginX: 0.0,
			OriginY: 0.15,
			ScaleX:  1.00,
			ScaleY:  1.00,
		},
		curAnim:           *animEnemy1Idle,
		attackCooldownSec: 0,
		turnedLeft:        turnedLeft,
	}
	if turnedLeft {
		enemy.drawOptionsAnim.ScaleX = -1.0
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

	// Update animation
	e.curAnim.Update(animDeltaTime)
	e.drawOptions.GeoM.Reset()
	e.drawOptions.GeoM.Translate(-tileLength/2.0, -tileLength)
	e.drawOptions.GeoM.Rotate(e.body.Angle())
	e.drawOptions.GeoM.Concat(cam.GetTranslation(pos.X, pos.Y).GeoM)
	if e.turnedLeft {
		e.drawOptionsAnim.ScaleX = -1.0
		e.drawOptionsAnim.OriginX = 1.0
	} else {
		e.drawOptionsAnim.ScaleX = 1.0
		e.drawOptionsAnim.OriginX = 0.0
	}
}

func (e *enemy) draw() {
	// e.curAnim.Draw(dst, &e.drawOptions)
	imageEnemy.Clear()
	e.curAnim.Draw(imageEnemy, &e.drawOptionsAnim)
	cam.Surface.DrawImage(imageEnemy, &e.drawOptions)
}

func enemyUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}

// Broken
// // Goroutine
// func (e *enemy) standUpBot() {
// 	const standUpForceY = -8000
// 	const standUpAngularVelocity = -4
// 	const checkIntervalSec = 3.0
// 	const checkEpsilon = 1.0

// 	vForce := cp.Vector{X: 0, Y: standUpForceY}

// 	ticker := time.NewTicker(time.Second * checkIntervalSec)
// 	for range ticker.C {
// 		if gamePaused {
// 			continue
// 		}
// 		angleDeg := e.body.Angle() * cp.DegreeConst
// 		angleDegMod := math.Mod(angleDeg, 180)

// 		if math.Abs(angleDegMod-90) < checkEpsilon {
// 			e.body.SetAngularVelocity(standUpAngularVelocity)
// 			e.body.SetForce(vForce)
// 		} else if math.Abs(angleDegMod+90) < checkEpsilon {
// 			e.body.SetAngularVelocity(-standUpAngularVelocity)
// 			e.body.SetForce(vForce)
// 		}
// 	}
// }
