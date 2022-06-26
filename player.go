// Copyright 2022 Anıl Konaç

package main

import (
	"bytes"
	"image/png"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/yohamta/ganim8/v2"
)

type playerState uint8

const (
	stateIdle playerState = iota
	stateWalking
	stateFiring
	stateJumping
	stateTotal
)

const (
	halfPi        = math.Pi / 2.0
	turnTolerance = 5 * cp.RadianConst
)

const (
	playerWidthTile  = 10.0 / 16.0
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
	imagePlayer     *ebiten.Image
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

	imagePlayer = ebiten.NewImage(16, 32)
}

type player struct {
	pos             cp.Vector
	posGun          cp.Vector
	size            cp.Vector
	sizeGun         cp.Vector
	angleGun        float64
	shape           *cp.Shape
	body            *cp.Body
	drawOptions     ebiten.DrawImageOptions
	drawOptionsAnim ganim8.DrawOptions
	drawOptionsGun  ebiten.DrawImageOptions
	curAnim         *ganim8.Animation
	onGround        bool
	gunRay          [2]cp.Vector
	gunForce        cp.Vector
	state           playerState
	stateGun        gunState
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
		drawOptionsAnim: ganim8.DrawOptions{
			OriginX: 0.0,
			OriginY: 0.1,
			ScaleX:  1.0,
			ScaleY:  1.0,
		},
		state:   stateIdle,
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
	if p.angleGun < -halfPi || p.angleGun > halfPi {
		p.posGun = p.pos.Add(cp.Vector{X: -posGunRelative.X, Y: posGunRelative.Y})
	} else {
		p.posGun = p.pos.Add(posGunRelative)
	}

	// Update gun angle
	distX := cursorX - p.posGun.X
	distY := cursorY - p.posGun.Y
	p.angleGun = math.Atan2(distY, distX)

	p.checkOnGround()

	// Raycast
	const rayLength = gunRange
	p.gunRay[0] = p.posGun
	p.gunRay[1] = p.gunRay[0].Add(cp.Vector{
		X: rayLength * math.Cos(p.angleGun), Y: rayLength * math.Sin(p.angleGun),
	})

	p.handleInputs(input, rayHitInfo)

	switch p.state {
	case stateWalking:
		p.curAnim = animPlayerWalk
	default:
		p.curAnim = animPlayerIdle
	}

	// v := p.body.Velocity()
	// fmt.Printf("Friction: %.2f\tVel X: %.2f\tVel Y: %.2f\n", p.shape.Friction(), v.X, v.Y)
	p.curAnim.Update(animDeltaTime)
	p.updateDrawOptions()

	// fmt.Printf("p.angleGun: %v\n", p.angleGun)
}

func (p *player) checkOnGround() {
	const groundNormalYThreshold = 0.8
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
	p.state = stateIdle
	// Handle inputs
	var surfaceV cp.Vector
	if input.right {
		surfaceV.X = -playerVelocity
		p.state = stateWalking
	} else if input.left {
		surfaceV.X = playerVelocity
		p.state = stateWalking
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
		newVelX := cp.Clamp(v.X-surfaceV.X*deltaTimeSec, -playerVelocity, playerVelocity)
		p.body.SetVelocity(newVelX, v.Y)
		p.state = stateJumping
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
		p.state = stateFiring
	} else {
		p.stateGun = gunStateIdle
	}
	// v := p.body.Velocity()
	// fmt.Printf("Velocity X: %.2f\tY: %.2f\t\tForce X: %.2f\tY:%.2f\n", v.X, v.Y, p.gunForce.X, p.gunForce.Y)
	// fmt.Printf("Force X:%.2f\tY:%.2f\n", p.gunForce.X, p.gunForce.Y)
}

func (p *player) updateDrawOptions() {
	p.drawOptions = *cam.GetTranslation(p.pos.X-tileLength/2.0, p.pos.Y-tileLength)
	// Player
	if p.angleGun < -halfPi || p.angleGun > halfPi {
		p.drawOptionsAnim.ScaleX = -1.0
		p.drawOptionsAnim.OriginX = 1.0
	} else {
		p.drawOptionsAnim.ScaleX = 1.0
		p.drawOptionsAnim.OriginX = 0.0

	}

	// Gun
	p.drawOptionsGun.GeoM.Reset()
	p.drawOptionsGun.GeoM.Translate(0, -p.sizeGun.Y/2.0)
	p.drawOptionsGun.GeoM.Rotate(p.angleGun)
	p.drawOptionsGun.GeoM.Concat(cam.GetTranslation(p.posGun.X, p.posGun.Y).GeoM)
}

func (p *player) draw() {
	// Draw player
	imagePlayer.Clear()
	p.curAnim.Draw(imagePlayer, &p.drawOptionsAnim)
	cam.Surface.DrawImage(imagePlayer, &p.drawOptions)

	// Draw gun
	if p.stateGun == gunStateAttract {
		cam.Surface.DrawImage(imageGunAttract, &p.drawOptionsGun)
	} else if p.stateGun == gunStateRepel {
		cam.Surface.DrawImage(imageGunRepel, &p.drawOptionsGun)
	} else {
		cam.Surface.DrawImage(imageGunIdle, &p.drawOptionsGun)
	}

	// ebitenutil.DrawLine(dst, p.gunRay[0].X, p.gunRay[0].Y, p.gunRay[1].X, p.gunRay[1].Y, colorCrosshair)
}

func playerUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}
