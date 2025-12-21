package domain

import (
	"errors"
	"math"
)

// ApplyDifficulty applies the given difficulty setting to a base LayoutConfig and
// returns a derived configuration. Validation of the setting should be done
// beforehand; this function assumes positive scales.
func ApplyDifficulty(base LayoutConfig, setting DifficultySetting) (LayoutConfig, error) {
	if setting.BallSpeedScale <= 0 ||
		settingsScaleNonPositive(setting) {
		return LayoutConfig{}, errors.New("difficulty scales must be positive")
	}

	derived := base
	derived.Difficulty = setting.Name

	// Scale core dimensions and speeds.
	derived.BallSpeed = base.BallSpeed * setting.BallSpeedScale
	derived.BallRadius = base.BallRadius * setting.BallRadiusScale
	derived.PaddleWidth = base.PaddleWidth * setting.PaddleWidthScale
	derived.PaddleSpeed = base.PaddleSpeed * setting.PaddleSpeedScale
	derived.BlockW = base.BlockW * setting.BlockSizeScale
	derived.BlockH = base.BlockH * setting.BlockSizeScale

	// Scale block count with rounding and enforce minimum of 1.
	scaledCount := int(math.Round(float64(base.BlockCount) * setting.BlockCountScale))
	if scaledCount < 1 {
		scaledCount = 1
	}
	derived.BlockCount = scaledCount

	// Scale MaxAttempts proportionally to block count, preserving the original per-block factor.
	perBlockAttempts := float64(base.MaxAttempts) / float64(base.BlockCount)
	derived.MaxAttempts = int(math.Ceil(perBlockAttempts * float64(derived.BlockCount)))
	if derived.MaxAttempts <= derived.BlockCount {
		derived.MaxAttempts = derived.BlockCount + 1
	}

	// Keep rows/cols as a reference layout; generation uses BlockCount.
	derived.BlockRows = base.BlockRows
	derived.BlockCols = base.BlockCols

	// Basic guard for non-positive dimensions.
	if derived.BlockW <= 0 || derived.BlockH <= 0 || derived.BallRadius <= 0 || derived.PaddleWidth <= 0 {
		return LayoutConfig{}, errors.New("scaled dimensions must be positive")
	}

	return derived, nil
}

func settingsScaleNonPositive(setting DifficultySetting) bool {
	return setting.BallRadiusScale <= 0 ||
		setting.PaddleWidthScale <= 0 ||
		setting.PaddleSpeedScale <= 0 ||
		setting.BlockSizeScale <= 0 ||
		setting.BlockCountScale <= 0
}
