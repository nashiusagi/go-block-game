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

type GameState struct {
	Blocks   []Block
	Ball     Ball
	Paddle   Paddle
	Score    int
	GameOver bool
}

type InputState struct {
	MoveLeft  bool
	MoveRight bool
}

func NewGameState(cfg LayoutConfig, blocks []Block) *GameState {
	return &GameState{
		Blocks: blocks,
		Ball: Ball{
			X:      cfg.ScreenW / 2,
			Y:      cfg.ScreenH / 2,
			Radius: cfg.BallRadius,
			VX:     cfg.BallSpeed * math.Cos(math.Pi/4),
			VY:     -cfg.BallSpeed * math.Sin(math.Pi/4),
		},
		Paddle: Paddle{
			X:      (cfg.ScreenW - cfg.PaddleWidth) / 2,
			Y:      cfg.PaddleY,
			Width:  cfg.PaddleWidth,
			Height: cfg.PaddleHeight,
		},
		Score:    0,
		GameOver: false,
	}
}

func Advance(state *GameState, input InputState, cfg LayoutConfig) {
	if state.GameOver {
		return
	}

	if input.MoveLeft && state.Paddle.X > 0 {
		state.Paddle.X -= cfg.PaddleSpeed
	}
	if input.MoveRight && state.Paddle.X < cfg.ScreenW-state.Paddle.Width {
		state.Paddle.X += cfg.PaddleSpeed
	}

	state.Ball.X += state.Ball.VX
	state.Ball.Y += state.Ball.VY

	if state.Ball.X-state.Ball.Radius <= 0 || state.Ball.X+state.Ball.Radius >= cfg.ScreenW {
		state.Ball.VX = -state.Ball.VX
	}

	if state.Ball.Y-state.Ball.Radius <= 0 {
		state.Ball.VY = -state.Ball.VY
	}

	if state.Ball.Y+state.Ball.Radius > cfg.ScreenH {
		state.GameOver = true
		return
	}

	if state.Ball.Y+state.Ball.Radius >= state.Paddle.Y &&
		state.Ball.Y-state.Ball.Radius <= state.Paddle.Y+state.Paddle.Height &&
		state.Ball.X+state.Ball.Radius >= state.Paddle.X &&
		state.Ball.X-state.Ball.Radius <= state.Paddle.X+state.Paddle.Width {
		hitPos := (state.Ball.X - state.Paddle.X) / state.Paddle.Width
		angle := math.Pi * (0.5 + hitPos*0.5)
		speed := reflectVelocity(state.Ball.VX, state.Ball.VY)
		state.Ball.VX = speed * math.Cos(angle)
		state.Ball.VY = -speed * math.Sin(angle)
		state.Ball.Y = state.Paddle.Y - state.Ball.Radius
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

		dx := state.Ball.X - blockCenterX
		dy := state.Ball.Y - blockCenterY

		if math.Abs(dx) < blockHalfWidth+state.Ball.Radius &&
			math.Abs(dy) < blockHalfHeight+state.Ball.Radius {
			block.Alive = false
			state.Score++

			if math.Abs(dx/blockHalfWidth) > math.Abs(dy/blockHalfHeight) {
				state.Ball.VX = -state.Ball.VX
			} else {
				state.Ball.VY = -state.Ball.VY
			}

			if dx > 0 {
				state.Ball.X = block.X + cfg.BlockW + state.Ball.Radius
			} else {
				state.Ball.X = block.X - state.Ball.Radius
			}
			if dy > 0 {
				state.Ball.Y = block.Y + cfg.BlockH + state.Ball.Radius
			} else {
				state.Ball.Y = block.Y - state.Ball.Radius
			}
		}
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
