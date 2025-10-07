package game

func startExplosion(gsm *GameStateManager, centerX, centerY int, areaEffect bool, fallingObject Tile) {
	// Start the explosion at the center
	gsm.SetTileAt(centerX, centerY, TileExplosion0)

	if !areaEffect {
		return
	}

	// Check if a gem/diamond fell on the enemy
	isDiamondExplosion := fallingObject == TileGem

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

		// If diamond explosion, convert destructible tiles to diamonds
		if isDiamondExplosion {
			// Don't convert walls/borders, but convert everything else
			if tile != TileBrickWall && tile != TileStoneWall {
				// Remove enemies from the list if present
				if tile == TileEnemy1 || tile == TileEnemy2 || tile == TileEnemy3 {
					enemies := gsm.GetEnemies()
					for i, enemy := range enemies {
						if enemy.X == x && enemy.Y == y {
							enemies = append(enemies[:i], enemies[i+1:]...)
							break
						}
					}
					gsm.SetEnemies(enemies)
				}
				gsm.SetTileAt(x, y, TileGem)
			}
		} else {
			// Regular explosion - destroy surrounding tiles and show explosion animation
			if tile == TileDirt || tile == TileGem || tile == TileRock {
				gsm.SetTileAt(x, y, TileExplosion0)
			} else if tile == TileEnemy1 || tile == TileEnemy2 || tile == TileEnemy3 {
				// Remove enemy from list
				enemies := gsm.GetEnemies()
				for i, enemy := range enemies {
					if enemy.X == x && enemy.Y == y {
						enemies = append(enemies[:i], enemies[i+1:]...)
						break
					}
				}
				gsm.SetEnemies(enemies)
				gsm.SetTileAt(x, y, TileExplosion0)
			} else if tile == TileEmpty {
				// Even empty tiles show explosion animation
				gsm.SetTileAt(x, y, TileExplosion0)
			}
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
				gsm.SetTileAt(x, y, TileEmpty) // Explosion ends
				if gsm.IsPlayerKilled() {
					gsm.ResetCurrentLevel()
					return // Exit after reset
				}
			}
		}
	}
}