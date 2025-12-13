package application

import (
	"testing"

	"block-game/src/config"
	"block-game/src/domain"
)

type fakeInput struct {
	state domain.InputState
}

func (f *fakeInput) Read() domain.InputState {
	return f.state
}

func TestNewGameUsecaseNilInput(t *testing.T) {
	cfg := config.DefaultLayoutConfig()
	if _, err := NewGameUsecase(cfg, domain.NewRandomSource(cfg.Seed), nil); err == nil {
		t.Fatalf("expected error when input is nil")
	}
}

func TestNewGameUsecaseFallbackOnError(t *testing.T) {
	cfg := config.DefaultLayoutConfig()
	cfg.MinPaddleGap = cfg.PaddleY + 10 // 強制的に生成失敗させる

	usecase, err := NewGameUsecase(cfg, domain.NewRandomSource(cfg.Seed), &fakeInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := cfg.BlockRows * cfg.BlockCols
	if len(usecase.State().Blocks) != expected {
		t.Fatalf("expected fallback blocks: %d got %d", expected, len(usecase.State().Blocks))
	}
}

func TestUpdateAppliesInput(t *testing.T) {
	cfg := config.DefaultLayoutConfig()
	fi := &fakeInput{state: domain.InputState{MoveRight: true}}
	usecase, err := NewGameUsecase(cfg, domain.NewRandomSource(cfg.Seed), fi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	initialX := usecase.State().Paddle.X
	if err := usecase.Update(); err != nil {
		t.Fatalf("update error: %v", err)
	}
	if usecase.State().Paddle.X <= initialX {
		t.Fatalf("expected paddle to move right")
	}
}
