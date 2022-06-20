package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

type input struct {
	cursorPos cp.Vector
	up/*, down*/ bool
	left, right bool

	attract bool
	repel   bool

	escape bool
}

func (i *input) update() {
	// Update cursor pos
	x, y := ebiten.CursorPosition()
	i.cursorPos = cp.Vector{X: float64(x), Y: float64(y)}

	// Update movement key states
	i.right = ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	i.left = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	i.up = ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeySpace)
	i.attract = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	i.repel = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	i.escape = ebiten.IsKeyPressed(ebiten.KeyEscape)
	// i.down = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyControlLeft)
}
