package main

import (
	"flag"
	"log"

	"block-game/internal/application"
	"block-game/internal/infrastructure/adapter"
	"block-game/internal/infrastructure/input"
	"block-game/internal/infrastructure/view"
	"block-game/pkg/config"
	"block-game/pkg/domain"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	diffFlag := flag.String("difficulty", string(domain.DifficultyNormal), "difficulty: EASY|NORMAL|HARD")
	flag.Parse()

	layout, applied, err := config.LayoutWithDifficulty(*diffFlag)
	if err != nil {
		log.Printf("difficulty selection error (%v), applied: %s", err, applied)
	}

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
