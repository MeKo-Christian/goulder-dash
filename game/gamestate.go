package game

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPlayerPosition = errors.New("invalid player position")
	ErrInvalidLevelIndex     = errors.New("invalid level index")
	ErrInvalidGemCounter     = errors.New("invalid gem counter")
	ErrInvalidEnemyPosition  = errors.New("invalid enemy position")
)

type PlayerState struct {
	X, Y               int
	Direction          Direction
	HoldsFallingObject bool
}

type GameState struct {
	Player            PlayerState
	TileMap           [GridHeight][GridWidth]Tile
	FrameCounter      int
	GemCounter        int
	CurrentLevel      LevelData
	CurrentLevelIndex int
	Enemies           []Enemy
	Levels            []LevelData
}

type GameStateManager struct {
	state *GameState
}

func NewGameStateManager() *GameStateManager {
	levels := []LevelData{
		createGeneratedLevel("Level 1", 42, 40, 20, 1),
		createGeneratedLevel("Level 2", 32, 45, 25, 0),
		createGeneratedLevel("Level 3", 12, 50, 30, 3),
		createGeneratedLevel("Level 4", 22, 55, 35, 4),
		createGeneratedLevel("Level 5", 52, 60, 40, 5),
		createGeneratedLevel("Level 6", 62, 65, 45, 6),
	}

	gameState := &GameState{
		Player: PlayerState{
			X:                  1,
			Y:                  1,
			Direction:          FacingDown,
			HoldsFallingObject: false,
		},
		TileMap:           levels[0].Grid,
		FrameCounter:      0,
		GemCounter:        0,
		CurrentLevel:      levels[0],
		CurrentLevelIndex: 0,
		Enemies:           []Enemy{},
		Levels:            levels,
	}

	manager := &GameStateManager{state: gameState}
	manager.initializeEnemies()

	return manager
}

func (gsm *GameStateManager) GetState() *GameState {
	return gsm.state
}

func (gsm *GameStateManager) GetPlayerPosition() (int, int) {
	return gsm.state.Player.X, gsm.state.Player.Y
}

func (gsm *GameStateManager) GetPlayerDirection() Direction {
	return gsm.state.Player.Direction
}

func (gsm *GameStateManager) GetPlayerHoldsFallingObject() bool {
	return gsm.state.Player.HoldsFallingObject
}

func (gsm *GameStateManager) GetTileAt(x, y int) Tile {
	if x < 0 || x >= GridWidth || y < 0 || y >= GridHeight {
		return TileStoneWall
	}

	return gsm.state.TileMap[y][x]
}

func (gsm *GameStateManager) SetTileAt(x, y int, tile Tile) bool {
	if !gsm.isValidPosition(x, y) {
		return false
	}

	gsm.state.TileMap[y][x] = tile

	return true
}

func (gsm *GameStateManager) SetPlayerPosition(x, y int) bool {
	if !gsm.isValidPosition(x, y) {
		return false
	}

	gsm.state.Player.X = x
	gsm.state.Player.Y = y

	return true
}

func (gsm *GameStateManager) SetPlayerDirection(direction Direction) {
	gsm.state.Player.Direction = direction
}

func (gsm *GameStateManager) SetPlayerHoldsFallingObject(holds bool) {
	gsm.state.Player.HoldsFallingObject = holds
}

func (gsm *GameStateManager) IncrementFrameCounter() {
	gsm.state.FrameCounter++
}

func (gsm *GameStateManager) GetFrameCounter() int {
	return gsm.state.FrameCounter
}

func (gsm *GameStateManager) IncrementGemCounter() {
	gsm.state.GemCounter++
}

func (gsm *GameStateManager) GetGemCounter() int {
	return gsm.state.GemCounter
}

func (gsm *GameStateManager) GetCurrentLevel() LevelData {
	return gsm.state.CurrentLevel
}

func (gsm *GameStateManager) GetEnemies() []Enemy {
	return gsm.state.Enemies
}

func (gsm *GameStateManager) SetEnemies(enemies []Enemy) {
	gsm.state.Enemies = enemies
}

func (gsm *GameStateManager) SetLevel(level [][]Tile) {
	var newTileMap [GridHeight][GridWidth]Tile
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			newTileMap[y][x] = level[y][x]
		}
	}
	gsm.state.TileMap = newTileMap
}

func (gsm *GameStateManager) GetTileMap() *[GridHeight][GridWidth]Tile {
	return &gsm.state.TileMap
}

func (gsm *GameStateManager) LoadNextLevel() {
	gsm.state.CurrentLevelIndex++
	if gsm.state.CurrentLevelIndex >= len(gsm.state.Levels) {
		gsm.state.CurrentLevelIndex = 0
	}

	gsm.resetLevel(gsm.state.CurrentLevelIndex)
}

func (gsm *GameStateManager) ResetCurrentLevel() {
	gsm.resetLevel(gsm.state.CurrentLevelIndex)
}

func (gsm *GameStateManager) ValidateState() error {
	// Validate player position
	if !gsm.isValidPosition(gsm.state.Player.X, gsm.state.Player.Y) {
		return fmt.Errorf("%w: (%d, %d)", ErrInvalidPlayerPosition, gsm.state.Player.X, gsm.state.Player.Y)
	}

	// Validate current level index
	if gsm.state.CurrentLevelIndex < 0 || gsm.state.CurrentLevelIndex >= len(gsm.state.Levels) {
		return fmt.Errorf("%w: %d", ErrInvalidLevelIndex, gsm.state.CurrentLevelIndex)
	}

	// Validate gem counter
	if gsm.state.GemCounter < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidGemCounter, gsm.state.GemCounter)
	}

	// Validate enemies
	for i, enemy := range gsm.state.Enemies {
		if !gsm.isValidPosition(enemy.X, enemy.Y) {
			return fmt.Errorf("%w %d: (%d, %d)", ErrInvalidEnemyPosition, i, enemy.X, enemy.Y)
		}
	}

	return nil
}

func (gsm *GameStateManager) resetLevel(levelIndex int) {
	gsm.state.CurrentLevel = gsm.state.Levels[levelIndex]
	gsm.state.TileMap = gsm.state.CurrentLevel.Grid
	gsm.state.Player.X = 1
	gsm.state.Player.Y = 1
	gsm.state.Player.Direction = FacingDown
	gsm.state.Player.HoldsFallingObject = false
	gsm.state.GemCounter = 0
	gsm.initializeEnemies()
}

func (gsm *GameStateManager) initializeEnemies() {
	gsm.state.Enemies = []Enemy{}

	for y := range GridHeight {
		for x := range GridWidth {
			if gsm.state.TileMap[y][x] == TileEnemy1 {
				gsm.state.Enemies = append(gsm.state.Enemies, Enemy{
					X:         x,
					Y:         y,
					Type:      TileEnemy1,
					Direction: FacingRight,
					MoveTimer: 8,
				})
			}
		}
	}
}

func (gsm *GameStateManager) isValidPosition(x, y int) bool {
	return x >= 0 && x < GridWidth && y >= 0 && y < GridHeight
}
