package view

import (
	"testing"

	"block-game/src/config"
	"block-game/src/domain"
	"github.com/hajimehoshi/ebiten/v2"
)

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
