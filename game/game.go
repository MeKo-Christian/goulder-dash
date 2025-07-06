package game

import (
	"strconv"

	"github.com/gonutz/prototype/draw"
)

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

func handlePlayerMovement(w draw.Window) {
	dx, dy := 0, 0
	moved := false
	rockMoved := false

	if w.WasKeyPressed(draw.KeyLeft) {
		dx = -1
		playerDirection = FacingLeft
		moved = true
	} else if w.WasKeyPressed(draw.KeyRight) {
		dx = 1
		playerDirection = FacingRight
		moved = true
	} else if w.WasKeyPressed(draw.KeyUp) {
		dy = -1
		playerDirection = FacingUp
		moved = true
	} else if w.WasKeyPressed(draw.KeyDown) {
		dy = 1
		playerDirection = FacingDown
		moved = true
	}

	if !moved {
		return
	}

	newX := playerX + dx
	newY := playerY + dy
	target := tileMap[newY][newX]

	// Walls are always blocked
	if target == TileBrickWall || target == TileStoneWall {
		return
	}

	// Prevent entry into closed exit
	if target == TileClosedExit {
		return
	}

	// Transition if player enters open exit
	if target == TileOpenExit {
		loadNextLevel()
		return
	}

	// Handle enemy collision
	if target == TileEnemy1 || target == TileEnemy2 || target == TileEnemy3 {
		// Player dies when walking into enemy
		tileMap[newY][newX] = TileExplosion0
		return
	}

	// Handle pushing rock
	if target == TileRock {
		// Only allow horizontal pushing
		if dy != 0 {
			return
		}

		pushX := newX + dx
		pushY := newY

		if tileMap[pushY][pushX] == TileEmpty {
			// Move rock
			tileMap[pushY][pushX] = TileRock
			tileMap[newY][newX] = TileEmpty
			rockMoved = true
		} else {
			// Can't push if not empty behind
			return
		}
	}

	// Check for gem
	if target == TileGem {
		collectGem()
	}

	// Reset on each move
	playerHoldsFallingObject = false

	// Check for support case
	if tileMap[newY][newX] == TileDirt || tileMap[newY][newX] == TileGem || rockMoved {
		if newY > 0 {
			above := tileMap[newY-1][newX]
			if above == TileRock || above == TileGem {
				playerHoldsFallingObject = true
			}
		}
	}

	// Move player
	tileMap[playerY][playerX] = TileEmpty
	playerX = newX
	playerY = newY
	tileMap[playerY][playerX] = TilePlayer
}

