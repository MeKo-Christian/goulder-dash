package main

import (
	"github.com/gonutz/prototype/draw"
	"github.com/meko-christian/goulder-dash/game"
)

func main() {
	err := draw.RunWindow("Goulder Dash", game.WindowWidth, game.WindowHeight, game.Update)
	if err != nil {
		panic(err)
	}
}
