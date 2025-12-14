package adapter

import (
	"testing"

	"block-game/internal/application"
	"block-game/internal/infrastructure/view"
	"block-game/pkg/config"
	"block-game/pkg/domain"
)

type fakeRenderer struct {
	calls int
}

func (f *fakeRenderer) Render(screen interface{}, state *domain.GameState) {
	f.calls++
}

func TestEbitenGameLayout(t *testing.T) {
	cfg := config.DefaultLayoutConfig()
	usecase, err := application.NewGameUsecase(cfg, domain.NewRandomSource(cfg.Seed), &fakeInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := view.NewRenderer(cfg)
	game := NewEbitenGame(usecase, r)

	w, h := game.Layout(0, 0)
	if w != int(cfg.ScreenW) || h != int(cfg.ScreenH) {
		t.Fatalf("unexpected layout size: %d x %d", w, h)
	}
}

// fakeInput is reused from application tests.
type fakeInput struct{}

func (f *fakeInput) Read() domain.InputState { return domain.InputState{} }
