package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const (
	enemyMass     = 2.0
	enemyFriction = 0.75
	enemyWidth    = playerWidth
	enemyHeight   = playerHeight
	enemyMoment   = 75
)

var (
	enemyImage = ebiten.NewImage(1, 1)
)

func init() {
	enemyImage.Fill(colorEnemy)
}

type enemy struct {
	pos         cp.Vector
	drawOptions ebiten.DrawImageOptions
	body        *cp.Body
	shape       *cp.Shape
}

func newEnemy(pos cp.Vector, space *cp.Space) *enemy {
	enemy := &enemy{
		pos: pos,
	}

	body := cp.NewBody(enemyMass, enemyMoment)
	body.SetPosition(cp.Vector{X: pos.X, Y: pos.Y})
	body.SetVelocityUpdateFunc(enemyUpdateVelocity)
	enemy.body = body
	enemy.shape = cp.NewBox(enemy.body, enemyWidth, enemyHeight, 0)
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
	e.drawOptions.GeoM.Scale(enemyWidth, enemyHeight)
	e.drawOptions.GeoM.Translate(-enemyWidth/2.0, -enemyHeight/2.0)
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
	const standUpForceY = -30000
	const standUpAngularVelocity = -35.0
	const checkIntervalSec = 2.0
	const checkEpsilon = 2.0
	ticker := time.NewTicker(time.Second * checkIntervalSec)
	for range ticker.C {
		angleDegMod := math.Mod(e.body.Angle()*cp.DegreeConst, 180)

		if math.Abs(angleDegMod-90) < checkEpsilon {
			e.body.SetAngularVelocity(-standUpAngularVelocity)
			e.body.SetForce(cp.Vector{X: 0, Y: standUpForceY})
		} else if math.Abs(angleDegMod+90) < checkEpsilon {
			e.body.SetAngularVelocity(standUpAngularVelocity)
			e.body.SetForce(cp.Vector{X: 0, Y: standUpForceY})
		}
	}
}
