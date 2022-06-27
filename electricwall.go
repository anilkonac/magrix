package main

import (
	"math"

	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"github.com/yohamta/ganim8/v2"
)

type electricWall struct {
	shape       *cp.Shape
	drawOptions ganim8.DrawOptions
	anim        *ganim8.Animation
}

func newElectricWall(obj *tiled.Object, space *cp.Space) *electricWall {
	radius := math.Min(obj.Width, obj.Height) / 2.0
	x2 := obj.X + obj.Width - radius
	y2 := obj.Y + obj.Height - radius
	shape := space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: obj.X + radius, Y: obj.Y + radius}, cp.Vector{X: x2, Y: y2}, radius))
	shape.SetElasticity(wallElasticity)
	shape.SetFriction(wallFriction)

	anim := animElectricBlue
	if obj.Properties.GetBool("isOrange") {
		anim = animElectricOrange
	}

	return &electricWall{
		shape: shape,
		drawOptions: ganim8.DrawOptions{
			X:       obj.X + obj.Width/2.0,
			Y:       obj.Y + obj.Height/2.0,
			OriginX: 0.5,
			OriginY: 0.5,
			ScaleX:  1.0,
			ScaleY:  2.0,
		},
		anim: anim,
	}
}

func (e *electricWall) update() {
	e.anim.Update(animDeltaTime)
}

func (e *electricWall) draw() {
	e.anim.Draw(imageObjects, &e.drawOptions)
}
