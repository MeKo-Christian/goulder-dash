package main

import (
	"github.com/gonutz/prototype/draw"
	"github.com/meko-christian/goulder-dash/game"
)

func main() {
	draw.RunWindow("Goulder Dash", game.WindowWidth, game.WindowHeight, game.Update)
}
