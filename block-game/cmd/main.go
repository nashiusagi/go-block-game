package main

import (
	"log"

	"block-game/application"
	"block-game/config"
	"block-game/domain"
	"block-game/infrastructure/adapter"
	"block-game/infrastructure/input"
	"block-game/infrastructure/view"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	layout := config.DefaultLayoutConfig()
	rnd := domain.NewRandomSource(layout.Seed)
	inputPort := input.NewEbitenInputAdapter()

	usecase, err := application.NewGameUsecase(layout, rnd, inputPort)
	if err != nil {
		log.Fatalf("failed to initialize game usecase: %v", err)
	}

	renderer := view.NewRenderer(layout)
	game := adapter.NewEbitenGame(usecase, renderer)

	ebiten.SetWindowSize(int(layout.ScreenW), int(layout.ScreenH))
	ebiten.SetWindowTitle("Block Game - ブロック崩し")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
