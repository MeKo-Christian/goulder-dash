package game

import "github.com/gonutz/prototype/draw"

func handlePlayerMovement(window draw.Window) {
	// Prevent movement if player is exploding
	currentTile := tileMap[playerY][playerX]
	if currentTile >= TileExplosion0 && currentTile <= TileExplosion5 {
		return
	}

	dx, dy := 0, 0
	moved := false

	switch {
	case window.WasKeyPressed(draw.KeyLeft):
		dx = -1
		playerDirection = FacingLeft
		moved = true
	case window.WasKeyPressed(draw.KeyRight):
		dx = 1
		playerDirection = FacingRight
		moved = true
	case window.WasKeyPressed(draw.KeyUp):
		dy = -1
		playerDirection = FacingUp
		moved = true
	case window.WasKeyPressed(draw.KeyDown):
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

	if !canPlayerMoveTo(target, newX, newY, dx, dy) {
		return
	}

	movePlayerTo(newX, newY, target, dx, dy)
}

func canPlayerMoveTo(target Tile, newX, newY, dx, dy int) bool {
	// Walls are always blocked
	if target == TileBrickWall || target == TileStoneWall {
		return false
	}

	// Prevent entry into closed exit
	if target == TileClosedExit {
		return false
	}

	// Handle pushing rock
	if target == TileRock {
		return canPushRock(newX, newY, dx, dy)
	}

	return true
}

func canPushRock(newX, newY, dx, dy int) bool {
	// Only allow horizontal pushing
	if dy != 0 {
		return false
	}

	pushX := newX + dx
	pushY := newY

	return tileMap[pushY][pushX] == TileEmpty
}

func movePlayerTo(newX, newY int, target Tile, dx, _ int) {
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

	rockMoved := handleRockPushing(newX, newY, dx, target)
	handleGemCollection(target)
	updatePlayerSupport(newX, newY, target, rockMoved)

	// Move player
	tileMap[playerY][playerX] = TileEmpty
	playerX = newX
	playerY = newY
	tileMap[playerY][playerX] = TilePlayer
}

func handleRockPushing(newX, newY, dx int, target Tile) bool {
	if target != TileRock {
		return false
	}

	pushX := newX + dx
	pushY := newY

	// Move rock
	tileMap[pushY][pushX] = TileRock
	tileMap[newY][newX] = TileEmpty

	return true
}

func handleGemCollection(target Tile) {
	if target == TileGem {
		collectGem()
	}
}

func updatePlayerSupport(newX, newY int, target Tile, rockMoved bool) {
	// Reset on each move
	playerHoldsFallingObject = false

	// Check for support case
	if target == TileDirt || target == TileGem || rockMoved {
		if newY > 0 {
			above := tileMap[newY-1][newX]
			if above == TileRock || above == TileGem {
				playerHoldsFallingObject = true
			}
		}
	}
}

func collectGem() {
	gemCounter++

	// Check if all gems collected
	if gemCounter >= currentLevel.GemTarget {
		// Open the exit
		for y := range GridHeight {
			for x := range GridWidth {
				if tileMap[y][x] == TileClosedExit {
					tileMap[y][x] = TileOpenExit
				}
			}
		}
	}
}
