package game

import "github.com/gonutz/prototype/draw"

const (
	KeyCodeLeft  = 1
	KeyCodeRight = 2
	KeyCodeUp    = 3
	KeyCodeDown  = 4
)

func handlePlayerMovement(window draw.Window, gsm *GameStateManager) {
	// Initialize player system if not already done
	if playerSystem == nil {
		playerSystem = NewPlayerSystem(gsm)
	}

	// Prevent movement if player is exploding
	if playerSystem.IsExploding() {
		return
	}

	type keyAction struct {
		key  draw.Key
		dx   int
		dy   int
		dir  Direction
		code int
	}

	actions := []keyAction{
		{draw.KeyLeft, -1, 0, FacingLeft, KeyCodeLeft},
		{draw.KeyRight, 1, 0, FacingRight, KeyCodeRight},
		{draw.KeyUp, 0, -1, FacingUp, KeyCodeUp},
		{draw.KeyDown, 0, 1, FacingDown, KeyCodeDown},
	}

	// Check for initial key presses (immediate response)
	for _, a := range actions {
		if window.WasKeyPressed(a.key) {
			playerSystem.SetDirection(a.dir)

			if playerSystem.Move(a.dx, a.dy) {
				playerSystem.StartKeyRepeat(a.code)
			}

			return
		}
	}

	// Check for key repeat (after initial delay)
	for _, a := range actions {
		if window.IsKeyDown(a.key) && playerSystem.CanRepeatKey(a.code) {
			playerSystem.SetDirection(a.dir)

			if playerSystem.Move(a.dx, a.dy) {
				playerSystem.SetRepeatInterval()
			}
		}
	}
}
