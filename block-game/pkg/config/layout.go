package config

import "block-game/pkg/domain"

const (
	ScreenWidth       = 800
	ScreenHeight      = 600
	BlockRows         = 5
	BlockCols         = 10
	BlockWidth        = 70
	BlockHeight       = 30
	BlockSpacing      = 5.0
	PaddleWidth       = 100
	PaddleHeight      = 20
	PaddleY           = ScreenHeight - 50
	PaddleSpeed       = 5.0
	BallRadius        = 10
	BallSpeed         = 5.0
	MinPaddleGap      = 180.0
	MaxAttemptsFactor = 10
	ItemDropChance    = 0.1
	ItemMaxCount      = 3
	ItemFallSpeed     = 3.0
	ItemWidth         = 16.0
	ItemHeight        = 12.0
	MaxBalls          = 8
)

func DefaultLayoutConfig() domain.LayoutConfig {
	return domain.LayoutConfig{
		ScreenW:        ScreenWidth,
		ScreenH:        ScreenHeight,
		BlockW:         BlockWidth,
		BlockH:         BlockHeight,
		BlockRows:      BlockRows,
		BlockCols:      BlockCols,
		BlockSpacing:   BlockSpacing,
		PaddleWidth:    PaddleWidth,
		PaddleHeight:   PaddleHeight,
		PaddleY:        PaddleY,
		PaddleSpeed:    PaddleSpeed,
		BallRadius:     BallRadius,
		BallSpeed:      BallSpeed,
		BlockCount:     BlockRows * BlockCols,
		MinPaddleGap:   MinPaddleGap,
		MaxAttempts:    MaxAttemptsFactor * BlockRows * BlockCols,
		MaxBalls:       MaxBalls,
		ItemDropChance: ItemDropChance,
		MaxItems:       ItemMaxCount,
		ItemWidth:      ItemWidth,
		ItemHeight:     ItemHeight,
		ItemFallSpeed:  ItemFallSpeed,
		Difficulty:     domain.DifficultyNormal,
		Seed:           nil,
	}
}
