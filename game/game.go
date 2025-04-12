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
	levels                   = []LevelData{
		createGeneratedLevel("Level 1", 42, 40, 20),
		createGeneratedLevel("Level 2", 32, 45, 25),
		createGeneratedLevel("Level 3", 12, 50, 30),
		createGeneratedLevel("Level 4", 22, 55, 35),
		createGeneratedLevel("Level 5", 52, 60, 40),
		createGeneratedLevel("Level 6", 62, 65, 45),
	}
)

func init() {
	currentLevel = levels[0]
	tileMap = currentLevel.Grid
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

func updatePhysics() {
	// Process bottom-up
	for y := GridHeight - 2; y >= 1; y-- {
		for x := 1; x < GridWidth-1; x++ {
			tile := tileMap[y][x]

			if tile >= TileExplosion0 && tile < TileExplosion5 {
				tileMap[y][x]++ // next frame
			} else if tile == TileExplosion5 {
				// reset level
				resetLevel(currentLevelIndex)
			}

			// Skip empty tiles
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
