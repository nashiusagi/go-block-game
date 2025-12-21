package domain

import (
	"testing"
)

func TestDefaultDifficultyProfileHasEntries(t *testing.T) {
	profile := DefaultDifficultyProfile()
	if len(profile.Settings) != 3 {
		t.Fatalf("expected 3 difficulty settings, got %d", len(profile.Settings))
	}
	if profile.Settings[DifficultyNormal].BallSpeedScale != 1.0 {
		t.Fatalf("normal ball speed scale mismatch")
	}
}

func TestValidateKnownDifficulty(t *testing.T) {
	profile := DefaultDifficultyProfile()
	validator := NewDifficultyValidator(DifficultyNormal)

	setting, applied, err := validator.Validate(profile, DifficultyHard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if applied != DifficultyHard {
		t.Fatalf("expected applied %s, got %s", DifficultyHard, applied)
	}
	if setting.BallSpeedScale != profile.Settings[DifficultyHard].BallSpeedScale {
		t.Fatalf("scale not applied")
	}
}

func TestValidateUnknownDifficultyFallsBack(t *testing.T) {
	profile := DefaultDifficultyProfile()
	validator := NewDifficultyValidator(DifficultyNormal)

	_, applied, err := validator.Validate(profile, Difficulty("UNKNOWN"))
	if err == nil {
		t.Fatalf("expected error for unknown difficulty")
	}
	if applied != DifficultyNormal {
		t.Fatalf("expected fallback to %s, got %s", DifficultyNormal, applied)
	}
}

func TestValidateInvalidScaleFallsBack(t *testing.T) {
	profile := DefaultDifficultyProfile()
	invalid := profile.Settings[DifficultyHard]
	invalid.BlockCountScale = 0
	profile.Settings[DifficultyHard] = invalid

	validator := NewDifficultyValidator(DifficultyNormal)
	setting, applied, err := validator.Validate(profile, DifficultyHard)
	if err == nil {
		t.Fatalf("expected validation error for invalid scale")
	}
	if applied != DifficultyNormal {
		t.Fatalf("expected fallback to %s, got %s", DifficultyNormal, applied)
	}
	if setting.Name != DifficultyNormal {
		t.Fatalf("expected fallback setting name %s, got %s", DifficultyNormal, setting.Name)
	}
}

func TestClampAboveMaxScale(t *testing.T) {
	profile := DefaultDifficultyProfile()
	custom := profile.Settings[DifficultyHard]
	custom.BallSpeedScale = 20.0
	profile.Settings[DifficultyHard] = custom

	validator := NewDifficultyValidator(DifficultyNormal)
	setting, applied, err := validator.Validate(profile, DifficultyHard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if applied != DifficultyHard {
		t.Fatalf("expected applied %s, got %s", DifficultyHard, applied)
	}
	if setting.BallSpeedScale != maxScale {
		t.Fatalf("expected clamped scale %f, got %f", maxScale, setting.BallSpeedScale)
	}
}
