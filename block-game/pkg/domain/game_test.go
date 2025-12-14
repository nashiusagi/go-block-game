package domain

import "testing"

func TestAdvancePaddleMovement(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})

	leftInput := InputState{MoveLeft: true}
	Advance(state, leftInput, cfg)
	afterLeft := state.Paddle.X
	if afterLeft >= (cfg.ScreenW-cfg.PaddleWidth)/2 {
		t.Fatalf("expected paddle to move left")
	}

	rightInput := InputState{MoveRight: true}
	Advance(state, rightInput, cfg)
	if state.Paddle.X <= afterLeft {
		t.Fatalf("expected paddle to move right")
	}
}

func TestAdvanceBallBouncesOnWall(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})
	state.Ball.X = state.Ball.Radius - 1
	state.Ball.VX = -cfg.BallSpeed

	Advance(state, InputState{}, cfg)
	if state.Ball.VX <= 0 {
		t.Fatalf("expected VX to invert on wall bounce")
	}
}

func TestAdvanceDestroysBlock(t *testing.T) {
	cfg := baseLayout()
	block := Block{X: 100, Y: 100, Alive: true}
	state := NewGameState(cfg, []Block{block})

	state.Ball.X = block.X + cfg.BlockW/2
	state.Ball.Y = block.Y - state.Ball.Radius - 1
	state.Ball.VX = 0
	state.Ball.VY = cfg.BallSpeed

	Advance(state, InputState{}, cfg)

	if state.Score != 1 {
		t.Fatalf("expected score to increment, got %d", state.Score)
	}
	if state.Blocks[0].Alive {
		t.Fatalf("expected block to be destroyed")
	}
	if state.Ball.VY >= 0 {
		t.Fatalf("expected VY to invert after collision")
	}
}

func TestAdvanceGameOverWhenBallFalls(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})
	state.Ball.Y = cfg.ScreenH + state.Ball.Radius + 1

	Advance(state, InputState{}, cfg)
	if !state.GameOver {
		t.Fatalf("expected game over when ball falls below screen")
	}
}
