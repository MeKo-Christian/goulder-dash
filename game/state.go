package game

func canEnemyMoveTo(gsm *GameStateManager, x, y int) bool {
	if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
		return false
	}

	tile := gsm.GetTileAt(x, y)

	return tile == TileEmpty || tile == TilePlayer
}

func getDirectionOffset(dir Direction) (int, int) {
	switch dir {
	case FacingRight:
		return 1, 0
	case FacingDown:
		return 0, 1
	case FacingLeft:
		return -1, 0
	case FacingUp:
		return 0, -1
	}

	return 0, 0
}

func turnClockwise(dir Direction) Direction {
	switch dir {
	case FacingRight:
		return FacingDown
	case FacingDown:
		return FacingLeft
	case FacingLeft:
		return FacingUp
	case FacingUp:
		return FacingRight
	}

	return FacingRight
}

func updateEnemies(gsm *GameStateManager) {
	enemies := gsm.GetEnemies()
	for i := range enemies {
		enemy := &enemies[i]

		if enemy.MoveTimer > 0 {
			enemy.MoveTimer--
			continue
		}

		// Butterfly wall-following
		if enemy.Type == TileEnemy1 {
			dir := enemy.Direction
			for range 4 {
				dx, dy := getDirectionOffset(dir)
				if canEnemyMoveTo(gsm, enemy.X+dx, enemy.Y+dy) {
					// Move forward
					enemy.Direction = dir
					moveEnemy(gsm, enemy)

					break
				}

				dir = turnClockwise(dir)
			}
		}

		enemy.MoveTimer = 8 // Move every 8 frames
	}

	gsm.SetEnemies(enemies)
}

func moveEnemy(gsm *GameStateManager, enemy *Enemy) {
	// Move to new position
	dx, dy := getDirectionOffset(enemy.Direction)
	newX := enemy.X + dx
	newY := enemy.Y + dy

	// Check for player collision before moving
	playerX, playerY := gsm.GetPlayerPosition()
	if newX == playerX && newY == playerY {
		// Player dies - clear enemy's old position and start explosion
		gsm.SetTileAt(enemy.X, enemy.Y, TileEmpty)
		startExplosion(gsm, newX, newY)

		return
	}

	// Clear old position
	gsm.SetTileAt(enemy.X, enemy.Y, TileEmpty)

	// Update enemy position
	enemy.X = newX
	enemy.Y = newY

	// Place enemy in new position
	gsm.SetTileAt(enemy.X, enemy.Y, enemy.Type)
}

