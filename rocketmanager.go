package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

const (
	rocketMass     = 1e-10
	rocketMoment   = 100
	rocketVelocity = 60
	rocketWidth    = 8
	rocketHeight   = 2
)

type rocket struct {
	body        *cp.Body
	shape       *cp.Shape
	target      cp.Vector
	drawOptions ganim8.DrawOptions
}

func newRocket(startPos, target cp.Vector, angle float64, space *cp.Space) *rocket {
	body := cp.NewBody(rocketMass, rocketMoment)
	body.SetPosition(startPos)
	body.SetVelocityUpdateFunc(rocketUpdateVelocity)
	body.SetVelocity(rocketVelocity, 0)
	// body.Set

	shape := cp.NewBox(body, rocketWidth, rocketHeight, 0)
	// TODO: Set elastic and frict ?

	space.AddBody(body)
	space.AddShape(shape)

	drawOpts := ganim8.DrawOptions{
		X:       startPos.X,
		Y:       startPos.Y,
		ScaleX:  1.0,
		ScaleY:  1.0,
		OriginX: 0.5,
		OriginY: 0.5,
	}

	return &rocket{body, shape, target, drawOpts}
}

type rocketManager struct {
	rockets []*rocket
}

func (m *rocketManager) update() {
	animRocket.Update(animDeltaTime)
	for _, rocket := range m.rockets {
		rocket.body.EachArbiter(func(arb *cp.Arbiter) {
			count := arb.Count()
			fmt.Printf("count: %v\n", count)
		})

		// Update position
		pos := rocket.body.Position()
		rocket.drawOptions.X = pos.X
		rocket.drawOptions.Y = pos.Y

		// Update velocity
		vel := rocket.body.Velocity()
		rocket.body.SetVelocity(vel.X, 0)
	}
}

func (m *rocketManager) draw(dst *ebiten.Image) {
	for _, rocket := range m.rockets {
		animRocket.Draw(dst, &rocket.drawOptions)
	}
}

func rocketUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
