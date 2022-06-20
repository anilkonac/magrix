package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const enemyFriction = 1.25

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

	body := cp.NewBody(playerMass, cp.INFINITY)
	body.SetPosition(cp.Vector{X: pos.X, Y: pos.Y})
	body.SetVelocityUpdateFunc(enemyUpdateVelocity)
	enemy.body = body
	enemy.shape = cp.NewBox(enemy.body, playerWidth, playerHeight, 0)
	enemy.shape.SetElasticity(playerElasticity)
	enemy.shape.SetFriction(enemyFriction)

	space.AddBody(enemy.body)
	space.AddShape(enemy.shape)

	return enemy
}

func (e *enemy) update(force *cp.Vector) {
	e.pos = e.body.Position()

	if force != nil {
		e.body.SetForce(*force)
	}

	e.drawOptions.GeoM.Reset()
	e.drawOptions.GeoM.Scale(playerWidth, playerHeight)
	e.drawOptions.GeoM.Translate(e.pos.X-playerWidth/2.0, e.pos.Y-playerHeight/2.0)

	f := e.shape.Friction()
	fmt.Printf("Enemy friction: %.2f\n", f)
}

func (e *enemy) draw(dst *ebiten.Image) {
	dst.DrawImage(imageGunActive, &e.drawOptions)
}

func enemyUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
