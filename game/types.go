package game

const (
	GridWidth  = 25
	GridHeight = 20

	TileSize     = 64
	TileDrawSize = 64
	TileCols     = 6
	TileRows     = 7
	SpriteSheet  = "assets/sprites.png"
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
	TileExplosion0
	TileExplosion1
	TileExplosion2
	TileExplosion3
	TileExplosion4
	TileExplosion5
)

type Direction int

const (
	FacingRight Direction = iota
	FacingDown
	FacingLeft
	FacingUp
)

var TileSpriteIndex = map[Tile]int{
	TileEmpty:      3 + 6*TileCols,
	TileDirt:       4 + 2*TileCols,
	TileBrickWall:  4 + 0*TileCols,
	TileStoneWall:  4 + 1*TileCols,
	TileRock:       5 + 0*TileCols,
	TileGem:        5 + 3*TileCols,
	TileClosedExit: 5 + 1*TileCols,
	TileOpenExit:   5 + 2*TileCols,
	TilePlayer:     0,
	TileEnemy1:     0 + 1*TileCols,
	TileEnemy2:     0 + 2*TileCols,
	TileEnemy3:     0 + 3*TileCols,
	TileExplosion0: 0 + 4*TileCols,
	TileExplosion1: 1 + 4*TileCols,
	TileExplosion2: 2 + 4*TileCols,
	TileExplosion3: 3 + 4*TileCols,
	TileExplosion4: 4 + 4*TileCols,
	TileExplosion5: 5 + 4*TileCols,
}

var (
	WindowWidth  = GridWidth * TileDrawSize
	WindowHeight = GridHeight * TileDrawSize
)
