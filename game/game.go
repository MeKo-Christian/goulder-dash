package game

import "github.com/gonutz/prototype/draw"

var gameStateManager *GameStateManager

func init() {
	gameStateManager = NewGameStateManager()
}

func Update(window draw.Window) {
	window.BlurImages(true)

	handlePlayerMovement(window, gameStateManager)

	gameStateManager.IncrementFrameCounter()

	if gameStateManager.GetFrameCounter()%10 == 0 {
		updatePhysics(gameStateManager)
	}

	// Update explosion animation every frame (faster than physics)
	if gameStateManager.GetFrameCounter()%5 == 0 {
		updateExplosions(gameStateManager)
	}

	// Update enemies every frame
	updateEnemies(gameStateManager)

	renderGame(window, gameStateManager)
}
