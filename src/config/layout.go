package config

import "block-game/src/domain"

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
)

func DefaultLayoutConfig() domain.LayoutConfig {
	return domain.LayoutConfig{
		ScreenW:      ScreenWidth,
		ScreenH:      ScreenHeight,
		BlockW:       BlockWidth,
		BlockH:       BlockHeight,
		BlockRows:    BlockRows,
		BlockCols:    BlockCols,
		BlockSpacing: BlockSpacing,
		PaddleWidth:  PaddleWidth,
		PaddleHeight: PaddleHeight,
		PaddleY:      PaddleY,
		PaddleSpeed:  PaddleSpeed,
		BallRadius:   BallRadius,
		BallSpeed:    BallSpeed,
		BlockCount:   BlockRows * BlockCols,
		MinPaddleGap: MinPaddleGap,
		MaxAttempts:  MaxAttemptsFactor * BlockRows * BlockCols,
		Seed:         nil,
	}
}
