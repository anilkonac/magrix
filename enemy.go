package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const (
	enemyFriction = 1.00
	enemyMass     = playerMass
	enemyWidth    = playerWidth
	enemyHeight   = playerHeight
	enemyMoment   = 100
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
