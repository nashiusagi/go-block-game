package domain

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type LayoutConfig struct {
	ScreenW, ScreenH          float64
	BlockW, BlockH            float64
	BlockRows, BlockCols      int
	BlockSpacing              float64
	PaddleWidth, PaddleHeight float64
	PaddleY                   float64
	PaddleSpeed               float64
	BallRadius                float64
	BallSpeed                 float64
	BlockCount                int
	MinPaddleGap              float64
	MaxAttempts               int
	Seed                      *int64
}

type RandomSource interface {
	Float64() float64
	Intn(n int) int
	Seed(seed int64)
}

type defaultRandomSource struct {
	r *rand.Rand
}

func (d *defaultRandomSource) Float64() float64 {
	return d.r.Float64()
}

func (d *defaultRandomSource) Intn(n int) int {
	return d.r.Intn(n)
}

func (d *defaultRandomSource) Seed(seed int64) {
	d.r.Seed(seed)
}

func NewRandomSource(seed *int64) RandomSource {
	seedVal := time.Now().UnixNano()
	if seed != nil {
		seedVal = *seed
	}
	return &defaultRandomSource{
		r: rand.New(rand.NewSource(seedVal)),
	}
}

type Block struct {
	X, Y  float64
	Alive bool
}

func rectsOverlap(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}

func GenerateBlocks(cfg LayoutConfig, rnd RandomSource) ([]Block, error) {
	if cfg.BlockCount <= 0 {
		return nil, errors.New("block count must be positive")
	}
	if cfg.MaxAttempts <= cfg.BlockCount {
		return nil, errors.New("max attempts must exceed block count")
	}
	if cfg.MinPaddleGap <= 0 {
		return nil, errors.New("min paddle gap must be positive")
	}
	if rnd == nil {
		rnd = NewRandomSource(cfg.Seed)
	}

	maxX := cfg.ScreenW - cfg.BlockW
	maxY := cfg.PaddleY - cfg.MinPaddleGap - cfg.BlockH
	if maxX <= 0 || maxY <= 0 {
		return nil, errors.New("invalid layout bounds")
	}

	blocks := make([]Block, 0, cfg.BlockCount)
	attempts := 0

	for len(blocks) < cfg.BlockCount {
		if attempts >= cfg.MaxAttempts {
			return nil, errors.New("block placement exceeded max attempts")
		}
		attempts++

		x := rnd.Float64() * maxX
		y := rnd.Float64() * maxY

		overlap := false
		for _, b := range blocks {
			if rectsOverlap(x, y, cfg.BlockW, cfg.BlockH, b.X, b.Y, cfg.BlockW, cfg.BlockH) {
				overlap = true
				break
			}
		}
		if overlap {
			continue
		}

		blocks = append(blocks, Block{
			X:     x,
			Y:     y,
			Alive: true,
		})
	}

	return blocks, nil
}

func GenerateGridFallback(cfg LayoutConfig) []Block {
	total := cfg.BlockRows * cfg.BlockCols
	if total <= 0 {
		return []Block{}
	}
	blocks := make([]Block, total)
	startX := (cfg.ScreenW - (float64(cfg.BlockCols)*(cfg.BlockW+cfg.BlockSpacing) - cfg.BlockSpacing)) / 2
	startY := 50.0
	for row := 0; row < cfg.BlockRows; row++ {
		for col := 0; col < cfg.BlockCols; col++ {
			idx := row*cfg.BlockCols + col
			blocks[idx] = Block{
				X:     startX + float64(col)*(cfg.BlockW+cfg.BlockSpacing),
				Y:     startY + float64(row)*(cfg.BlockH+cfg.BlockSpacing),
				Alive: true,
			}
		}
	}
	return blocks
}

func reflectVelocity(vx, vy float64) float64 {
	return math.Sqrt(vx*vx + vy*vy)
}
