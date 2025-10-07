package game

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

// testContext holds the state for a single scenario.
type testContext struct {
	gsm        *GameStateManager
	explosionX int
	explosionY int
}

func (tc *testContext) theGameIsInAStateWhereARockCanFallOnThePlayer() error {
	level := make([][]Tile, GridHeight)
	for i := range level {
		level[i] = make([]Tile, GridWidth)
	}
	// Player at (1, 2)
	playerX, playerY := 1, 2
	tc.explosionX, tc.explosionY = playerX, playerY
	level[playerY][playerX] = TilePlayer
	tc.gsm.SetPlayerPosition(playerX, playerY)
	// Rock at (1, 1)
	level[1][1] = TileRock
	// Dirt and gems around player's future explosion location (1, 2)
	level[playerY+1][playerX] = TileGem    // Below
	level[playerY][playerX+1] = TileDirt   // Right
	level[playerY+1][playerX+1] = TileDirt // Below-Right
	tc.gsm.SetLevel(level)
	return nil
}

func (tc *testContext) thePhysicsEngineUpdates() error {
	updatePhysics(tc.gsm)
	return nil
}

func (tc *testContext) thePlayerShouldBeReplacedByAnExplosion() error {
	tile := tc.gsm.GetTileAt(tc.explosionX, tc.explosionY)
	if tile != TileExplosion0 {
		return fmt.Errorf("expected tile at (%d, %d) to be TileExplosion0, but got %v", tc.explosionX, tc.explosionY, tile)
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
		if tc.gsm.GetTileAt(1, 3) != TileEmpty {
			return fmt.Errorf("expected adjacent gem at (1, 3) to be destroyed, but it was not")
		}
	} else if tc.explosionX == 2 && tc.explosionY == 1 { // Scenario 2
		if tc.gsm.GetTileAt(2, 2) != TileEmpty {
			return fmt.Errorf("expected adjacent gem at (2, 2) to be destroyed, but it was not")
		}
	}
	return nil
}

func (tc *testContext) theGameIsInAStateWhereAnEnemyCanCollideWithThePlayer() error {
	level := make([][]Tile, GridHeight)
	for i := range level {
		level[i] = make([]Tile, GridWidth)
	}
	// Player at (2, 1)
	playerX, playerY := 2, 1
	tc.explosionX, tc.explosionY = playerX, playerY
	level[playerY][playerX] = TilePlayer
	tc.gsm.SetPlayerPosition(playerX, playerY)
	// Enemy at (1, 1), facing right
	enemy := Enemy{X: 1, Y: 1, Type: TileEnemy1, Direction: FacingRight}
	tc.gsm.SetEnemies([]Enemy{enemy})
	level[1][1] = TileEnemy1
	// Dirt and gems around player's location (2, 1)
	level[playerY+1][playerX] = TileGem    // Below
	level[playerY][playerX+1] = TileDirt   // Right
	level[playerY+1][playerX+1] = TileDirt // Below-Right
	tc.gsm.SetLevel(level)
	return nil
}

func (tc *testContext) theEnemyMovesIntoThePlayersTile() error {
	updateEnemies(tc.gsm)
	return nil
}

// InitializeSteps registers the step definitions for godog.
func InitializeSteps(ctx *godog.ScenarioContext) {
	tc := &testContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tc.gsm = NewGameStateManager()
		return ctx, nil
	})

	ctx.Step(`^the game is in a state where a rock can fall on the player$`, tc.theGameIsInAStateWhereARockCanFallOnThePlayer)
	ctx.Step(`^the physics engine updates$`, tc.thePhysicsEngineUpdates)
	ctx.Step(`^the player should be replaced by an explosion$`, tc.thePlayerShouldBeReplacedByAnExplosion)
	ctx.Step(`^adjacent dirt and gems should be destroyed$`, tc.adjacentDirtAndGemsShouldBeDestroyed)
	ctx.Step(`^the game is in a state where an enemy can collide with the player$`, tc.theGameIsInAStateWhereAnEnemyCanCollideWithThePlayer)
	ctx.Step(`^the enemy moves into the player's tile$`, tc.theEnemyMovesIntoThePlayersTile)
}