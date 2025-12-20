package domain

import "math"

type Ball struct {
	X, Y   float64
	VX, VY float64
	Radius float64
}

// BallService はボールの移動・衝突を扱うドメインサービス
type BallService struct{}

func NewBallService() BallService {
	return BallService{}
}

// Advance は全ボールの移動と衝突処理を行う
func (BallService) Advance(state *GameState, cfg LayoutConfig, rnd RandomSource) {
	if rnd == nil {
		rnd = NewRandomSource(cfg.Seed)
	}

	newBalls := make([]Ball, 0, len(state.Balls))
	for _, ball := range state.Balls {
		ball.X += ball.VX
		ball.Y += ball.VY

		if ball.X-ball.Radius <= 0 || ball.X+ball.Radius >= cfg.ScreenW {
			ball.VX = -ball.VX
			// 壁の内側に押し戻すことで連続反射による滑りを防ぐ
			if ball.X-ball.Radius < 0 {
				ball.X = ball.Radius
			} else if ball.X+ball.Radius > cfg.ScreenW {
				ball.X = cfg.ScreenW - ball.Radius
			}
		}

		if ball.Y-ball.Radius <= 0 {
			ball.VY = -ball.VY
			// 上壁との衝突後に位置を補正
			if ball.Y-ball.Radius < 0 {
				ball.Y = ball.Radius
			}
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
