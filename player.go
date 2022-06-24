// Copyright 2022 Anıl Konaç

package main

import (
	"bytes"
	"fmt"
	"image/png"
	"math"
	"time"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

const (
	playerWidthTile  = 1
	playerHeightTile = 1.6
	gunWidthTile     = 1
	gunHeightTile    = 1.0 / 3.0
)

const (
	gravity          = 750.0
	playerMass       = 0.75
	playerElasticity = 0.0
	// Taken from cp-examples/player and modified
	playerFriction        = playerGroundAccel / (4 * gravity)
	playerVelocity        = 150.0
	playerGroundAccelTime = 0.05
	playerGroundAccel     = playerVelocity / playerGroundAccelTime
	jumpHeightTile        = 1.5
	//
)

const (
	gunRange     = cameraWidth + cameraHeight
	gunForceMult = 15
	gunForceMax  = 750
	gunMinAlpha  = 1e-5 // required to prevent player pos to go NaN
)

type gunState uint8

const (
	gunStateIdle gunState = iota
	gunStateAttract
	gunStateRepel
)

var (
	//go:embed assets/gun_idle.png
	gunIdleBytes []byte
	//go:embed assets/gun_attract.png
	gunAttractBytes []byte
	//go:embed assets/gun_repel.png
	gunRepelBytes   []byte
	imageGunIdle    *ebiten.Image
	imageGunAttract *ebiten.Image
	imageGunRepel   *ebiten.Image
)

var posGunRelative cp.Vector

func init() {
	img, err := png.Decode(bytes.NewReader(gunIdleBytes))
	panicErr(err)
	imageGunIdle = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(gunAttractBytes))
	panicErr(err)
	imageGunAttract = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(gunRepelBytes))
	panicErr(err)
	imageGunRepel = ebiten.NewImageFromImage(img)
}

type player struct {
	pos            cp.Vector
	posGun         cp.Vector
	size           cp.Vector
	sizeGun        cp.Vector
	angleGun       float64
	shape          *cp.Shape
	body           *cp.Body
	drawOptions    ganim8.DrawOptions
	drawOptionsGun ebiten.DrawImageOptions
	curAnim        *ganim8.Animation
	onGround       bool
	gunRay         [2]cp.Vector
	gunForce       cp.Vector
	stateGun       gunState
}

func newPlayer(pos cp.Vector, space *cp.Space) *player {
	player := &player{
		pos: pos,
		size: cp.Vector{
			X: playerWidthTile * tileLength,
			Y: playerHeightTile * tileLength,
		},
		sizeGun: cp.Vector{
			X: gunWidthTile * tileLength,
			Y: gunHeightTile * tileLength,
		},
		drawOptions: ganim8.DrawOptions{
			OriginX: 0.5,
			OriginY: 0.6,
			ScaleX:  1.0,
			ScaleY:  1.0,
		},
		curAnim: animPlayerIdle,
	}

	player.body = cp.NewBody(playerMass, cp.INFINITY)
	player.body.SetPosition(cp.Vector{X: pos.X, Y: pos.Y})
	player.body.SetVelocityUpdateFunc(playerUpdateVelocity)
	player.shape = cp.NewBox(player.body, playerWidthTile*tileLength, playerHeightTile*tileLength, 0)
	player.shape.SetElasticity(playerElasticity)

	space.AddBody(player.body)
	space.AddShape(player.shape)

	posGunRelative = cp.Vector{X: tileLength / 7.0, Y: -tileLength / 4.0}

	return player
}

func (p *player) update(input *input, rayHitInfo *cp.SegmentQueryInfo) {
	// Update position
	p.pos = p.body.Position()

	// Update gun position
	if p.angleGun < -1.5 || p.angleGun > 1.5 {
		p.posGun = p.pos.Add(cp.Vector{-posGunRelative.X, posGunRelative.Y})
	} else {
		p.posGun = p.pos.Add(posGunRelative)
	}

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
	p.curAnim.Update(time.Millisecond * animDeltaTime)
	p.updateGeometryMatrices()

	// fmt.Printf("p.angleGun: %v\n", p.angleGun)
}

func (p *player) checkOnGround() {
	const groundNormalYThreshold = 0.5
	// Grab the grounding normal from last frame - Taken from cp-examples/player and modified
	groundNormal := cp.Vector{}
	p.body.EachArbiter(func(arb *cp.Arbiter) {
		n := arb.Normal() //.Neg()

		if n.Y > groundNormal.Y {
			groundNormal = n
		}
	})
	p.onGround = groundNormal.Y > groundNormalYThreshold
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
		jumpV := math.Sqrt(2.0 * jumpHeightTile * tileLength * gravity) // Taken from cp-examples/player
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

	// v := p.body.Velocity()
	// fmt.Printf("Velocity X: %.2f\tY: %.2f\t\tForce X: %.2f\tY:%.2f\n", v.X, v.Y, p.gunForce.X, p.gunForce.Y)
}

func (p *player) updateGeometryMatrices() {
	// Player
	p.drawOptions.X = p.pos.X
	p.drawOptions.Y = p.pos.Y
	if p.angleGun < -1.5 || p.angleGun > 1.5 {
		p.drawOptions.ScaleX = -1.0
	} else {
		p.drawOptions.ScaleX = 1.0
	}

	// Gun
	p.drawOptionsGun.GeoM.Reset()
	// TODO: Fix gun drawing when player turns to the left
	// if p.angleGun < -1.5 || p.angleGun > 1.5 {
	// 	p.drawOptionsGun.GeoM.Scale(0, -1.0)
	// }
	// angleDeg := p.angleGun * cp.DegreeConst
	// if angleDeg < -90.0 {
	// 	angleDeg = -180.0 - angleDeg
	// } else if angleDeg > 90 {
	// 	angleDeg = 180 - angleDeg
	// }
	// fmt.Printf("angleDeg: %v\n", angleDeg)

	// angleDeg := p.angleGun * cp.DegreeConst
	// var scale float64 = 1.0
	angle := p.angleGun
	// if angle < -math.Pi/2.0 {
	// 	// angle = -math.Pi - angle
	// 	scale = -1.0
	// } else if angle > math.Pi/2.0 {
	// 	// angle = math.Pi - angle
	// 	scale = -1.0
	// }
	fmt.Printf("angle: %v\n", angle)
	// p.drawOptionsGun.GeoM.Scale(scale, scale)
	p.drawOptionsGun.GeoM.Translate(0, -p.sizeGun.Y/2.0)
	p.drawOptionsGun.GeoM.Rotate(angle)
	p.drawOptionsGun.GeoM.Translate(p.posGun.X, p.posGun.Y)
}

func (p *player) draw(dst *ebiten.Image) {
	p.curAnim.Draw(dst, &p.drawOptions)

	// Draw prototype gun
	if p.stateGun == gunStateAttract {
		dst.DrawImage(imageGunAttract, &p.drawOptionsGun)
	} else if p.stateGun == gunStateRepel {
		dst.DrawImage(imageGunRepel, &p.drawOptionsGun)
	} else {
		dst.DrawImage(imageGunIdle, &p.drawOptionsGun)
	}

	// ebitenutil.DrawLine(dst, p.gunRay[0].X, p.gunRay[0].Y, p.gunRay[1].X, p.gunRay[1].Y, colorCrosshair)
}

func playerUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
