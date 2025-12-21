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

// --- Paddle-enlarge item tests ---

func TestTryDropItemPaddleEnlargeIndependent(t *testing.T) {
	cfg := baseLayout()
	cfg.ItemDropChance = 0.0      // no multiball
	cfg.PaddleEnlargeChance = 1.0 // always paddle-enlarge
	block := Block{X: 100, Y: 100, Alive: true}
	state := NewGameState(cfg, []Block{block})

	rnd := NewRandomSource(nil)
	tryDropItem(state, cfg, &block, rnd)

	if len(state.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(state.Items))
	}
	if state.Items[0].Type != ItemTypePaddleEnlarge {
		t.Fatalf("expected paddle-enlarge item, got %v", state.Items[0].Type)
	}
}

func TestTryDropItemBothTypesCanSpawn(t *testing.T) {
	cfg := baseLayout()
	cfg.ItemDropChance = 1.0
	cfg.PaddleEnlargeChance = 1.0
	cfg.MaxItems = 10
	block := Block{X: 100, Y: 100, Alive: true}
	state := NewGameState(cfg, []Block{block})

	rnd := NewRandomSource(nil)
	tryDropItem(state, cfg, &block, rnd)

	if len(state.Items) != 2 {
		t.Fatalf("expected 2 items (both types), got %d", len(state.Items))
	}

	hasMultiball := false
	hasPaddleEnlarge := false
	for _, item := range state.Items {
		if item.Type == ItemTypeMultiball {
			hasMultiball = true
		}
		if item.Type == ItemTypePaddleEnlarge {
			hasPaddleEnlarge = true
		}
	}
	if !hasMultiball || !hasPaddleEnlarge {
		t.Fatalf("expected both item types, multiball=%v paddleEnlarge=%v", hasMultiball, hasPaddleEnlarge)
	}
}

func TestApplyPaddleEnlarge(t *testing.T) {
	cfg := baseLayout()
	cfg.PaddleEnlargeDuration = 300
	cfg.PaddleEnlargeMultiplier = 3.0
	state := NewGameState(cfg, []Block{})
	originalWidth := state.Paddle.Width

	applyPaddleEnlarge(state, cfg)

	if !state.PaddleEffect.Active {
		t.Fatal("expected PaddleEffect.Active to be true")
	}
	if state.PaddleEffect.RemainingTicks != 300 {
		t.Fatalf("expected RemainingTicks=300, got %d", state.PaddleEffect.RemainingTicks)
	}
	if state.PaddleEffect.BaseWidth != originalWidth {
		t.Fatalf("expected BaseWidth=%v, got %v", originalWidth, state.PaddleEffect.BaseWidth)
	}
	if state.Paddle.Width != originalWidth*3.0 {
		t.Fatalf("expected paddle width=%v, got %v", originalWidth*3.0, state.Paddle.Width)
	}
}

func TestApplyPaddleEnlargeRePickupResetsTimer(t *testing.T) {
	cfg := baseLayout()
	cfg.PaddleEnlargeDuration = 300
	cfg.PaddleEnlargeMultiplier = 3.0
	state := NewGameState(cfg, []Block{})

	applyPaddleEnlarge(state, cfg)
	enlargedWidth := state.Paddle.Width

	// Simulate some time passing
	for i := 0; i < 100; i++ {
		updatePaddleEffect(state)
	}
	if state.PaddleEffect.RemainingTicks != 200 {
		t.Fatalf("expected RemainingTicks=200 after 100 ticks, got %d", state.PaddleEffect.RemainingTicks)
	}

	// Re-pickup
	applyPaddleEnlarge(state, cfg)

	if state.PaddleEffect.RemainingTicks != 300 {
		t.Fatalf("expected RemainingTicks reset to 300, got %d", state.PaddleEffect.RemainingTicks)
	}
	if state.Paddle.Width != enlargedWidth {
		t.Fatalf("expected paddle width unchanged at %v, got %v", enlargedWidth, state.Paddle.Width)
	}
}

func TestUpdatePaddleEffectCountdown(t *testing.T) {
	cfg := baseLayout()
	cfg.PaddleEnlargeDuration = 10
	cfg.PaddleEnlargeMultiplier = 3.0
	state := NewGameState(cfg, []Block{})
	originalWidth := state.Paddle.Width

	applyPaddleEnlarge(state, cfg)

	// Count down to 0
	for i := 0; i < 10; i++ {
		if !state.PaddleEffect.Active {
			t.Fatalf("effect should still be active at tick %d", i)
		}
		updatePaddleEffect(state)
	}

	if state.PaddleEffect.Active {
		t.Fatal("expected PaddleEffect.Active to be false after countdown")
	}
	if state.Paddle.Width != originalWidth {
		t.Fatalf("expected paddle width reverted to %v, got %v", originalWidth, state.Paddle.Width)
	}
}

func TestPaddleEnlargeItemPickupFlow(t *testing.T) {
	cfg := baseLayout()
	cfg.ItemDropChance = 0.0
	cfg.PaddleEnlargeChance = 1.0
	cfg.PaddleEnlargeDuration = 5
	cfg.PaddleEnlargeMultiplier = 3.0

	block := Block{X: 100, Y: 100, Alive: true}
	state := NewGameState(cfg, []Block{block})
	originalWidth := state.Paddle.Width

	// Position ball to hit block
	state.Balls[0].X = block.X + cfg.BlockW/2
	state.Balls[0].Y = block.Y - state.Balls[0].Radius - 1
	state.Balls[0].VX = 0
	state.Balls[0].VY = cfg.BallSpeed

	rnd := NewRandomSource(nil)
	Advance(state, InputState{}, cfg, rnd)

	if len(state.Items) == 0 {
		t.Fatal("expected paddle-enlarge item to drop")
	}
	if state.Items[0].Type != ItemTypePaddleEnlarge {
		t.Fatalf("expected paddle-enlarge type, got %v", state.Items[0].Type)
	}

	// Directly call applyPaddleEnlarge to simulate pickup (updateItems handles collision)
	applyPaddleEnlarge(state, cfg)

	if !state.PaddleEffect.Active {
		t.Fatal("expected effect to be active after pickup")
	}
	if state.Paddle.Width != originalWidth*3.0 {
		t.Fatalf("expected width %v, got %v", originalWidth*3.0, state.Paddle.Width)
	}

	// Wait for effect to expire
	for i := 0; i < 10; i++ {
		updatePaddleEffect(state)
	}

	if state.PaddleEffect.Active {
		t.Fatal("expected effect to expire")
	}
	if state.Paddle.Width != originalWidth {
		t.Fatalf("expected width reverted to %v, got %v", originalWidth, state.Paddle.Width)
	}
}
