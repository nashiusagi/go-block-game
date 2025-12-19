package domain

import "testing"

func TestBallService_BounceAndClampWall(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})
	svc := NewBallService()

	// 左壁へ向かうよう設定し、反射と位置補正を検証
	b := &state.Balls[0]
	b.X = b.Radius - 1
	b.Y = cfg.ScreenH / 2
	b.VX = -cfg.BallSpeed
	b.VY = 0

	svc.Advance(state, cfg, NewRandomSource(nil))

	if state.Balls[0].VX <= 0 {
		t.Fatalf("expected VX to invert on wall bounce")
	}
	if state.Balls[0].X < state.Balls[0].Radius {
		t.Fatalf("expected X to be clamped inside wall, got %f", state.Balls[0].X)
	}
}

func TestBallService_DestroysBlock(t *testing.T) {
	cfg := baseLayout()
	block := Block{X: 100, Y: 100, Alive: true}
	state := NewGameState(cfg, []Block{block})
	svc := NewBallService()

	b := &state.Balls[0]
	b.X = block.X + cfg.BlockW/2
	b.Y = block.Y - b.Radius - 1
	b.VX = 0
	b.VY = cfg.BallSpeed

	svc.Advance(state, cfg, NewRandomSource(nil))

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

func TestBallService_RemovesFallenBall(t *testing.T) {
	cfg := baseLayout()
	state := NewGameState(cfg, []Block{})
	svc := NewBallService()

	b := &state.Balls[0]
	b.Y = cfg.ScreenH + b.Radius + 1
	b.VY = cfg.BallSpeed

	svc.Advance(state, cfg, NewRandomSource(nil))

	if len(state.Balls) != 0 {
		t.Fatalf("expected fallen ball to be removed, got %d balls", len(state.Balls))
	}
}
