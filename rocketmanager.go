package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

const (
	rocketMass               = 0.1
	rocketMoment             = 100
	rocketVelocity           = 120
	rocketWidth              = 8
	rocketHeight             = 2
	explosionTotalDurationMs = durationExplosionMs * 14
)

type explosion struct {
	drawOptions ganim8.DrawOptions
	elapsedMs   int64
	animation   ganim8.Animation
}

func newExplosion(pos cp.Vector) *explosion {
	return &explosion{
		drawOptions: ganim8.DrawOptions{
			X:       pos.X,
			Y:       pos.Y,
			ScaleX:  1.0,
			ScaleY:  1.0,
			OriginX: 0.5,
			OriginY: 0.5,
		},
		animation: *animExplosion,
	}
}

type rocket struct {
	body        *cp.Body
	shape       *cp.Shape
	drawOptions ganim8.DrawOptions
}

func newRocket(startPos cp.Vector, angle float64, space *cp.Space) *rocket {
	body := cp.NewBody(rocketMass, rocketMoment)
	body.SetPosition(startPos)
	body.SetVelocityUpdateFunc(rocketUpdateVelocity)
	body.SetAngle(angle)
	velMult := 1.0
	scaleX := 1.0
	if angle <= -halfPi || angle >= halfPi {
		velMult = -1.0
		scaleX = -1.0
	}
	body.SetVelocity(rocketVelocity*velMult, 0)
	// body.Set

	shape := cp.NewBox(body, rocketWidth, rocketHeight, 0)
	// TODO: Set elastic and frict ?

	space.AddBody(body)
	space.AddShape(shape)

	drawOpts := ganim8.DrawOptions{
		X:       startPos.X,
		Y:       startPos.Y,
		ScaleX:  scaleX,
		ScaleY:  1.0,
		OriginX: 0.5,
		OriginY: 0.5,
	}

	return &rocket{body, shape, drawOpts}
}

type rocketManager struct {
	rockets    []*rocket
	explosions []*explosion
}

func (m *rocketManager) update( /*playerPos *cp.Vector*/ ) (hitBodies []*cp.Body) {
	animRocket.Update(animDeltaTime)

	rocketsToBeDeleted := make([]*rocket, 0, 8)
	for _, rocket := range m.rockets {
		var rocketHit bool
		var hitBody *cp.Body
		rocket.body.EachArbiter(func(arb *cp.Arbiter) {
			if arb.IsFirstContact() {
				bodyA, bodyB := arb.Bodies()
				if bodyA != rocket.body {
					hitBody = bodyA
				} else {
					hitBody = bodyB
				}
				rocketHit = true
			}

		})

		if rocketHit {
			m.explosions = append(m.explosions, newExplosion(rocket.body.Position()))
			hitBodies = append(hitBodies, hitBody)
			rocketsToBeDeleted = append(rocketsToBeDeleted, rocket)
			continue
		}

		// Update position
		pos := rocket.body.Position()
		rocket.drawOptions.X = pos.X
		rocket.drawOptions.Y = pos.Y

		// Eliminate gravity
		// velocityPercent := rocket.body.Velocity().Length() / rocketVelocity // To eliminate floating stopped rockets
		rocket.body.SetForce(cp.Vector{X: 0, Y: -gravity * rocketMass /* * velocityPercent*/})
	}

	// // TODO: Object pooling?
	// // Delete hit rockets
	for iRocket, rocket := range m.rockets {
		for _, rocketTarget := range rocketsToBeDeleted {
			if rocket == rocketTarget {
				copy(m.rockets[iRocket:], m.rockets[iRocket+1:])
				m.rockets[len(m.rockets)-1] = nil
				m.rockets = m.rockets[:len(m.rockets)-1]
			}
		}
	}

	// Update explosion animations
	explosionsToBeDeleted := make([]*explosion, 0, 8)
	for _, explo := range m.explosions {
		explo.animation.Update(animDeltaTime)
		explo.elapsedMs += animDeltaTime.Milliseconds()
		if explo.elapsedMs >= explosionTotalDurationMs {
			explosionsToBeDeleted = append(explosionsToBeDeleted, explo)
		}
	}

	// Delete ended explosion animation
	for iExplo, explo := range m.explosions {
		for _, exploTarget := range explosionsToBeDeleted {
			if explo == exploTarget {
				copy(m.explosions[iExplo:], m.explosions[iExplo+1:])
				m.explosions[len(m.explosions)-1] = nil
				m.explosions = m.explosions[:len(m.explosions)-1]
			}
		}

	}

	return
}

func (m *rocketManager) draw(dst *ebiten.Image) {
	// Draw rockets
	for _, rocket := range m.rockets {
		animRocket.Draw(dst, &rocket.drawOptions)
	}

	// Draw explosions
	for _, explo := range m.explosions {
		explo.animation.Draw(dst, &explo.drawOptions)

	}
}

func rocketUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
