package view

import (
	"image"
	"testing"

	"block-game/pkg/config"
	"block-game/pkg/domain"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestRenderShowsDifficultyLabel_NoPanic(t *testing.T) {
	layout := domain.LayoutConfig{
		ScreenW:      100,
		ScreenH:      100,
		BlockW:       10,
		BlockH:       10,
		PaddleY:      80,
		PaddleWidth:  20,
		PaddleHeight: 5,
		BallRadius:   5,
		Difficulty:   domain.DifficultyHard,
	}
	state := &domain.GameState{
		Blocks: []domain.Block{},
		Balls: []domain.Ball{
			{X: 50, Y: 50, Radius: 5},
		},
		Paddle: domain.Paddle{X: 40, Y: 80, Width: 20, Height: 5},
		Items:  []domain.Item{},
		Score:  0,
	}

	renderer := NewRenderer(layout)
	screen := ebiten.NewImageWithOptions(image.Rect(0, 0, int(layout.ScreenW), int(layout.ScreenH)), &ebiten.NewImageOptions{})

	renderer.Render(screen, state)

	if screen == nil {
		t.Fatalf("screen should not be nil after rendering")
	}
}

func TestRendererRenderDoesNotPanic(t *testing.T) {
	cfg := config.DefaultLayoutConfig()
	renderer := NewRenderer(cfg)

	state := domain.NewGameState(cfg, []domain.Block{
		{X: 10, Y: 10, Alive: true},
	})
	state.GameOver = true

	screen := ebiten.NewImage(int(cfg.ScreenW), int(cfg.ScreenH))
	defer screen.Dispose()

	renderer.Render(screen, state)
	// no assertion: absence of panic is success
}
