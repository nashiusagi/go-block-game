package config

import (
	"block-game/pkg/domain"
	"testing"
)

func TestDefaultLayoutConfigValues(t *testing.T) {
	cfg := DefaultLayoutConfig()

	if cfg.BlockCount != BlockRows*BlockCols {
		t.Fatalf("unexpected block count: %d", cfg.BlockCount)
	}
	if cfg.MaxAttempts != MaxAttemptsFactor*BlockRows*BlockCols {
		t.Fatalf("unexpected max attempts: %d", cfg.MaxAttempts)
	}
	if cfg.ScreenW != ScreenWidth || cfg.ScreenH != ScreenHeight {
		t.Fatalf("unexpected screen size: %f x %f", cfg.ScreenW, cfg.ScreenH)
	}
	if cfg.PaddleY != PaddleY {
		t.Fatalf("unexpected paddle Y: %f", cfg.PaddleY)
	}
	if cfg.MinPaddleGap != MinPaddleGap {
		t.Fatalf("unexpected min paddle gap: %f", cfg.MinPaddleGap)
	}
}

func TestLayoutWithDifficultyHard(t *testing.T) {
	cfg, applied, err := LayoutWithDifficulty("HARD")
	if err != nil && applied == domain.DifficultyNormal {
		t.Fatalf("unexpected failure applying HARD: %v", err)
	}
	if applied != domain.DifficultyHard {
		t.Fatalf("expected applied HARD, got %s", applied)
	}
	if cfg.BallSpeed <= BallSpeed {
		t.Fatalf("ball speed not increased for HARD")
	}
	if cfg.BlockCount <= BlockRows*BlockCols {
		t.Fatalf("block count not increased for HARD")
	}
}

func TestLayoutWithDifficultyInvalidFallsBack(t *testing.T) {
	cfg, applied, err := LayoutWithDifficulty("UNKNOWN")
	if applied != domain.DifficultyNormal {
		t.Fatalf("expected fallback to NORMAL, got %s", applied)
	}
	if err == nil {
		t.Fatalf("expected validation error for invalid difficulty")
	}
	// Should still return a valid config (defaults applied).
	if cfg.Difficulty != domain.DifficultyNormal {
		t.Fatalf("config difficulty should be NORMAL on fallback")
	}
}