func collectGem() {
	gemCounter++

	// Check if all gems collected
	if gemCounter >= currentLevel.GemTarget {
		// Open the exit
		for y := 0; y < GridHeight; y++ {
			for x := 0; x < GridWidth; x++ {
				if tileMap[y][x] == TileClosedExit {
					tileMap[y][x] = TileOpenExit
				}
			}
		}
	}
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

func turnCounterClockwise(dir Direction) Direction {
	switch dir {
	case FacingRight:
		return FacingUp
	case FacingUp:
		return FacingLeft
	case FacingLeft:
		return FacingDown
	case FacingDown:
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
	// Clear old position
	tileMap[enemy.Y][enemy.X] = TileEmpty

	// Move to new position
	dx, dy := getDirectionOffset(enemy.Direction)
	enemy.X += dx
	enemy.Y += dy

	// Check for player collision
	if enemy.X == playerX && enemy.Y == playerY {
		// Player dies
		tileMap[enemy.Y][enemy.X] = TileExplosion0
		return
	}

	// Place enemy in new position
	tileMap[enemy.Y][enemy.X] = enemy.Type
}

func updateExplosions() {
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
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

func freeSurroundingBlocks(centerX, centerY int) {
	// Define the 8 surrounding positions
	offsets := [][]int{
		{-1, -1}, {0, -1}, {1, -1}, // left up, up, right up
		{-1, 0}, {1, 0},            // left, right
		{-1, 1}, {0, 1}, {1, 1},    // left down, down, right down
	}

	for _, offset := range offsets {
		x := centerX + offset[0]
		y := centerY + offset[1]

		// Check bounds
		if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
			continue
		}

		// Free the block (except for hard walls)
		if tileMap[y][x] != TileStoneWall {
			tileMap[y][x] = TileEmpty
		}
	}
}

func createSurroundingDiamonds(centerX, centerY int) {
	// Define the 8 surrounding positions
	offsets := [][]int{
		{-1, -1}, {0, -1}, {1, -1}, // left up, up, right up
		{-1, 0}, {1, 0},            // left, right
		{-1, 1}, {0, 1}, {1, 1},    // left down, down, right down
	}

	for _, offset := range offsets {
		x := centerX + offset[0]
		y := centerY + offset[1]

		// Check bounds
		if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
			continue
		}

		// Create diamonds (except for hard walls)
		if tileMap[y][x] != TileStoneWall {
			tileMap[y][x] = TileGem
		}
	}
}

func updatePhysics() {
	// Process bottom-up
	for y := GridHeight - 2; y >= 1; y-- {
		for x := 1; x < GridWidth-1; x++ {
			tile := tileMap[y][x]

			// Skip empty tiles and explosions (handled separately)
			if tile != TileRock && tile != TileGem {
				continue
			}

			// FALL STRAIGHT
			if tileMap[y+1][x] == TileEmpty {
				tileMap[y+1][x] = tile
				tileMap[y][x] = TileEmpty
				continue
			}

			// ROLL RIGHT
			if (tileMap[y+1][x] == TileRock || tileMap[y+1][x] == TileGem) &&
				tileMap[y][x+1] == TileEmpty &&
				tileMap[y+1][x+1] == TileEmpty {
				tileMap[y+1][x+1] = tile
				tileMap[y][x] = TileEmpty
				continue
			}

			// ROLL LEFT
			if (tileMap[y+1][x] == TileRock || tileMap[y+1][x] == TileGem) &&
				tileMap[y][x-1] == TileEmpty &&
				tileMap[y+1][x-1] == TileEmpty {
				tileMap[y+1][x-1] = tile
				tileMap[y][x] = TileEmpty
				continue
			}

			// FALL ON PLAYER
			if tileMap[y+1][x] == TilePlayer {
				if !playerHoldsFallingObject {
					// Player dies
					tileMap[y+1][x] = TileExplosion0
					tileMap[y][x] = TileEmpty
				}
			}

			// FALL ON ENEMY
			if tileMap[y+1][x] == TileEnemy1 || tileMap[y+1][x] == TileEnemy2 || tileMap[y+1][x] == TileEnemy3 {
				// Remove enemy from list
				for i, enemy := range enemies {
					if enemy.X == x && enemy.Y == y+1 {
						enemies = append(enemies[:i], enemies[i+1:]...)
						break
					}
				}
				
				if tile == TileRock {
					// Rock falling on enemy: free surrounding blocks (except hard walls)
					freeSurroundingBlocks(x, y+1)
					tileMap[y+1][x] = TileEmpty
					tileMap[y][x] = TileEmpty
				} else if tile == TileGem {
					// Gem falling on enemy: kill enemy and create diamonds in surrounding blocks
					createSurroundingDiamonds(x, y+1)
					tileMap[y+1][x] = TileGem
					tileMap[y][x] = TileEmpty
				}
			}
		}
	}
}

func Update(w draw.Window) {
	w.BlurImages(true)

	handlePlayerMovement(w)

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

	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			tile := tileMap[y][x]
			spriteIndex := TileSpriteIndex[tile]

			if tile == TilePlayer {
				switch playerDirection {
				case FacingRight:
					spriteIndex = 0
				case FacingDown:
					spriteIndex = 1
				case FacingLeft:
					spriteIndex = 2
				case FacingUp:
					spriteIndex = 3
				}
			}

			sx := (spriteIndex % TileCols) * TileSize
			sy := (spriteIndex / TileCols) * TileSize

			err := w.DrawImageFilePart(
				SpriteSheet,
				sx, sy, TileSize, TileSize,
				x*TileDrawSize, y*TileDrawSize, TileDrawSize, TileDrawSize,
				0,
			)
			if err != nil {
				w.DrawText("Failed to load sprite!", 10, 10, draw.Red)
			}
		}
	}

	text := currentLevel.Name + " - Gems: " + strconv.Itoa(gemCounter) + " / " + strconv.Itoa(currentLevel.GemTarget)
	w.DrawText(text, 8, 8, draw.White)
}
