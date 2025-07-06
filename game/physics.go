package game

func updatePhysics() {
	// Process bottom-up
	for y := GridHeight - 2; y >= 1; y-- {
		for x := 1; x < GridWidth-1; x++ {
			tile := tileMap[y][x]

			// Skip empty tiles and explosions (handled separately)
			if tile != TileRock && tile != TileGem {
				continue
			}

			handleFallingObject(x, y, tile)
		}
	}
}

func handleFallingObject(x, y int, tile Tile) {
	// Fall straight down
	if tileMap[y+1][x] == TileEmpty {
		tileMap[y+1][x] = tile
		tileMap[y][x] = TileEmpty

		return
	}

	// Roll right
	if canRollRight(x, y) {
		tileMap[y+1][x+1] = tile
		tileMap[y][x] = TileEmpty

		return
	}

	// Roll left
	if canRollLeft(x, y) {
		tileMap[y+1][x-1] = tile
		tileMap[y][x] = TileEmpty

		return
	}

	// Handle collisions
	handleObjectCollision(x, y, tile)
}

func canRollRight(x, y int) bool {
	return (tileMap[y+1][x] == TileRock || tileMap[y+1][x] == TileGem) &&
		tileMap[y][x+1] == TileEmpty &&
		tileMap[y+1][x+1] == TileEmpty
}

func canRollLeft(x, y int) bool {
	return (tileMap[y+1][x] == TileRock || tileMap[y+1][x] == TileGem) &&
		tileMap[y][x-1] == TileEmpty &&
		tileMap[y+1][x-1] == TileEmpty
}

func handleObjectCollision(x, y int, tile Tile) {
	target := tileMap[y+1][x]

	// FALL ON PLAYER
	if target == TilePlayer {
		if !playerHoldsFallingObject {
			// Player dies
			tileMap[y+1][x] = TileExplosion0
			tileMap[y][x] = TileEmpty
		}

		return
	}

	// FALL ON ENEMY
	if target == TileEnemy1 || target == TileEnemy2 || target == TileEnemy3 {
		handleObjectFallOnEnemy(x, y, tile)
	}
}

func handleObjectFallOnEnemy(x, y int, tile Tile) {
	// Remove enemy from list
	for i, enemy := range enemies {
		if enemy.X == x && enemy.Y == y+1 {
			enemies = append(enemies[:i], enemies[i+1:]...)
			break
		}
	}

	switch tile {
	case TileRock:
		// Rock falling on enemy: free surrounding blocks (except hard walls)
		freeSurroundingBlocks(x, y+1)
		tileMap[y+1][x] = TileEmpty
		tileMap[y][x] = TileEmpty
	case TileGem:
		// Gem falling on enemy: kill enemy and create diamonds in surrounding blocks
		createSurroundingDiamonds(x, y+1)
		tileMap[y+1][x] = TileGem
		tileMap[y][x] = TileEmpty
	}
}

func freeSurroundingBlocks(centerX, centerY int) {
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
		if tileMap[y][x] != TileStoneWall {
			tileMap[y][x] = TileEmpty
		}
	}
}

func createSurroundingDiamonds(centerX, centerY int) {
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
		if tileMap[y][x] != TileStoneWall {
			tileMap[y][x] = TileGem
		}
	}
}