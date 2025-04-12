package main

import (
	"math/rand"

	"github.com/gonutz/prototype/draw"
)

const (
	tileSize     = 64 // original sprite tile size
	tileDrawSize = 32 // scale down to 32x32 for smoother view
	tileCols     = 6
	tileRows     = 7
	gridWidth    = 25
	gridHeight   = 20
	spriteSheet  = "assets/sprites.png"
)

type Tile int

const (
	TileEmpty Tile = iota
	TileDirt
	TileBrickWall
	TileStoneWall
	TileRock
	TileGem
	TileClosedExit
	TileOpenExit
	TilePlayer
	TileEnemy1
	TileEnemy2
	TileEnemy3
	TileExplosionStart // used for animation
	TileExplosionEnd   // just for range logic
)

type Direction int

const (
	FacingRight Direction = iota
	FacingDown
	FacingLeft
	FacingUp
)

var tileSpriteIndex = map[Tile]int{
	TileEmpty:          3 + 6*tileCols,
	TileDirt:           4 + 2*tileCols,
	TileBrickWall:      4 + 0*tileCols,
	TileStoneWall:      4 + 1*tileCols,
	TileRock:           5 + 0*tileCols,
	TileGem:            5 + 3*tileCols,
	TileClosedExit:     5 + 1*tileCols,
	TileOpenExit:       5 + 2*tileCols,
	TilePlayer:         0, // Default: right-facing
	TileEnemy1:         0 + 1*tileCols,
	TileEnemy2:         0 + 2*tileCols,
	TileEnemy3:         0 + 3*tileCols,
	TileExplosionStart: 0 + 4*tileCols,
	TileExplosionEnd:   5 + 4*tileCols,
}

var playerX, playerY int
var playerDirection = FacingDown

var frameCounter int

var tileIndices [gridHeight][gridWidth]int
var tileMap [gridHeight][gridWidth]Tile

var level = [gridHeight][gridWidth]Tile{}

func setupLevel1() [gridHeight][gridWidth]Tile {
	var level [gridHeight][gridWidth]Tile

	// Fill map with borders and dirt
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			if y == 0 || y == gridHeight-1 || x == 0 || x == gridWidth-1 {
				level[y][x] = TileBrickWall
			} else {
				level[y][x] = TileDirt
			}
		}
	}

	playerX = 1
	playerY = 1

	// Place player at [1][1]
	level[playerX][playerY] = TilePlayer

	// Place exit near bottom-right, just inside the wall
	level[gridHeight-2][gridWidth-2] = TileClosedExit

	// Seed RNG for consistent placement
	r := rand.New(rand.NewSource(42))

	// Place 40 rocks
	placeRandomTiles(&level, r, TileRock, 40)

	// Place 20 gems
	placeRandomTiles(&level, r, TileGem, 20)

	return level
}

func placeRandomTiles(level *[gridHeight][gridWidth]Tile, r *rand.Rand, tile Tile, count int) {
	placed := 0
	for placed < count {
		x := r.Intn(gridWidth-2) + 1
		y := r.Intn(gridHeight-2) + 1

		if level[y][x] == TileDirt {
			level[y][x] = tile
			placed++
		}
	}
}

func setupLevel2() [gridHeight][gridWidth]Tile {
	var level [gridHeight][gridWidth]Tile

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			if y == 0 || y == gridHeight-1 || x == 0 || x == gridWidth-1 {
				level[y][x] = TileBrickWall
			} else {
				level[y][x] = TileDirt
			}
		}
	}

	// Place player at [1][1]
	level[1][1] = TilePlayer

	// TODO

	return level
}

func loadLevel(n int) [gridHeight][gridWidth]Tile {
	switch n {
	case 1:
		return setupLevel1()
	default:
		return setupLevel1()
	}
}

func handlePlayerMovement(w draw.Window) {
	dx, dy := 0, 0
	moved := false

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
		} else {
			// Can't push if not empty behind
			return
		}
	}

	// Move player
	tileMap[playerY][playerX] = TileEmpty
	playerX = newX
	playerY = newY
	tileMap[playerY][playerX] = TilePlayer
}

func updatePhysics() {
	// Process bottom-up
	for y := gridHeight - 2; y >= 1; y-- {
		for x := 1; x < gridWidth-1; x++ {
			tile := tileMap[y][x]
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
		}
	}
}

func main() {
	windowWidth := gridWidth * tileDrawSize
	windowHeight := gridHeight * tileDrawSize

	tileMap = loadLevel(1)

	draw.RunWindow("Goulder Dash", windowWidth, windowHeight, func(w draw.Window) {
		w.BlurImages(true) // enable smoothing

		handlePlayerMovement(w)

		frameCounter++
		if frameCounter%10 == 0 {
			updatePhysics()
		}

		for y := 0; y < gridHeight; y++ {
			for x := 0; x < gridWidth; x++ {
				tile := tileMap[y][x]
				spriteIndex := tileSpriteIndex[tile]

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
				} else {
					spriteIndex = tileSpriteIndex[tile]
				}

				sx := (spriteIndex % tileCols) * tileSize
				sy := (spriteIndex / tileCols) * tileSize

				err := w.DrawImageFilePart(
					spriteSheet,
					sx, sy, tileSize, tileSize,
					x*tileDrawSize, y*tileDrawSize, tileDrawSize, tileDrawSize,
					0,
				)

				if err != nil {
					w.DrawText("Failed to load sprite!", 10, 10, draw.Red)
				}
			}
		}
	})
}
