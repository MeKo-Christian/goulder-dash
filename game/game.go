package game

import "github.com/gonutz/prototype/draw"

func Update(window draw.Window) {
	window.BlurImages(true)

	handlePlayerMovement(window)

	frameCounter++
	if frameCounter%10 == 0 {
		updatePhysics()
	}

	// Update explosion animation every frame (faster than physics)
	if frameCounter%5 == 0 {
		updateExplosions()
	}

	// Update enemies every frame
	updateEnemies()

	renderGame(window)
}
