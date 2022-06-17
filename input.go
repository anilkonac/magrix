package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

type input struct {
	cursorPos   cp.Vector
	up, down    bool
	left, right bool
}

func (i *input) update() {
	// Update cursor pos
	x, y := ebiten.CursorPosition()
	i.cursorPos = cp.Vector{X: float64(x), Y: float64(y)}

	// Update movement key states
	i.right = ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	i.left = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
}
