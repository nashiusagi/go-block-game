package domain

import (
	"math"
	"testing"
)

func baseLayoutForDifficulty() LayoutConfig {
	return LayoutConfig{
		ScreenW:      800,
		ScreenH:      600,
		BlockW:       70,
		BlockH:       30,
		BlockRows:    5,
		BlockCols:    10,
		BlockSpacing: 5,
		PaddleWidth:  100,
		PaddleHeight: 20,
		PaddleY:      550,
		PaddleSpeed:  5,
		BallRadius:   10,
		BallSpeed:    5,
		BlockCount:   50,
		MinPaddleGap: 180,
		MaxAttempts:  500,
		MaxBalls:     8,
	}
}

func TestApplyDifficultyScalesValues(t *testing.T) {
	base := baseLayoutForDifficulty()
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
		t.Fatalf("unexpected error: %v", err)
	}

	if derived.Difficulty != DifficultyHard {
		t.Fatalf("expected difficulty %s, got %s", DifficultyHard, derived.Difficulty)
	}

	if derived.BallSpeed != base.BallSpeed*setting.BallSpeedScale {
		t.Fatalf("ball speed not scaled")
	}
	if derived.BallRadius != base.BallRadius*setting.BallRadiusScale {
		t.Fatalf("ball radius not scaled")
	}
	if derived.PaddleWidth != base.PaddleWidth*setting.PaddleWidthScale {
		t.Fatalf("paddle width not scaled")
	}
	if derived.PaddleSpeed != base.PaddleSpeed*setting.PaddleSpeedScale {
		t.Fatalf("paddle speed not scaled")
	}
	if derived.BlockW != base.BlockW*setting.BlockSizeScale || derived.BlockH != base.BlockH*setting.BlockSizeScale {
		t.Fatalf("block size not scaled")
	}

	expectedCount := int(math.Round(float64(base.BlockCount) * setting.BlockCountScale))
	if derived.BlockCount != expectedCount {
		t.Fatalf("block count scaled incorrectly: want %d, got %d", expectedCount, derived.BlockCount)
	}

	perBlock := float64(base.MaxAttempts) / float64(base.BlockCount)
	expectedAttempts := int(math.Ceil(perBlock * float64(expectedCount)))
	if expectedAttempts <= expectedCount {
		expectedAttempts = expectedCount + 1
	}
	if derived.MaxAttempts != expectedAttempts {
		t.Fatalf("max attempts not scaled: want %d, got %d", expectedAttempts, derived.MaxAttempts)
	}
}

func TestApplyDifficultyRejectsNonPositiveScale(t *testing.T) {
	base := baseLayoutForDifficulty()
	setting := DifficultySetting{
		Name:             DifficultyEasy,
		BallSpeedScale:   0,
		BallRadiusScale:  1,
		PaddleWidthScale: 1,
		PaddleSpeedScale: 1,
		BlockSizeScale:   1,
		BlockCountScale:  1,
	}

	if _, err := ApplyDifficulty(base, setting); err == nil {
		t.Fatalf("expected error for non-positive scale")
	}
}
