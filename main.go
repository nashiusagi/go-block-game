package main

import (
	"log"

	"block-game/src/application"
	"block-game/src/config"
	"block-game/src/domain"
	"block-game/src/infrastructure/adapter"
	"block-game/src/infrastructure/input"
	"block-game/src/infrastructure/view"

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
