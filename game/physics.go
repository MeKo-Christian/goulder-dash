package game

func updatePhysics(gsm *GameStateManager) {
	// Process bottom-up
	for y := GridHeight - 2; y >= 1; y-- {
		for x := 1; x < GridWidth-1; x++ {
			tile := gsm.GetTileAt(x, y)

			// Skip empty tiles and explosions (handled separately)
			if tile != TileRock && tile != TileGem {
				continue
			}

			handleFallingObject(gsm, x, y, tile)
		}
	}
}

func handleFallingObject(gsm *GameStateManager, x, y int, tile Tile) {
	// Fall straight down
	if gsm.GetTileAt(x, y+1) == TileEmpty {
		gsm.SetTileAt(x, y+1, tile)
		gsm.SetTileAt(x, y, TileEmpty)

		return
	}

	// Roll right
	if canRollRight(gsm, x, y) {
		gsm.SetTileAt(x+1, y+1, tile)
		gsm.SetTileAt(x, y, TileEmpty)

		return
	}

	// Roll left
	if canRollLeft(gsm, x, y) {
		gsm.SetTileAt(x-1, y+1, tile)
		gsm.SetTileAt(x, y, TileEmpty)

		return
	}

	// Handle collisions
	handleObjectCollision(gsm, x, y, tile)
}

func canRollRight(gsm *GameStateManager, x, y int) bool {
	return (gsm.GetTileAt(x, y+1) == TileRock || gsm.GetTileAt(x, y+1) == TileGem) &&
		gsm.GetTileAt(x+1, y) == TileEmpty &&
		gsm.GetTileAt(x+1, y+1) == TileEmpty
}

func canRollLeft(gsm *GameStateManager, x, y int) bool {
	return (gsm.GetTileAt(x, y+1) == TileRock || gsm.GetTileAt(x, y+1) == TileGem) &&
		gsm.GetTileAt(x-1, y) == TileEmpty &&
		gsm.GetTileAt(x-1, y+1) == TileEmpty
}

func handleObjectCollision(gsm *GameStateManager, x, y int, tile Tile) {
	target := gsm.GetTileAt(x, y+1)

	// FALL ON PLAYER
	if target == TilePlayer {
		if !gsm.GetPlayerHoldsFallingObject() {
			// Player dies
			startExplosion(gsm, x, y+1)
			gsm.SetTileAt(x, y, TileEmpty)
		}

		return
	}

	// FALL ON ENEMY
	if target == TileEnemy1 || target == TileEnemy2 || target == TileEnemy3 {
		handleObjectFallOnEnemy(gsm, x, y, tile)
	}
}

func handleObjectFallOnEnemy(gsm *GameStateManager, x, y int, tile Tile) {
	// Remove enemy from list
	enemies := gsm.GetEnemies()
	for i, enemy := range enemies {
		if enemy.X == x && enemy.Y == y+1 {
			enemies = append(enemies[:i], enemies[i+1:]...)
			break
		}
	}

	gsm.SetEnemies(enemies)

	switch tile {
	case TileRock:
		// Rock falling on enemy: free surrounding blocks (except hard walls)
		freeSurroundingBlocks(gsm, x, y+1)
		gsm.SetTileAt(x, y+1, TileEmpty)
		gsm.SetTileAt(x, y, TileEmpty)
	case TileGem:
		// Gem falling on enemy: kill enemy and create diamonds in surrounding blocks
		createSurroundingDiamonds(gsm, x, y+1)
		gsm.SetTileAt(x, y+1, TileGem)
		gsm.SetTileAt(x, y, TileEmpty)
	default:
		// Handle any other tile types that shouldn't fall on enemies
		gsm.SetTileAt(x, y+1, TileEmpty)
		gsm.SetTileAt(x, y, TileEmpty)
	}
}

func freeSurroundingBlocks(gsm *GameStateManager, centerX, centerY int) {
	// Define the 8 surrounding positions
	offsets := [][]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}

	for _, offset := range offsets {
		x := centerX + offset[0]
		y := centerY + offset[1]

		// Check bounds
		if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
			continue
		}

		// Free the block (except for hard walls)
		if gsm.GetTileAt(x, y) != TileStoneWall {
			gsm.SetTileAt(x, y, TileEmpty)
		}
	}
}

func createSurroundingDiamonds(gsm *GameStateManager, centerX, centerY int) {
	// Define the 8 surrounding positions
	offsets := [][]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}

	for _, offset := range offsets {
		x := centerX + offset[0]
		y := centerY + offset[1]

		// Check bounds
		if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
			continue
		}

		// Create diamonds (except for hard walls)
		if gsm.GetTileAt(x, y) != TileStoneWall {
			gsm.SetTileAt(x, y, TileGem)
		}
	}
}
