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

func createGeneratedLevel(name string, seed int64, rockCount, gemCount int) LevelData {
	var grid [GridHeight][GridWidth]Tile

	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			if y == 0 || y == GridHeight-1 || x == 0 || x == GridWidth-1 {
				grid[y][x] = TileBrickWall
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

func resetLevel(n int) {
	currentLevel = levels[n]
	tileMap = currentLevel.Grid
	playerX, playerY = 1, 1
	playerDirection = FacingDown
	playerHoldsFallingObject = false
	gemCounter = 0
}
