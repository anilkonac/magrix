// Copyright 2022 Anıl Konaç

package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const (
	playerWidth  = 40.0
	playerHeight = 100.0
	gunWidth     = playerHeight / 2.0
	gunHeight    = playerWidth / 3.0
)

const (
	gravity          = 1000.0
	playerMass       = 1
	playerElasticity = 0.0
	// Taken from cp-examples/player and modified
	playerFriction        = playerGroundAccel / (2 * gravity)
	playerVelocity        = 400.0
	playerGroundAccelTime = 0.1
	playerGroundAccel     = playerVelocity / playerGroundAccelTime
	playerAirAccelTime    = 0.5
	playerAirAccel        = playerVelocity / playerAirAccelTime
	jumpHeight            = 60.0
	//
)

const (
	gunRange     = screenWidth + screenHeight
	gunForceMult = 45
	gunForceMax  = 1500
	gunMinAlpha  = 1e-5 // required to prevent player pos to go NaN
)

type gunState uint8

const (
	gunStateIdle gunState = iota
	gunStateAttract
	gunStateRepel
)

var (
	imagePlayer     = ebiten.NewImage(1, 1)
	imageGun        = ebiten.NewImage(1, 1)
	imageGunAttract = ebiten.NewImage(1, 1)
	imageGunRepel   = ebiten.NewImage(1, 1)
)

func init() {
	imagePlayer.Fill(colorPlayer)
	imageGun.Fill(colorGun)
	imageGunAttract.Fill(colorGunAttract)
	imageGunRepel.Fill(colorGunRepel)
}

type player struct {
	pos            cp.Vector
	posGun         cp.Vector
	angleGun       float64
	shape          *cp.Shape
	body           *cp.Body
	drawOptions    ebiten.DrawImageOptions
	drawOptionsGun ebiten.DrawImageOptions
	onGround       bool
	gunRay         [2]cp.Vector
	gunForce       cp.Vector
	stateGun       gunState
}

func newPlayer(pos cp.Vector, space *cp.Space) *player {
	player := &player{
		pos: pos,
	}

	player.body = cp.NewBody(playerMass, cp.INFINITY)
	player.body.SetPosition(cp.Vector{X: pos.X, Y: pos.Y})
	player.body.SetVelocityUpdateFunc(playerUpdateVelocity)
	player.shape = cp.NewBox(player.body, playerWidth, playerHeight, 0)
	player.shape.SetElasticity(playerElasticity)

	space.AddBody(player.body)
	space.AddShape(player.shape)

	return player
}

func (p *player) update(input *input, rayHitInfo *cp.SegmentQueryInfo) {
	// Update position
	p.pos = p.body.Position()

	// Update gun position
	p.posGun = p.pos

	// Update gun angle
	distX := input.cursorPos.X - p.posGun.X
	distY := input.cursorPos.Y - p.posGun.Y
	p.angleGun = math.Atan2(distY, distX)

	p.checkOnGround()

	// Raycast
	const rayLength = gunRange
	p.gunRay[0] = p.posGun
	p.gunRay[1] = p.gunRay[0].Add(cp.Vector{
		X: rayLength * math.Cos(p.angleGun), Y: rayLength * math.Sin(p.angleGun),
	})

	p.handleInputs(input, rayHitInfo)

	// v := p.body.Velocity()
	// fmt.Printf("Friction: %.2f\tVel X: %.2f\tVel Y: %.2f\n", p.shape.Friction(), v.X, v.Y)
	p.updateGeometryMatrices()
}

func (p *player) checkOnGround() {
	// Grab the grounding normal from last frame - Taken from cp-examples/player
	groundNormal := cp.Vector{}
	p.body.EachArbiter(func(arb *cp.Arbiter) {
		n := arb.Normal() //.Neg()

		if n.Y > groundNormal.Y {
			groundNormal = n
		}
	})
	p.onGround = groundNormal.Y > 0
}

func (p *player) handleInputs(input *input, rayHitInfo *cp.SegmentQueryInfo) {
	// Handle inputs
	var surfaceV cp.Vector
	if input.right {
		surfaceV.X = -playerVelocity
	} else if input.left {
		surfaceV.X = playerVelocity
	}
	p.shape.SetSurfaceV(surfaceV)
	if p.onGround {
		p.shape.SetFriction(playerFriction)
	} else {
		p.shape.SetFriction(0)
	}

	if input.up && p.onGround {
		jumpV := math.Sqrt(2.0 * jumpHeight * gravity) // Taken from cp-examples/player
		p.body.SetVelocityVector(p.body.Velocity().Add(cp.Vector{X: 0, Y: -jumpV}))
	}
	// Apply air control if not on ground
	if !p.onGround {
		v := p.body.Velocity()
		newVelX := cp.Clamp(v.X-surfaceV.X*deltaTime, -playerVelocity, playerVelocity)
		p.body.SetVelocity(newVelX, v.Y)
	}

	// Apply magnetic force if fire is pressed
	p.gunForce = cp.Vector{}
	if (input.attract || input.repel) && rayHitInfo.Alpha >= gunMinAlpha {
		forceDirection := rayHitInfo.Point.Sub(p.pos).Normalize()
		p.gunForce = forceDirection.Mult(gunForceMult).Mult(1 / (rayHitInfo.Alpha * rayHitInfo.Alpha))
		p.gunForce = p.gunForce.Clamp(gunForceMax)

		p.stateGun = gunStateAttract
		if input.repel {
			p.gunForce = p.gunForce.Neg()
			p.stateGun = gunStateRepel
		}
		p.body.SetForce(p.gunForce)
	} else {
		p.stateGun = gunStateIdle
	}
	// fmt.Printf("Player X: %.2f\tY:%.2f\tForce X: %.2f\tY:%.2f\n", p.pos.X, p.pos.Y, force.X, force.Y)
}

func (p *player) updateGeometryMatrices() {
	// Player
	p.drawOptions.GeoM.Reset()
	p.drawOptions.GeoM.Scale(playerWidth, playerHeight)
	p.drawOptions.GeoM.Translate(p.pos.X-playerWidth/2.0, p.pos.Y-playerHeight/2.0)

	// Gun
	p.drawOptionsGun.GeoM.Reset()
	p.drawOptionsGun.GeoM.Scale(gunWidth, gunHeight)
	p.drawOptionsGun.GeoM.Translate(0, -gunHeight/2.0)
	p.drawOptionsGun.GeoM.Rotate(p.angleGun)
	p.drawOptionsGun.GeoM.Translate(p.posGun.X, p.posGun.Y)
}

func (p *player) draw(dst *ebiten.Image) {
	// Draw prototype player
	dst.DrawImage(imagePlayer, &p.drawOptions)

	// Draw prototype gun
	if p.stateGun == gunStateAttract {
		dst.DrawImage(imageGunAttract, &p.drawOptionsGun)
	} else if p.stateGun == gunStateRepel {
		dst.DrawImage(imageGunRepel, &p.drawOptionsGun)
	} else {
		dst.DrawImage(imageGun, &p.drawOptionsGun)
	}

	// ebitenutil.DrawLine(dst, p.gunRay[0].X, p.gunRay[0].Y, p.gunRay[1].X, p.gunRay[1].Y, colorCrosshair)
}

func playerUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}