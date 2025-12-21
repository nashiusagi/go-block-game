package domain

import (
	"testing"
)

type mockRandom struct {
	floats []float64
	idx    int
}

func (m *mockRandom) Float64() float64 {
	val := m.floats[m.idx%len(m.floats)]
	m.idx++
	return val
}

func (m *mockRandom) Intn(n int) int {
	return 0
}

func (m *mockRandom) Seed(seed int64) {}

func baseLayout() LayoutConfig {
	return LayoutConfig{
		ScreenW:        800,
		ScreenH:        600,
		BlockW:         70,
		BlockH:         30,
		BlockRows:      5,
		BlockCols:      10,
		BlockSpacing:   5,
		PaddleWidth:    100,
		PaddleHeight:   20,
		PaddleY:        550,
		PaddleSpeed:    5,
		BallRadius:     10,
		BallSpeed:      5,
		BlockCount:     10,
		MinPaddleGap:   180,
		MaxAttempts:    200,
		MaxBalls:       8,
		ItemDropChance: 0.1,
		MaxItems:       3,
		ItemWidth:      16,
		ItemHeight:     12,
		ItemFallSpeed:  3,
		Seed:           nil,
	}
}

func TestGenerateBlocksRespectsBounds(t *testing.T) {
	cfg := baseLayout()
	cfg.BlockCount = 2
	cfg.MaxAttempts = 50
	// deterministic placement spread across area (x,y pairs)
	rnd := &mockRandom{floats: []float64{0.05, 0.05, 0.6, 0.6}}

	blocks, err := GenerateBlocks(cfg, rnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != cfg.BlockCount {
		t.Fatalf("expected %d blocks, got %d", cfg.BlockCount, len(blocks))
	}
	for i, b := range blocks {
		if b.X < 0 || b.X+cfg.BlockW > cfg.ScreenW {
			t.Fatalf("block %d X out of bounds: %f", i, b.X)
		}
		if b.Y < 0 || b.Y+cfg.BlockH > cfg.PaddleY-cfg.MinPaddleGap {
			t.Fatalf("block %d Y out of bounds: %f", i, b.Y)
		}
	}
}

func TestGenerateBlocksInvalidConfig(t *testing.T) {
	cfg := baseLayout()
	cfg.BlockCount = 0
	if _, err := GenerateBlocks(cfg, nil); err == nil {
		t.Fatalf("expected error for zero block count")
	}

	cfg = baseLayout()
	cfg.MaxAttempts = cfg.BlockCount
	if _, err := GenerateBlocks(cfg, nil); err == nil {
		t.Fatalf("expected error for insufficient attempts")
	}

	cfg = baseLayout()
	cfg.MinPaddleGap = -1
	if _, err := GenerateBlocks(cfg, nil); err == nil {
		t.Fatalf("expected error for negative gap")
	}
}

func TestGenerateBlocksInsufficientSpace(t *testing.T) {
	cfg := baseLayout()
	cfg.BlockW = 400
	cfg.BlockH = 400
	cfg.BlockCount = 5
	cfg.MaxAttempts = 20

	if _, err := GenerateBlocks(cfg, nil); err == nil {
		t.Fatalf("expected error due to insufficient space for blocks")
	}
}

func TestGenerateGridFallback(t *testing.T) {
	cfg := baseLayout()
	blocks := GenerateGridFallback(cfg)
	expected := cfg.BlockRows * cfg.BlockCols
	if len(blocks) != expected {
		t.Fatalf("expected %d blocks, got %d", expected, len(blocks))
	}
	// first block should be near left/top inside screen
	if blocks[0].X < 0 || blocks[0].Y < 0 {
		t.Fatalf("fallback block has negative position: (%f,%f)", blocks[0].X, blocks[0].Y)
	}
}

func TestGenerateBlocksHardLayoutIntegration(t *testing.T) {
	base := baseLayout()
	base.BlockCount = 50
	base.MaxAttempts = 500
	// Apply HARD-like scales: smaller blocks, more count, faster ball
	setting := DifficultySetting{
		Name:             DifficultyHard,
		BallSpeedScale:   1.2,
		BallRadiusScale:  0.9,
		PaddleWidthScale: 0.9,
		PaddleSpeedScale: 0.9,
		BlockSizeScale:   0.9,
		BlockCountScale:  1.3,
	}
	derived, err := ApplyDifficulty(base, setting)
	if err != nil {
		t.Fatalf("apply difficulty failed: %v", err)
	}

	blocks, err := GenerateBlocks(derived, nil)
	if err != nil {
		t.Fatalf("generate blocks failed for hard layout: %v", err)
	}
	if len(blocks) != derived.BlockCount {
		t.Fatalf("expected %d blocks, got %d", derived.BlockCount, len(blocks))
	}
}
