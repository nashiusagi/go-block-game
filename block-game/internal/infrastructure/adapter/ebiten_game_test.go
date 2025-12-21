package adapter

import (
	"testing"

	"block-game/pkg/config"
	"block-game/pkg/domain"
)

// fakeInput satisfies application.InputPort for tests.
type fakeInput struct{}

func (f *fakeInput) Read() domain.InputState { return domain.InputState{} }

func TestLayoutUsesBaseConfigBeforeStart(t *testing.T) {
	game := NewEbitenGame(&fakeInput{})

	w, h := game.Layout(0, 0)
	base := config.DefaultLayoutConfig()
	if w != int(base.ScreenW) || h != int(base.ScreenH) {
		t.Fatalf("unexpected layout size: %d x %d", w, h)
	}
}

func TestStartGameUsesSelectedDifficulty(t *testing.T) {
	game := NewEbitenGame(&fakeInput{})
	game.selectedDiff = domain.DifficultyHard

	if err := game.startGame(); err != nil {
		t.Fatalf("startGame returned error: %v", err)
	}

	if game.usecase == nil || game.renderer == nil {
		t.Fatalf("usecase or renderer not initialized")
	}
	if game.currentLayout().Difficulty != domain.DifficultyHard {
		t.Fatalf("expected HARD difficulty, got %s", game.currentLayout().Difficulty)
	}
}

func TestStartGameFallbackOnInvalidDifficulty(t *testing.T) {
	game := NewEbitenGame(&fakeInput{})
	game.selectedDiff = domain.Difficulty("UNKNOWN")

	if err := game.startGame(); err != nil {
		t.Fatalf("startGame returned error: %v", err)
	}

	if game.currentLayout().Difficulty != domain.DifficultyNormal {
		t.Fatalf("expected fallback to NORMAL, got %s", game.currentLayout().Difficulty)
	}
	if game.statusMsg == "" {
		t.Fatalf("expected status message on fallback, got empty")
	}
}
