package config

import "testing"

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
