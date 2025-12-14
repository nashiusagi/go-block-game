package main

import (
	"math"
	"testing"

	"block-game/internal/application"
	"block-game/pkg/config"
	"block-game/pkg/domain"
)

func floatEquals(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func rectanglesOverlap(a domain.Block, b domain.Block, w, h float64) bool {
	return a.X < b.X+w && a.X+w > b.X && a.Y < b.Y+h && a.Y+h > b.Y
}

func TestGenerateBlocksDeterministic(t *testing.T) {
	seed := int64(42)
	cfg := config.DefaultLayoutConfig()
	cfg.BlockCount = 12
	cfg.Seed = &seed

	blocks, err := domain.GenerateBlocks(cfg, domain.NewRandomSource(cfg.Seed))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != cfg.BlockCount {
		t.Fatalf("expected %d blocks, got %d", cfg.BlockCount, len(blocks))
	}
	for i, b := range blocks {
		if b.X < 0 || b.X+cfg.BlockW > cfg.ScreenW {
			t.Fatalf("block %d out of X bounds: %.2f", i, b.X)
		}
		if b.Y < 0 || b.Y+cfg.BlockH > cfg.PaddleY-cfg.MinPaddleGap {
			t.Fatalf("block %d out of Y bounds: %.2f", i, b.Y)
		}
		if b.Y > cfg.PaddleY-cfg.MinPaddleGap+1e-6 {
			t.Fatalf("block %d violates min paddle gap: %.2f", i, b.Y)
		}
	}
	for i := 0; i < len(blocks); i++ {
		for j := i + 1; j < len(blocks); j++ {
			if rectanglesOverlap(blocks[i], blocks[j], cfg.BlockW, cfg.BlockH) {
				t.Fatalf("blocks %d and %d overlap", i, j)
			}
		}
	}
	blocks2, err := domain.GenerateBlocks(cfg, domain.NewRandomSource(cfg.Seed))
	if err != nil {
		t.Fatalf("unexpected error on second run: %v", err)
	}
	for i := range blocks {
		if !floatEquals(blocks[i].X, blocks2[i].X, 1e-9) || !floatEquals(blocks[i].Y, blocks2[i].Y, 1e-9) {
			t.Fatalf("deterministic seed mismatch at %d: got (%.6f, %.6f) vs (%.6f, %.6f)", i, blocks[i].X, blocks[i].Y, blocks2[i].X, blocks2[i].Y)
		}
	}
}

func TestGenerateBlocksMaxAttemptsExceeded(t *testing.T) {
	cfg := config.DefaultLayoutConfig()
	cfg.ScreenW = cfg.BlockW * 0.5
	cfg.ScreenH = cfg.BlockH * 0.5
	cfg.BlockCount = 2
	cfg.PaddleY = cfg.BlockH * 2
	cfg.MaxAttempts = 5

	if _, err := domain.GenerateBlocks(cfg, domain.NewRandomSource(cfg.Seed)); err == nil {
		t.Fatalf("expected error due to impossible placement, got nil")
	}
}

type fakeInput struct{}

func (f *fakeInput) Read() domain.InputState {
	return domain.InputState{}
}

func TestGameUsecaseIntegration(t *testing.T) {
	layout := config.DefaultLayoutConfig()
	usecase, err := application.NewGameUsecase(layout, domain.NewRandomSource(layout.Seed), &fakeInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := layout.BlockRows * layout.BlockCols
	if len(usecase.State().Blocks) != expected {
		t.Fatalf("expected %d blocks, got %d", expected, len(usecase.State().Blocks))
	}
	for i, b := range usecase.State().Blocks {
		if b.X < 0 || b.X+layout.BlockW > layout.ScreenW {
			t.Fatalf("block %d out of X bounds: %.2f", i, b.X)
		}
		if b.Y < 0 || b.Y+layout.BlockH > layout.PaddleY-layout.MinPaddleGap {
			t.Fatalf("block %d out of Y bounds: %.2f", i, b.Y)
		}
	}
}
