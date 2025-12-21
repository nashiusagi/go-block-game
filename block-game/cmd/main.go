package main

import (
	"flag"
	"log"

	"block-game/internal/infrastructure/adapter"
	"block-game/internal/infrastructure/input"
	"block-game/pkg/config"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	flag.Parse()

	baseLayout := config.DefaultLayoutConfig()
	inputPort := input.NewEbitenInputAdapter()
	game := adapter.NewEbitenGame(inputPort)

	ebiten.SetWindowSize(int(baseLayout.ScreenW), int(baseLayout.ScreenH))
	ebiten.SetWindowTitle("Block Game - ブロック崩し")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
