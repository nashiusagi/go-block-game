package adapter

import (
	"testing"

	"block-game/src/application"
	"block-game/src/config"
	"block-game/src/domain"
	"block-game/src/infrastructure/view"
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
