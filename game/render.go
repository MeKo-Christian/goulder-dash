package game

import (
	"strconv"

	"github.com/gonutz/prototype/draw"
)

func renderGame(window draw.Window, gsm *GameStateManager) {
	tileMap := gsm.GetTileMap()
	playerDirection := gsm.GetPlayerDirection()

	for y := range GridHeight {
		for x := range GridWidth {
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

			err := window.DrawImageFilePart(
				SpriteSheet,
				sx, sy, TileSize, TileSize,
				x*TileDrawSize, y*TileDrawSize, TileDrawSize, TileDrawSize,
				0,
			)
			if err != nil {
				window.DrawText("Failed to load sprite!", 10, 10, draw.Red)
			}
		}
	}

	currentLevel := gsm.GetCurrentLevel()
	gemCounter := gsm.GetGemCounter()
	text := currentLevel.Name + " - Gems: " + strconv.Itoa(gemCounter) + " / " + strconv.Itoa(currentLevel.GemTarget)
	window.DrawText(text, 8, 8, draw.White)
}
