package domain

import "testing"

func TestAdvancePaddleMovement(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})

	leftInput := InputState{MoveLeft: true}
	Advance(state, leftInput, cfg, NewRandomSource(nil))
	afterLeft := state.Paddle.X
	if afterLeft >= (cfg.ScreenW-cfg.PaddleWidth)/2 {
		t.Fatalf("expected paddle to move left")
	}

	rightInput := InputState{MoveRight: true}
	Advance(state, rightInput, cfg, NewRandomSource(nil))
	if state.Paddle.X <= afterLeft {
		t.Fatalf("expected paddle to move right")
	}
}

func TestAdvanceBallBouncesOnWall(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})
	state.Balls[0].X = state.Balls[0].Radius - 1
	state.Balls[0].VX = -cfg.BallSpeed

	Advance(state, InputState{}, cfg, NewRandomSource(nil))
	if state.Balls[0].VX <= 0 {
		t.Fatalf("expected VX to invert on wall bounce")
	}
}

func TestAdvanceDestroysBlock(t *testing.T) {
	cfg := baseLayout()
	block := Block{X: 100, Y: 100, Alive: true}
	state := NewGameState(cfg, []Block{block})

	state.Balls[0].X = block.X + cfg.BlockW/2
	state.Balls[0].Y = block.Y - state.Balls[0].Radius - 1
	state.Balls[0].VX = 0
	state.Balls[0].VY = cfg.BallSpeed

	Advance(state, InputState{}, cfg, NewRandomSource(nil))

	if state.Score != 1 {
		t.Fatalf("expected score to increment, got %d", state.Score)
	}
	if state.Blocks[0].Alive {
		t.Fatalf("expected block to be destroyed")
	}
	if state.Balls[0].VY >= 0 {
		t.Fatalf("expected VY to invert after collision")
	}
}

func TestAdvanceGameOverWhenBallFalls(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})
	state.Balls[0].Y = cfg.ScreenH + state.Balls[0].Radius + 1

	Advance(state, InputState{}, cfg, NewRandomSource(nil))
	if !state.GameOver {
		t.Fatalf("expected game over when ball falls below screen")
	}
}

func TestItemDropAndPickupTriggersMultiball(t *testing.T) {
	cfg := baseLayout()
	cfg.ItemDropChance = 1.0
	state := NewGameState(cfg, []Block{
		{X: 100, Y: 100, Alive: true},
		{X: 200, Y: 100, Alive: true},
	})

	// position ball to hit the block
	state.Balls[0].X = 100 + cfg.BlockW/2
	state.Balls[0].Y = 100 - state.Balls[0].Radius - 1
	state.Balls[0].VX = 0
	state.Balls[0].VY = cfg.BallSpeed

	rnd := NewRandomSource(nil)
	Advance(state, InputState{}, cfg, rnd)

	if len(state.Items) == 0 {
		t.Fatalf("expected an item to drop")
	}

	// place item on paddle to trigger pickup
	state.Items[0].Y = state.Paddle.Y
	state.Items[0].X = state.Paddle.X
	Advance(state, InputState{}, cfg, rnd)

	if len(state.Balls) < 2 {
		t.Fatalf("expected balls to increase after pickup, got %d", len(state.Balls))
	}
}
