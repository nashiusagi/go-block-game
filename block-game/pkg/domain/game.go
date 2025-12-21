package domain

import "math"

var ballService = NewBallService()

type Paddle struct {
	X, Y   float64
	Width  float64
	Height float64
}

// ItemType represents the type of a falling item.
type ItemType int

const (
	ItemTypeMultiball ItemType = iota
	ItemTypePaddleEnlarge
)

type Item struct {
	X, Y   float64
	Width  float64
	Height float64
	VY     float64
	Active bool
	Type   ItemType
}

// PaddleEffect tracks the temporary paddle enlargement state.
type PaddleEffect struct {
	Active         bool
	RemainingTicks int     // 60FPS: 5 seconds = 300 ticks
	BaseWidth      float64 // original paddle width before effect
	Multiplier     float64 // e.g., 3.0
}

type GameState struct {
	Blocks       []Block
	Balls        []Ball
	Paddle       Paddle
	Items        []Item
	PaddleEffect PaddleEffect
	Score        int
	GameOver     bool
}

type InputState struct {
	MoveLeft  bool
	MoveRight bool
}

func NewGameState(cfg LayoutConfig, blocks []Block) *GameState {
	initialBall := Ball{
		X:      cfg.ScreenW / 2,
		Y:      cfg.ScreenH / 2,
		Radius: cfg.BallRadius,
		VX:     cfg.BallSpeed * math.Cos(math.Pi/4),
		VY:     -cfg.BallSpeed * math.Sin(math.Pi/4),
	}
	return &GameState{
		Blocks: blocks,
		Balls:  []Ball{initialBall},
		Paddle: Paddle{
			X:      (cfg.ScreenW - cfg.PaddleWidth) / 2,
			Y:      cfg.PaddleY,
			Width:  cfg.PaddleWidth,
			Height: cfg.PaddleHeight,
		},
		Items:    []Item{},
		Score:    0,
		GameOver: false,
	}
}

func Advance(state *GameState, input InputState, cfg LayoutConfig, rnd RandomSource) {
	if state.GameOver {
		return
	}

	if input.MoveLeft && state.Paddle.X > 0 {
		state.Paddle.X -= cfg.PaddleSpeed
	}
	if input.MoveRight && state.Paddle.X < cfg.ScreenW-state.Paddle.Width {
		state.Paddle.X += cfg.PaddleSpeed
	}

	updateItems(state, cfg)

	ballService.Advance(state, cfg, rnd)

	if len(state.Balls) == 0 {
		state.GameOver = true
		return
	}

	allBlocksDestroyed := true
	for _, block := range state.Blocks {
		if block.Alive {
			allBlocksDestroyed = false
			break
		}
	}
	if allBlocksDestroyed && len(state.Blocks) > 0 {
		state.GameOver = true
	}
}

func updateItems(state *GameState, cfg LayoutConfig) {
	active := state.Items[:0]
	for _, item := range state.Items {
		if !item.Active {
			continue
		}
		item.Y += item.VY

		if rectsOverlap(item.X, item.Y, item.Width, item.Height, state.Paddle.X, state.Paddle.Y, state.Paddle.Width, state.Paddle.Height) {
			// Apply effect based on item type
			switch item.Type {
			case ItemTypeMultiball:
				applyMultiball(state, cfg)
			case ItemTypePaddleEnlarge:
				applyPaddleEnlarge(state, cfg)
			}
			item.Active = false
		} else if item.Y > cfg.ScreenH {
			item.Active = false
		}

		if item.Active {
			active = append(active, item)
		}
	}
	state.Items = active
}

func tryDropItem(state *GameState, cfg LayoutConfig, block *Block, rnd RandomSource) {
	// Multiball item lottery (independent)
	if len(state.Items) < cfg.MaxItems && rnd.Float64() < cfg.ItemDropChance {
		spawnItem(state, cfg, block, ItemTypeMultiball)
	}
	// Paddle-enlarge item lottery (independent)
	if len(state.Items) < cfg.MaxItems && rnd.Float64() < cfg.PaddleEnlargeChance {
		spawnItem(state, cfg, block, ItemTypePaddleEnlarge)
	}
}

// spawnItem creates a new item of the given type at the block's position.
func spawnItem(state *GameState, cfg LayoutConfig, block *Block, itemType ItemType) {
	state.Items = append(state.Items, Item{
		X:      block.X + cfg.BlockW/2 - cfg.ItemWidth/2,
		Y:      block.Y + cfg.BlockH/2 - cfg.ItemHeight/2,
		Width:  cfg.ItemWidth,
		Height: cfg.ItemHeight,
		VY:     cfg.ItemFallSpeed,
		Active: true,
		Type:   itemType,
	})
}

func applyMultiball(state *GameState, cfg LayoutConfig) {
	if len(state.Balls) == 0 {
		return
	}

	target := len(state.Balls) * 2
	if target > cfg.MaxBalls {
		target = cfg.MaxBalls
	}
	if target <= len(state.Balls) {
		return
	}

	newBalls := make([]Ball, 0, target)
	newBalls = append(newBalls, state.Balls...)

	for _, b := range state.Balls {
		if len(newBalls) >= target {
			break
		}
		dup := b
		dup.VX = -dup.VX
		newBalls = append(newBalls, dup)
	}

	state.Balls = newBalls
}

// applyPaddleEnlarge activates the paddle enlargement effect.
// If already active, it resets the duration timer.
func applyPaddleEnlarge(state *GameState, cfg LayoutConfig) {
	effect := &state.PaddleEffect
	if !effect.Active {
		// First activation: save base width and enlarge
		effect.BaseWidth = state.Paddle.Width
		effect.Multiplier = cfg.PaddleEnlargeMultiplier
		state.Paddle.Width = effect.BaseWidth * effect.Multiplier
	}
	// (Re)set timer
	effect.Active = true
	effect.RemainingTicks = cfg.PaddleEnlargeDuration
}
