package game

import (
	"math/rand"
)

type LevelData struct {
	Name      string
	Grid      [GridHeight][GridWidth]Tile
	GemTarget int
	RockCount int
	GemCount  int
	Seed      int64
}

func createGeneratedLevel(name string, seed int64, rockCount, gemCount, enemyCount int) LevelData {
	var grid [GridHeight][GridWidth]Tile

	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			if y == 0 || y == GridHeight-1 || x == 0 || x == GridWidth-1 {
				grid[y][x] = TileStoneWall
			} else {
				grid[y][x] = TileDirt
			}
		}
	}

	grid[1][1] = TilePlayer
	grid[GridHeight-2][GridWidth-2] = TileClosedExit

	r := rand.New(rand.NewSource(seed))
	placeRandomTiles(&grid, r, TileRock, rockCount)
	placeRandomTiles(&grid, r, TileGem, gemCount)
	if enemyCount > 0 {
		placeEnemiesWithSpace(&grid, r, TileEnemy1, enemyCount)
	}

	return LevelData{
		Name:      name,
		Grid:      grid,
		GemTarget: gemCount,
		RockCount: rockCount,
		GemCount:  0,
		Seed:      seed,
	}
}

func placeRandomTiles(level *[GridHeight][GridWidth]Tile, r *rand.Rand, tile Tile, count int) {
	placed := 0
	for placed < count {
		x := r.Intn(GridWidth-2) + 1
		y := r.Intn(GridHeight-2) + 1

		if level[y][x] == TileDirt {
			level[y][x] = tile
			placed++
		}
	}
}

func placeEnemiesWithSpace(level *[GridHeight][GridWidth]Tile, r *rand.Rand, tile Tile, count int) {
	placed := 0
	attempts := 0
	maxAttempts := count * 50 // Prevent infinite loops

	for placed < count && attempts < maxAttempts {
		attempts++
		x := r.Intn(GridWidth-4) + 2 // Leave room for clearing space
		y := r.Intn(GridHeight-4) + 2

		// Check if enemy position is available
		if level[y][x] != TileDirt {
			continue
		}

		// Check if we can clear space around the enemy
		canPlace := true
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				checkX, checkY := x+dx, y+dy
				if checkX < 1 || checkX >= GridWidth-1 || checkY < 1 || checkY >= GridHeight-1 {
					canPlace = false
					break
				}
				// Don't place if there are walls or player/exit nearby
				if level[checkY][checkX] == TileBrickWall || level[checkY][checkX] == TileStoneWall ||
					level[checkY][checkX] == TilePlayer || level[checkY][checkX] == TileClosedExit {
					canPlace = false
					break
				}
			}
			if !canPlace {
				break
			}
		}

		// check if two above the position is empty
		if canPlace && y > 1 && level[y-2][x] != TileDirt || level[y-2][x] == TileRock {
			canPlace = false
		}

		if canPlace {
			// Place enemy
			level[y][x] = tile

			// Clear space around enemy (convert dirt to empty in a cross pattern)
			directions := [][]int{{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}, {-1, 0}}
			for _, dir := range directions {
				nx, ny := x+dir[0], y+dir[1]
				if level[ny][nx] != TileBrickWall {
					level[ny][nx] = TileEmpty
				}
			}
			placed++
		}
	}
}

func resetLevel(n int) {
	currentLevel = levels[n]
	tileMap = currentLevel.Grid
	playerX, playerY = 1, 1
	playerDirection = FacingDown
	playerHoldsFallingObject = false
	gemCounter = 0

	// Initialize enemies from tilemap
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
