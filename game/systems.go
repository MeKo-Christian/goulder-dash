package game

var playerSystem *PlayerSystem

func updateSystems(gsm *GameStateManager) {
	// Initialize player system if it doesn't exist
	if playerSystem == nil {
		playerSystem = NewPlayerSystem(gsm)
	}

	// Update player key repeat timer
	playerSystem.UpdateKeyRepeat()

	// Update physics every 10 frames
	if gsm.GetFrameCounter()%10 == 0 {
		updatePhysics(gsm)
	}

	// Update explosion animation every 5 frames
	if gsm.GetFrameCounter()%5 == 0 {
		updateExplosions(gsm)
	}

	// Update enemies every frame
	updateEnemies(gsm)
}