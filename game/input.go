package game

import "github.com/gonutz/prototype/draw"

var playerSystem *PlayerSystem

const (
	KeyCodeLeft  = 1
	KeyCodeRight = 2
	KeyCodeUp    = 3
	KeyCodeDown  = 4
)

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

	// Check for initial key presses (immediate response)
	if window.WasKeyPressed(draw.KeyLeft) {
		playerSystem.SetDirection(FacingLeft)
		if playerSystem.Move(-1, 0) {
			playerSystem.StartKeyRepeat(KeyCodeLeft)
		}
		return
	}
	if window.WasKeyPressed(draw.KeyRight) {
		playerSystem.SetDirection(FacingRight)
		if playerSystem.Move(1, 0) {
			playerSystem.StartKeyRepeat(KeyCodeRight)
		}
		return
	}
	if window.WasKeyPressed(draw.KeyUp) {
		playerSystem.SetDirection(FacingUp)
		if playerSystem.Move(0, -1) {
			playerSystem.StartKeyRepeat(KeyCodeUp)
		}
		return
	}
	if window.WasKeyPressed(draw.KeyDown) {
		playerSystem.SetDirection(FacingDown)
		if playerSystem.Move(0, 1) {
			playerSystem.StartKeyRepeat(KeyCodeDown)
		}
		return
	}

	// Check for key repeat (after initial delay)
	if window.IsKeyDown(draw.KeyLeft) && playerSystem.CanRepeatKey(KeyCodeLeft) {
		playerSystem.SetDirection(FacingLeft)
		if playerSystem.Move(-1, 0) {
			playerSystem.SetRepeatInterval()
		}
	} else if window.IsKeyDown(draw.KeyRight) && playerSystem.CanRepeatKey(KeyCodeRight) {
		playerSystem.SetDirection(FacingRight)
		if playerSystem.Move(1, 0) {
			playerSystem.SetRepeatInterval()
		}
	} else if window.IsKeyDown(draw.KeyUp) && playerSystem.CanRepeatKey(KeyCodeUp) {
		playerSystem.SetDirection(FacingUp)
		if playerSystem.Move(0, -1) {
			playerSystem.SetRepeatInterval()
		}
	} else if window.IsKeyDown(draw.KeyDown) && playerSystem.CanRepeatKey(KeyCodeDown) {
		playerSystem.SetDirection(FacingDown)
		if playerSystem.Move(0, 1) {
			playerSystem.SetRepeatInterval()
		}
	}
}
