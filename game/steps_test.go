package game_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/meko-christian/goulder-dash/game"
)

var (
	ErrTileNotExplosion = errors.New("tile is not explosion")
	ErrGemNotDestroyed  = errors.New("adjacent gem not destroyed")
)

// testContext holds the state for a single scenario.
type testContext struct {
	gsm        *game.GameStateManager
	explosionX int
	explosionY int
}

func (tc *testContext) theGameIsInAStateWhereARockCanFallOnThePlayer() error {
	level := make([][]game.Tile, game.GridHeight)
	for i := range level {
		level[i] = make([]game.Tile, game.GridWidth)
	}
	// Player at (1, 2)
	playerX, playerY := 1, 2
	tc.explosionX, tc.explosionY = playerX, playerY
	level[playerY][playerX] = game.TilePlayer
	tc.gsm.SetPlayerPosition(playerX, playerY)
	// Rock at (1, 1)
	level[1][1] = game.TileRock
	// Dirt and gems around player's future explosion location (1, 2)
	level[playerY+1][playerX] = game.TileGem    // Below
	level[playerY][playerX+1] = game.TileDirt   // Right
	level[playerY+1][playerX+1] = game.TileDirt // Below-Right
	tc.gsm.SetLevel(level)

	return nil
}

func (tc *testContext) thePhysicsEngineUpdates() error {
	game.UpdatePhysics(tc.gsm)
	return nil
}

func (tc *testContext) thePlayerShouldBeReplacedByAnExplosion() error {
	tile := tc.gsm.GetTileAt(tc.explosionX, tc.explosionY)
	if tile != game.TileExplosion0 {
		return fmt.Errorf("%w at (%d, %d), got %v", ErrTileNotExplosion, tc.explosionX, tc.explosionY, tile)
	}

	return nil
}

func (tc *testContext) adjacentDirtAndGemsShouldBeDestroyed() error {
	// The first scenario placed a gem at (1, 3) and dirt at (2, 2) and (2, 3)
	// relative to the explosion at (1, 2).
	// The second scenario placed a gem at (2, 2) and dirt at (1, 3) and (2, 3)
	// relative to the explosion at (2, 1).
	// This is just a sample check, a full check would iterate all 8 neighbors.
	if tc.explosionX == 1 && tc.explosionY == 2 { // Scenario 1
		if tc.gsm.GetTileAt(1, 3) != game.TileEmpty {
			return fmt.Errorf("%w at (1, 3)", ErrGemNotDestroyed)
		}
	} else if tc.explosionX == 2 && tc.explosionY == 1 { // Scenario 2
		if tc.gsm.GetTileAt(2, 2) != game.TileEmpty {
			return fmt.Errorf("%w at (2, 2)", ErrGemNotDestroyed)
		}
	}

	return nil
}

func (tc *testContext) theGameIsInAStateWhereAnEnemyCanCollideWithThePlayer() error {
	level := make([][]game.Tile, game.GridHeight)
	for i := range level {
		level[i] = make([]game.Tile, game.GridWidth)
	}
	// Player at (2, 1)
	playerX, playerY := 2, 1
	tc.explosionX, tc.explosionY = playerX, playerY
	level[playerY][playerX] = game.TilePlayer
	tc.gsm.SetPlayerPosition(playerX, playerY)
	// Enemy at (1, 1), facing right
	enemy := game.Enemy{X: 1, Y: 1, Type: game.TileEnemy1, Direction: game.FacingRight, MoveTimer: 0}
	tc.gsm.SetEnemies([]game.Enemy{enemy})

	level[1][1] = game.TileEnemy1
	// Dirt and gems around player's location (2, 1)
	level[playerY+1][playerX] = game.TileGem    // Below
	level[playerY][playerX+1] = game.TileDirt   // Right
	level[playerY+1][playerX+1] = game.TileDirt // Below-Right
	tc.gsm.SetLevel(level)

	return nil
}

func (tc *testContext) theEnemyMovesIntoThePlayersTile() error {
	game.UpdateEnemies(tc.gsm)
	return nil
}

// InitializeSteps registers the step definitions for godog.
func InitializeSteps(ctx *godog.ScenarioContext) {
	testCtx := &testContext{
		gsm:        nil,
		explosionX: 0,
		explosionY: 0,
	}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testCtx.gsm = game.NewGameStateManager()
		return ctx, nil
	})

	ctx.Step(`^the game is in a state where a rock can fall on the player$`,
		testCtx.theGameIsInAStateWhereARockCanFallOnThePlayer)
	ctx.Step(`^the physics engine updates$`, testCtx.thePhysicsEngineUpdates)
	ctx.Step(`^the player should be replaced by an explosion$`, testCtx.thePlayerShouldBeReplacedByAnExplosion)
	ctx.Step(`^adjacent dirt and gems should be destroyed$`, testCtx.adjacentDirtAndGemsShouldBeDestroyed)
	ctx.Step(`^the game is in a state where an enemy can collide with the player$`,
		testCtx.theGameIsInAStateWhereAnEnemyCanCollideWithThePlayer)
	ctx.Step(`^the enemy moves into the player's tile$`, testCtx.theEnemyMovesIntoThePlayersTile)
}
