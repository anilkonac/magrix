package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"github.com/yohamta/ganim8/v2"
)

type terminal struct {
	shape       *cp.Shape
	drawOptions ganim8.DrawOptions
	spr         *ganim8.Sprite
	triggered   bool
}

func newTerminal(obj *tiled.Object, space *cp.Space) *terminal {
	var shape *cp.Shape
	if obj.Properties.GetBool("blocking") {
		radius := math.Min(obj.Width, obj.Height) / 2.0
		x2 := obj.X + obj.Width - radius
		y2 := obj.Y + obj.Height - radius
		shape = space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: obj.X + radius, Y: obj.Y + radius}, cp.Vector{X: x2, Y: y2}, radius))
		// shape.SetElasticity(wallElasticity)
		// shape.SetFriction(wallFriction)
	}

	var drawOpts ebiten.DrawImageOptions
	drawOpts.GeoM.Reset()
	drawOpts.GeoM.Translate(obj.X+obj.Width/2.0, obj.Y+obj.Height/2.0)

	spr := spriteTerminalBlue
	if obj.Properties.GetBool("isOrange") {
		spr = spriteTerminalOrange
	}
	return &terminal{
		shape: shape,
		drawOptions: ganim8.DrawOptions{
			X:       obj.X + obj.Width/2.0,
			Y:       obj.Y + obj.Height/2.0,
			OriginX: 0.5,
			OriginY: 0.5,
			ScaleX:  1.0,
			ScaleY:  1.0,
		},
		spr: spr,
	}
}

func (t *terminal) trigger() {
	t.triggered = true
}

func (t *terminal) draw() {
	var index int
	if t.triggered {
		index = 1
	}
	t.spr.Draw(imageObjects, index, &t.drawOptions)
}
