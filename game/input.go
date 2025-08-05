package game

import "github.com/gonutz/prototype/draw"

var playerSystem *PlayerSystem

func init() {
	// playerSystem will be initialized when gameStateManager is created
}

func handlePlayerMovement(window draw.Window, gsm *GameStateManager) {
	// Initialize player system if not already done
	if playerSystem == nil {
		playerSystem = NewPlayerSystem(gsm)
	}

	// Prevent movement if player is exploding
	if playerSystem.IsExploding() {
		return
	}

	dx, dy := 0, 0
	moved := false

	switch {
	case window.WasKeyPressed(draw.KeyLeft):
		dx = -1

		playerSystem.SetDirection(FacingLeft)

		moved = true
	case window.WasKeyPressed(draw.KeyRight):
		dx = 1

		playerSystem.SetDirection(FacingRight)

		moved = true
	case window.WasKeyPressed(draw.KeyUp):
		dy = -1

		playerSystem.SetDirection(FacingUp)

		moved = true
	case window.WasKeyPressed(draw.KeyDown):
		dy = 1

		playerSystem.SetDirection(FacingDown)

		moved = true
	}

	if !moved {
		return
	}

	playerSystem.Move(dx, dy)
}
