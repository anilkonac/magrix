package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

const (
	rocketMass               = 0.25
	rocketMoment             = 10
	rocketVelocity           = 120
	rocketWidth              = 8
	rocketHeight             = 2
	explosionTotalDurationMs = durationExplosionMs * 14
)

var (
	imageRocket    = ebiten.NewImage(8, 3)
	imageExplosion = ebiten.NewImage(32, 32)
)

type explosion struct {
	pos             cp.Vector
	drawOptions     ebiten.DrawImageOptions
	drawOptionsAnim ganim8.DrawOptions
	elapsedMs       int64
	animation       ganim8.Animation
}

func newExplosion(pos cp.Vector) *explosion {
	explo := &explosion{
		pos: pos,
		drawOptionsAnim: ganim8.DrawOptions{
			// X:       pos.X,
			// Y:       pos.Y,
			ScaleX:  1.0,
			ScaleY:  1.0,
			OriginX: 0.0,
			OriginY: 0.0,
		},
		animation: *animExplosion,
	}

	return explo
}

type rocket struct {
	body            *cp.Body
	shape           *cp.Shape
	drawOptions     ebiten.DrawImageOptions
	drawOptionsAnim ganim8.DrawOptions
}

func newRocket(startPos cp.Vector, angle float64, space *cp.Space) *rocket {
	body := cp.NewBody(rocketMass, rocketMoment)
	body.SetPosition(startPos)
	body.SetVelocityUpdateFunc(rocketUpdateVelocity)
	body.SetAngle(angle)
	body.SetVelocity(rocketVelocity*math.Cos(angle), 0)

	shape := cp.NewBox(body, rocketWidth, rocketHeight, 0)
	// TODO: Set elasticity and friction ?

	space.AddBody(body)
	space.AddShape(shape)

	drawOptsAnim := ganim8.DrawOptions{
		// X:       startPos.X,
		// Y:       startPos.Y,
		ScaleX:  1.0,
		ScaleY:  1.0,
		OriginX: 0.0,
		OriginY: 0.0,
	}

	return &rocket{
		body:            body,
		shape:           shape,
		drawOptionsAnim: drawOptsAnim,
	}
}

type rocketManager struct {
	rockets    []*rocket
	explosions []*explosion
	space      *cp.Space
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

		// Eliminate gravity
		// velocityPercent := rocket.body.Velocity().Length() / rocketVelocity // To eliminate floating stopped rockets
		rocket.body.SetForce(cp.Vector{X: 0, Y: -gravity * rocketMass /* * velocityPercent*/})

		// Update draw options
		pos := rocket.body.Position()
		rocket.drawOptions.GeoM.Reset()
		// rocket.drawOptions.GeoM.Translate(-4, -1.5)
		rocket.drawOptions.GeoM.Rotate(rocket.body.Angle())
		rocket.drawOptions.GeoM.Concat(cam.GetTranslation(pos.X, pos.Y).GeoM)
	}

	// TODO: Object pooling?
	// Delete hit rockets
	for iRocket, rocket := range m.rockets {
		for _, rocketTarget := range rocketsToBeDeleted {
			if rocket == rocketTarget {
				// Delete from slice
				copy(m.rockets[iRocket:], m.rockets[iRocket+1:])
				m.rockets[len(m.rockets)-1] = nil
				m.rockets = m.rockets[:len(m.rockets)-1]

				// Delete from space
				m.space.RemoveShape(rocket.shape)
				m.space.RemoveBody(rocket.body)

				rocketTarget = nil
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
			continue
		}
		explo.drawOptions = *cam.GetTranslation(explo.pos.X-16, explo.pos.Y-16)
	}

	// Delete ended explosion animation
	for iExplo, explo := range m.explosions {
		for _, exploTarget := range explosionsToBeDeleted {
			if explo == exploTarget {
				copy(m.explosions[iExplo:], m.explosions[iExplo+1:])
				m.explosions[len(m.explosions)-1] = nil
				m.explosions = m.explosions[:len(m.explosions)-1]

				exploTarget = nil
			}
		}

	}

	return
}

func (m *rocketManager) draw() {
	// Draw rockets
	for _, rocket := range m.rockets {
		// TODO: Fix not drawing rocket
		// imageRocket.Clear()
		imageRocket.Fill(colorGunAttract)
		animRocket.Draw(imageRocket, &rocket.drawOptionsAnim)
		cam.Surface.DrawImage(imageRocket, &rocket.drawOptions)
	}

	// Draw explosions
	for _, explo := range m.explosions {
		imageExplosion.Clear()
		explo.animation.Draw(imageExplosion, &explo.drawOptionsAnim)
		cam.Surface.DrawImage(imageExplosion, &explo.drawOptions)
	}
}

func rocketUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
