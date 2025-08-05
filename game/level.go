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

	for y := range GridHeight {
		for x := range GridWidth {
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

func placeEnemiesWithSpace(level *[GridHeight][GridWidth]Tile, randGen *rand.Rand, tile Tile, count int) {
	placed := 0
	attempts := 0
	maxAttempts := count * 50 // Prevent infinite loops

	for placed < count && attempts < maxAttempts {
		attempts++
		x := randGen.Intn(GridWidth-4) + 2 // Leave room for clearing space
		y := randGen.Intn(GridHeight-4) + 2

		if canPlaceEnemyAt(level, x, y) {
			level[y][x] = tile
			clearSpaceAroundEnemy(level, x, y)

			placed++
		}
	}
}

func canPlaceEnemyAt(level *[GridHeight][GridWidth]Tile, x, y int) bool {
	// Check if enemy position is available
	if level[y][x] != TileDirt {
		return false
	}

	// Check if we can clear space around the enemy
	if !canClearSpaceAround(level, x, y) {
		return false
	}

	// Check if two above the position is empty
	if y > 1 && (level[y-2][x] != TileDirt || level[y-2][x] == TileRock) {
		return false
	}

	return true
}

func canClearSpaceAround(level *[GridHeight][GridWidth]Tile, x, y int) bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			checkX, checkY := x+dx, y+dy
			if !isValidPosition(checkX, checkY) || isBlockingTile(level[checkY][checkX]) {
				return false
			}
		}
	}

	return true
}

func isValidPosition(x, y int) bool {
	return x >= 1 && x < GridWidth-1 && y >= 1 && y < GridHeight-1
}

func isBlockingTile(tile Tile) bool {
	return tile == TileBrickWall || tile == TileStoneWall ||
		tile == TilePlayer || tile == TileClosedExit
}

func clearSpaceAroundEnemy(level *[GridHeight][GridWidth]Tile, x, y int) {
	directions := [][]int{{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}, {-1, 0}}
	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if level[ny][nx] != TileBrickWall {
			level[ny][nx] = TileEmpty
		}
	}
}
