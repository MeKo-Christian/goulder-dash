package game

type PlayerSystem struct {
	gsm               *GameStateManager
	keyRepeatTimer    int
	keyRepeatInterval int
	lastKey           int // Track which key was last pressed
}

func NewPlayerSystem(gsm *GameStateManager) *PlayerSystem {
	return &PlayerSystem{
		gsm:               gsm,
		keyRepeatTimer:    0,
		keyRepeatInterval: 0,
		lastKey:           0,
	}
}

func (ps *PlayerSystem) GetPosition() (int, int) {
	return ps.gsm.GetPlayerPosition()
}

func (ps *PlayerSystem) GetDirection() Direction {
	return ps.gsm.GetPlayerDirection()
}

func (ps *PlayerSystem) SetDirection(direction Direction) {
	ps.gsm.SetPlayerDirection(direction)
}

func (ps *PlayerSystem) UpdateKeyRepeat() {
	if ps.keyRepeatTimer > 0 {
		ps.keyRepeatTimer--
	}

	if ps.keyRepeatInterval > 0 {
		ps.keyRepeatInterval--
	}
}

func (ps *PlayerSystem) StartKeyRepeat(keyCode int) {
	ps.lastKey = keyCode
	ps.keyRepeatTimer = 15 // Initial delay of 15 frames (~250ms at 60fps)
	ps.keyRepeatInterval = 0
}

func (ps *PlayerSystem) CanRepeatKey(keyCode int) bool {
	return ps.lastKey == keyCode && ps.keyRepeatTimer == 0 && ps.keyRepeatInterval == 0
}

func (ps *PlayerSystem) SetRepeatInterval() {
	ps.keyRepeatInterval = 4 // 4 frames between repeats (~67ms at 60fps)
}

func (ps *PlayerSystem) Move(dx, dy int) bool {
	x, y := ps.GetPosition()
	newX, newY := x+dx, y+dy

	target := ps.gsm.GetTileAt(newX, newY)

	if !ps.canMoveTo(target, newX, newY, dx, dy) {
		return false
	}

	ps.moveTo(newX, newY, target, dx, dy)

	return true
}

func (ps *PlayerSystem) IsExploding() bool {
	x, y := ps.GetPosition()
	tile := ps.gsm.GetTileAt(x, y)

	return tile >= TileExplosion0 && tile <= TileExplosion5
}

func (ps *PlayerSystem) canMoveTo(target Tile, newX, newY, dx, dy int) bool {
	// Walls are always blocked
	if target == TileBrickWall || target == TileStoneWall {
		return false
	}

	// Prevent entry into closed exit
	if target == TileClosedExit {
		return false
	}

	// Handle pushing rock
	if target == TileRock {
		return ps.canPushRock(newX, newY, dx, dy)
	}

	return true
}

func (ps *PlayerSystem) canPushRock(newX, newY, dx, dy int) bool {
	// Only allow horizontal pushing
	if dy != 0 {
		return false
	}

	pushX := newX + dx
	pushY := newY

	return ps.gsm.GetTileAt(pushX, pushY) == TileEmpty
}

func (ps *PlayerSystem) moveTo(newX, newY int, target Tile, dx, _ int) {
	// Transition if player enters open exit
	if target == TileOpenExit {
		ps.gsm.LoadNextLevel()
		return
	}

	// Handle enemy collision
	if target == TileEnemy1 || target == TileEnemy2 || target == TileEnemy3 {
		// Player dies when walking into enemy
		ps.gsm.SetTileAt(newX, newY, TileExplosion0)
		return
	}

	rockMoved := ps.handleRockPushing(newX, newY, dx, target)
	ps.handleGemCollection(target)
	ps.updateSupport(newX, newY, target, rockMoved)

	// Move player
	playerX, playerY := ps.GetPosition()
	ps.gsm.SetTileAt(playerX, playerY, TileEmpty)
	ps.gsm.SetPlayerPosition(newX, newY)
	ps.gsm.SetTileAt(newX, newY, TilePlayer)
}

func (ps *PlayerSystem) handleRockPushing(newX, newY, dx int, target Tile) bool {
	if target != TileRock {
		return false
	}

	pushX := newX + dx
	pushY := newY

	// Move rock
	ps.gsm.SetTileAt(pushX, pushY, TileRock)
	ps.gsm.SetTileAt(newX, newY, TileEmpty)

	return true
}

func (ps *PlayerSystem) handleGemCollection(target Tile) {
	if target == TileGem {
		ps.gsm.IncrementGemCounter()

		// Check if all gems collected
		if ps.gsm.GetGemCounter() >= ps.gsm.GetCurrentLevel().GemTarget {
			// Open the exit
			for y := range GridHeight {
				for x := range GridWidth {
					if ps.gsm.GetTileAt(x, y) == TileClosedExit {
						ps.gsm.SetTileAt(x, y, TileOpenExit)
					}
				}
			}
		}
	}
}

func (ps *PlayerSystem) updateSupport(newX, newY int, target Tile, rockMoved bool) {
	// Reset on each move
	ps.gsm.SetPlayerHoldsFallingObject(false)

	// Check for support case
	if target == TileDirt || target == TileGem || rockMoved {
		if newY > 0 {
			above := ps.gsm.GetTileAt(newX, newY-1)
			if above == TileRock || above == TileGem {
				ps.gsm.SetPlayerHoldsFallingObject(true)
			}
		}
	}
}
