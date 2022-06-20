// Copyright 2022 Anıl Konaç

package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

const (
	wallElasticity = 1
	wallFriction   = 1
	wallWidth      = 30
	wallRadius     = wallWidth / 2.0
)

var (
	imageWall = ebiten.NewImage(1, 1)
)

func init() {
	imageWall.Fill(colorWall)
}

type wall struct {
	shape       *cp.Shape
	drawOptions ebiten.DrawImageOptions
}

// v1 has to be either left or top position
func newWall(x1, y1, x2, y2, radius float64, space *cp.Space) *wall {
	shape := space.AddShape(cp.NewSegment(space.StaticBody, cp.Vector{X: x1, Y: y1}, cp.Vector{X: x2, Y: y2}, radius))
	shape.SetElasticity(wallElasticity)
	shape.SetFriction(wallFriction)

	wallWidth := 2 * radius
	width := math.Max(x2-x1+wallWidth, wallWidth)
	height := math.Max(y2-y1+wallWidth, wallWidth)

	wall := new(wall)
	wall.shape = shape
	wall.drawOptions.GeoM.Reset()
	wall.drawOptions.GeoM.Scale(width, height)
	wall.drawOptions.GeoM.Translate(x1-radius, y1-radius)

	return wall
}

func (w *wall) draw(screen *ebiten.Image) {
	screen.DrawImage(imageWall, &w.drawOptions)
}
