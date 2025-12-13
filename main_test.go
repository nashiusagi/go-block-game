package main

import (
	"math"
	"testing"
)

func floatEquals(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func rectanglesOverlap(a Block, b Block, w, h float64) bool {
	return rectsOverlap(a.X, a.Y, w, h, b.X, b.Y, w, h)
}

func TestGenerateBlocksDeterministic(t *testing.T) {
	seed := int64(42)
	cfg := LayoutConfig{
		ScreenW:      screenWidth,
		ScreenH:      screenHeight,
		BlockW:       blockWidth,
		BlockH:       blockHeight,
		BlockCount:   12,
		MinPaddleGap: minPaddleGap,
		PaddleY:      float64(screenHeight - 50),
		MaxAttempts:  500,
		Seed:         &seed,
	}
	blocks, err := GenerateBlocks(cfg, newRandomSource(cfg.Seed))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != cfg.BlockCount {
		t.Fatalf("expected %d blocks, got %d", cfg.BlockCount, len(blocks))
	}
	// 全ブロックが境界内・ギャップ内にあること
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
	// 重なりなしを検証
	for i := 0; i < len(blocks); i++ {
		for j := i + 1; j < len(blocks); j++ {
			if rectanglesOverlap(blocks[i], blocks[j], cfg.BlockW, cfg.BlockH) {
				t.Fatalf("blocks %d and %d overlap", i, j)
			}
		}
	}
	// 決定的シードで再現性を確認
	blocks2, err := GenerateBlocks(cfg, newRandomSource(cfg.Seed))
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
	cfg := LayoutConfig{
		ScreenW:      blockWidth * 0.5, // 意図的に狭い
		ScreenH:      blockHeight * 0.5,
		BlockW:       blockWidth,
		BlockH:       blockHeight,
		BlockCount:   2,
		MinPaddleGap: minPaddleGap,
		PaddleY:      blockHeight * 2, // 低い位置でギャップが満たせない
		MaxAttempts:  5,
		Seed:         nil,
	}
	if _, err := GenerateBlocks(cfg, newRandomSource(cfg.Seed)); err == nil {
		t.Fatalf("expected error due to impossible placement, got nil")
	}
}

func TestInitBlocksIntegration(t *testing.T) {
	g := &Game{}
	g.initPaddle()
	g.initBlocks()

	expected := blockRows * blockCols
	if len(g.blocks) != expected {
		t.Fatalf("expected %d blocks, got %d", expected, len(g.blocks))
	}
	for i, b := range g.blocks {
		if b.X < 0 || b.X+blockWidth > screenWidth {
			t.Fatalf("block %d out of X bounds: %.2f", i, b.X)
		}
		if b.Y < 0 || b.Y+blockHeight > g.paddle.Y-minPaddleGap {
			t.Fatalf("block %d out of Y bounds: %.2f", i, b.Y)
		}
	}
}
