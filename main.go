package main

import (
	"github.com/gonutz/prototype/draw"
	"github.com/cwbudde/goulder-dash/game"
)

func main() {
	err := draw.RunWindow("Goulder Dash", game.WindowWidth, game.WindowHeight, game.Update)
	if err != nil {
		panic(err)
	}
}
