package game

func startExplosion(gsm *GameStateManager, centerX, centerY int) {
	// Start the explosion at the center
	gsm.SetTileAt(centerX, centerY, TileExplosion0)

	// Define the 8 surrounding positions
	offsets := [][]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, offset := range offsets {
		x := centerX + offset[0]
		y := centerY + offset[1]

		// Check bounds
		if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
			continue
		}

		tile := gsm.GetTileAt(x, y)
		// Destroy surrounding dirt and gems.
		if tile == TileDirt || tile == TileGem {
			gsm.SetTileAt(x, y, TileEmpty)
		}
	}
}

func updateExplosions(gsm *GameStateManager) {
	for y := range GridHeight {
		for x := range GridWidth {
			tile := gsm.GetTileAt(x, y)

			if tile >= TileExplosion0 && tile < TileExplosion5 {
				gsm.SetTileAt(x, y, tile+1) // next frame
			} else if tile == TileExplosion5 {
				// reset level
				gsm.ResetCurrentLevel()
				return // Exit after reset
			}
		}
	}
}