package domain

import (
	"errors"
	"fmt"
)

// Difficulty represents the game difficulty level.
type Difficulty string

const (
	DifficultyEasy   Difficulty = "EASY"
	DifficultyNormal Difficulty = "NORMAL"
	DifficultyHard   Difficulty = "HARD"

	defaultDifficulty = DifficultyNormal
	maxScale          = 10.0
)

// DifficultySetting holds scaling factors applied to the base LayoutConfig.
// Each scale must be greater than 0; values above maxScale are clamped.
type DifficultySetting struct {
	Name             Difficulty
	BallSpeedScale   float64
	BallRadiusScale  float64
	PaddleWidthScale float64
	PaddleSpeedScale float64
	BlockSizeScale   float64
	BlockCountScale  float64
}

// DifficultyProfile stores the available settings and default selection.
type DifficultyProfile struct {
	Settings map[Difficulty]DifficultySetting
	Default  Difficulty
}

// DefaultDifficultyProfile returns predefined settings for EASY/NORMAL/HARD.
func DefaultDifficultyProfile() DifficultyProfile {
	settings := map[Difficulty]DifficultySetting{
		DifficultyEasy: {
			Name:             DifficultyEasy,
			BallSpeedScale:   0.8,
			BallRadiusScale:  1.0,
			PaddleWidthScale: 1.1,
			PaddleSpeedScale: 1.1,
			BlockSizeScale:   1.0,
			BlockCountScale:  1.0,
		},
		DifficultyNormal: {
			Name:             DifficultyNormal,
			BallSpeedScale:   1.0,
			BallRadiusScale:  1.0,
			PaddleWidthScale: 1.0,
			PaddleSpeedScale: 1.0,
			BlockSizeScale:   1.0,
			BlockCountScale:  1.0,
		},
		DifficultyHard: {
			Name:             DifficultyHard,
			BallSpeedScale:   1.2,
			BallRadiusScale:  0.9,
			PaddleWidthScale: 0.9,
			PaddleSpeedScale: 0.9,
			BlockSizeScale:   0.9,
			BlockCountScale:  1.3,
		},
	}
	return DifficultyProfile{
		Settings: settings,
		Default:  defaultDifficulty,
	}
}

// DifficultyValidator validates difficulty selection and clamps scales.
type DifficultyValidator struct {
	Default Difficulty
}

// NewDifficultyValidator creates a validator with the provided default difficulty.
func NewDifficultyValidator(defaultDiff Difficulty) DifficultyValidator {
	if defaultDiff == "" {
		defaultDiff = defaultDifficulty
	}
	return DifficultyValidator{Default: defaultDiff}
}

// Validate returns a safe DifficultySetting. Unknown/invalid inputs fall back to the default.
// The returned Difficulty reflects the applied setting (default when fallback occurs).
func (v DifficultyValidator) Validate(profile DifficultyProfile, selected Difficulty) (DifficultySetting, Difficulty, error) {
	settings := profile.Settings
	if len(settings) == 0 {
		return DifficultySetting{}, v.DefaultDifficulty(profile), errors.New("no difficulty settings defined")
	}

	setting, ok := settings[selected]
	if !ok {
		fallback := v.DefaultDifficulty(profile)
		return settings[fallback], fallback, fmt.Errorf("unknown difficulty: %s", selected)
	}

	clamped, err := clampSetting(setting)
	if err != nil {
		fallback := v.DefaultDifficulty(profile)
		safe := profile.Settings[fallback]
		clampedFallback, _ := clampSetting(safe)
		return clampedFallback, fallback, err
	}

	return clamped, selected, nil
}

// DefaultDifficulty resolves the effective default difficulty.
func (v DifficultyValidator) DefaultDifficulty(profile DifficultyProfile) Difficulty {
	if profile.Default != "" {
		return profile.Default
	}
	return v.Default
}

func clampSetting(setting DifficultySetting) (DifficultySetting, error) {
	var err error
	setting.BallSpeedScale, err = clampScale(setting.BallSpeedScale)
	if err != nil {
		return DifficultySetting{}, fmt.Errorf("ball speed scale: %w", err)
	}
	setting.BallRadiusScale, err = clampScale(setting.BallRadiusScale)
	if err != nil {
		return DifficultySetting{}, fmt.Errorf("ball radius scale: %w", err)
	}
	setting.PaddleWidthScale, err = clampScale(setting.PaddleWidthScale)
	if err != nil {
		return DifficultySetting{}, fmt.Errorf("paddle width scale: %w", err)
	}
	setting.PaddleSpeedScale, err = clampScale(setting.PaddleSpeedScale)
	if err != nil {
		return DifficultySetting{}, fmt.Errorf("paddle speed scale: %w", err)
	}
	setting.BlockSizeScale, err = clampScale(setting.BlockSizeScale)
	if err != nil {
		return DifficultySetting{}, fmt.Errorf("block size scale: %w", err)
	}
	setting.BlockCountScale, err = clampScale(setting.BlockCountScale)
	if err != nil {
		return DifficultySetting{}, fmt.Errorf("block count scale: %w", err)
	}
	return setting, nil
}

func clampScale(scale float64) (float64, error) {
	if scale <= 0 {
		return 0, errors.New("scale must be positive")
	}
	if scale > maxScale {
		return maxScale, nil
	}
	return scale, nil
}
