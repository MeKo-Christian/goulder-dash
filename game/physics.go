package game

func UpdatePhysics(gsm *GameStateManager) {
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
			gsm.SetPlayerKilled(true)
			startExplosion(gsm, x, y+1, false, tile)
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

	// An enemy was crushed, trigger an area-effect explosion.
	// Pass the falling object type so we know if it's a diamond explosion
	startExplosion(gsm, x, y+1, true, tile)
	gsm.SetTileAt(x, y, TileEmpty)
}
