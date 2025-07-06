package game

var (
	playerX, playerY         = 1, 1
	playerDirection          = FacingDown
	playerHoldsFallingObject = false
	frameCounter             = 0
	gemCounter               = 0
	tileMap                  [GridHeight][GridWidth]Tile
	currentLevel             LevelData
	currentLevelIndex        = 0
	enemies                  []Enemy
	levels                   = []LevelData{
		createGeneratedLevel("Level 1", 42, 40, 20, 1),
		createGeneratedLevel("Level 2", 32, 45, 25, 0),
		createGeneratedLevel("Level 3", 12, 50, 30, 3),
		createGeneratedLevel("Level 4", 22, 55, 35, 4),
		createGeneratedLevel("Level 5", 52, 60, 40, 5),
		createGeneratedLevel("Level 6", 62, 65, 45, 6),
	}
)

func init() {
	currentLevel = levels[0]
	tileMap = currentLevel.Grid

	// Initialize enemies for the first level
	enemies = nil

	for y := range GridHeight {
		for x := range GridWidth {
			if tileMap[y][x] == TileEnemy1 {
				enemies = append(enemies, Enemy{
					X: x, Y: y,
					Type:      TileEnemy1,
					Direction: FacingRight,
					MoveTimer: 8,
				})
			}
		}
	}
}

func loadNextLevel() {
	currentLevelIndex++
	if currentLevelIndex >= len(levels) {
		currentLevelIndex = 0
	}

	resetLevel(currentLevelIndex)
}

func canEnemyMoveTo(x, y int) bool {
	if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
		return false
	}

	tile := tileMap[y][x]

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

func updateEnemies() {
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
				if canEnemyMoveTo(enemy.X+dx, enemy.Y+dy) {
					// Move forward
					enemy.Direction = dir
					moveEnemy(enemy)

					break
				}

				dir = turnClockwise(dir)
			}
		}

		enemy.MoveTimer = 8 // Move every 8 frames
	}
}

func moveEnemy(enemy *Enemy) {
	// Move to new position
	dx, dy := getDirectionOffset(enemy.Direction)
	newX := enemy.X + dx
	newY := enemy.Y + dy

	// Check for player collision before moving
	if newX == playerX && newY == playerY {
		// Player dies - clear enemy's old position and start explosion
		tileMap[enemy.Y][enemy.X] = TileEmpty
		tileMap[newY][newX] = TileExplosion0
		return
	}

	// Clear old position
	tileMap[enemy.Y][enemy.X] = TileEmpty

	// Update enemy position
	enemy.X = newX
	enemy.Y = newY

	// Place enemy in new position
	tileMap[enemy.Y][enemy.X] = enemy.Type
}

func updateExplosions() {
	for y := range GridHeight {
		for x := range GridWidth {
			tile := tileMap[y][x]

			if tile >= TileExplosion0 && tile < TileExplosion5 {
				tileMap[y][x]++ // next frame
			} else if tile == TileExplosion5 {
				// reset level
				resetLevel(currentLevelIndex)
				return // Exit after reset
			}
		}
	}
}
