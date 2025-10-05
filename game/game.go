package game

import "github.com/gonutz/prototype/draw"

var gameStateManager *GameStateManager

func init() {
	gameStateManager = NewGameStateManager()
}

func Update(window draw.Window) {
	window.BlurImages(true)

	updateSystems(gameStateManager)
	handlePlayerMovement(window, gameStateManager)

	gameStateManager.IncrementFrameCounter()

	renderGame(window, gameStateManager)
}
