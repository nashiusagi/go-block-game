package domain

import "math"

type Ball struct {
	X, Y   float64
	VX, VY float64
	Radius float64
}

type Paddle struct {
	X, Y   float64
	Width  float64
	Height float64
}

type Item struct {
	X, Y   float64
	Width  float64
	Height float64
	VY     float64
	Active bool
}

type GameState struct {
	Blocks   []Block
	Balls    []Ball
	Paddle   Paddle
	Items    []Item
	Score    int
	GameOver bool
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

	updateBalls(state, cfg, rnd)

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
			applyMultiball(state, cfg)
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

func updateBalls(state *GameState, cfg LayoutConfig, rnd RandomSource) {
	if rnd == nil {
		rnd = NewRandomSource(cfg.Seed)
	}

	newBalls := make([]Ball, 0, len(state.Balls))
	for _, ball := range state.Balls {
		ball.X += ball.VX
		ball.Y += ball.VY

		if ball.X-ball.Radius <= 0 || ball.X+ball.Radius >= cfg.ScreenW {
			ball.VX = -ball.VX
		}

		if ball.Y-ball.Radius <= 0 {
			ball.VY = -ball.VY
		}

		if ball.Y+ball.Radius > cfg.ScreenH {
			continue
		}

		if ball.Y+ball.Radius >= state.Paddle.Y &&
			ball.Y-ball.Radius <= state.Paddle.Y+state.Paddle.Height &&
			ball.X+ball.Radius >= state.Paddle.X &&
			ball.X-ball.Radius <= state.Paddle.X+state.Paddle.Width {
			hitPos := (ball.X - state.Paddle.X) / state.Paddle.Width
			angle := math.Pi * (0.5 + hitPos*0.5)
			speed := reflectVelocity(ball.VX, ball.VY)
			ball.VX = speed * math.Cos(angle)
			ball.VY = -speed * math.Sin(angle)
			ball.Y = state.Paddle.Y - ball.Radius
		}

		for i := range state.Blocks {
			if !state.Blocks[i].Alive {
				continue
			}

			block := &state.Blocks[i]
			blockCenterX := block.X + cfg.BlockW/2
			blockCenterY := block.Y + cfg.BlockH/2
			blockHalfWidth := cfg.BlockW / 2
			blockHalfHeight := cfg.BlockH / 2

			dx := ball.X - blockCenterX
			dy := ball.Y - blockCenterY

			if math.Abs(dx) < blockHalfWidth+ball.Radius &&
				math.Abs(dy) < blockHalfHeight+ball.Radius {
				block.Alive = false
				state.Score++
				tryDropItem(state, cfg, block, rnd)

				if math.Abs(dx/blockHalfWidth) > math.Abs(dy/blockHalfHeight) {
					ball.VX = -ball.VX
				} else {
					ball.VY = -ball.VY
				}

				if dx > 0 {
					ball.X = block.X + cfg.BlockW + ball.Radius
				} else {
					ball.X = block.X - ball.Radius
				}
				if dy > 0 {
					ball.Y = block.Y + cfg.BlockH + ball.Radius
				} else {
					ball.Y = block.Y - ball.Radius
				}
			}
		}

		newBalls = append(newBalls, ball)
	}

	state.Balls = newBalls
}

func tryDropItem(state *GameState, cfg LayoutConfig, block *Block, rnd RandomSource) {
	if len(state.Items) >= cfg.MaxItems {
		return
	}
	if rnd.Float64() >= cfg.ItemDropChance {
		return
	}
	state.Items = append(state.Items, Item{
		X:      block.X + cfg.BlockW/2 - cfg.ItemWidth/2,
		Y:      block.Y + cfg.BlockH/2 - cfg.ItemHeight/2,
		Width:  cfg.ItemWidth,
		Height: cfg.ItemHeight,
		VY:     cfg.ItemFallSpeed,
		Active: true,
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
