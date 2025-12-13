package main

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth       = 800
	screenHeight      = 600
	blockRows         = 5
	blockCols         = 10
	blockWidth        = 70
	blockHeight       = 30
	blockSpacing      = 5
	paddleWidth       = 100
	paddleHeight      = 20
	ballRadius        = 10
	ballSpeed         = 5.0
	minPaddleGap      = 180.0
	maxAttemptsFactor = 10
)

type LayoutConfig struct {
	ScreenW, ScreenH float64
	BlockW, BlockH   float64
	BlockCount       int
	MinPaddleGap     float64
	PaddleY          float64
	MaxAttempts      int
	Seed             *int64
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

func newRandomSource(seed *int64) RandomSource {
	seedVal := time.Now().UnixNano()
	if seed != nil {
		seedVal = *seed
	}
	return &defaultRandomSource{
		r: rand.New(rand.NewSource(seedVal)),
	}
}

func newLayoutConfig(paddleY float64) LayoutConfig {
	return LayoutConfig{
		ScreenW:      screenWidth,
		ScreenH:      screenHeight,
		BlockW:       blockWidth,
		BlockH:       blockHeight,
		BlockCount:   blockRows * blockCols,
		MinPaddleGap: minPaddleGap,
		PaddleY:      paddleY,
		MaxAttempts:  maxAttemptsFactor * blockRows * blockCols,
		Seed:         nil,
	}
}

func rectsOverlap(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}

type Block struct {
	X, Y  float64
	Alive bool
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
		rnd = newRandomSource(cfg.Seed)
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

type Game struct {
	blocks   []Block
	ball     Ball
	paddle   Paddle
	score    int
	gameOver bool
}

func NewGame() *Game {
	g := &Game{}
	g.initBall()
	g.initPaddle()
	g.initBlocks()
	return g
}

func (g *Game) initBlocks() {
	cfg := newLayoutConfig(g.paddle.Y)
	blocks, err := GenerateBlocks(cfg, newRandomSource(cfg.Seed))
	if err != nil {
		// フォールバック: 固定グリッド配置に切り替える
		g.blocks = make([]Block, blockRows*blockCols)
		startX := (screenWidth - (blockCols*(blockWidth+blockSpacing) - blockSpacing)) / 2
		startY := 50
		for row := 0; row < blockRows; row++ {
			for col := 0; col < blockCols; col++ {
				idx := row*blockCols + col
				g.blocks[idx] = Block{
					X:     float64(startX + col*(blockWidth+blockSpacing)),
					Y:     float64(startY + row*(blockHeight+blockSpacing)),
					Alive: true,
				}
			}
		}
		return
	}
	g.blocks = blocks
}

func (g *Game) initBall() {
	g.ball = Ball{
		X:      float64(screenWidth) / 2,
		Y:      float64(screenHeight) / 2,
		Radius: ballRadius,
		VX:     ballSpeed * math.Cos(math.Pi/4),
		VY:     -ballSpeed * math.Sin(math.Pi/4),
	}
}

func (g *Game) initPaddle() {
	g.paddle = Paddle{
		X:      float64(screenWidth-paddleWidth) / 2,
		Y:      float64(screenHeight - 50),
		Width:  paddleWidth,
		Height: paddleHeight,
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	// パドルの移動
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.paddle.X > 0 {
		g.paddle.X -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.paddle.X < float64(screenWidth)-g.paddle.Width {
		g.paddle.X += 5
	}

	// ボールの移動
	g.ball.X += g.ball.VX
	g.ball.Y += g.ball.VY

	// 壁との衝突判定（左右）
	if g.ball.X-g.ball.Radius <= 0 || g.ball.X+g.ball.Radius >= float64(screenWidth) {
		g.ball.VX = -g.ball.VX
	}

	// 天井との衝突判定
	if g.ball.Y-g.ball.Radius <= 0 {
		g.ball.VY = -g.ball.VY
	}

	// ボールがフィールドから落ちた場合
	if g.ball.Y+g.ball.Radius > float64(screenHeight) {
		g.gameOver = true
		return nil
	}

	// パドルとの衝突判定
	if g.ball.Y+g.ball.Radius >= g.paddle.Y &&
		g.ball.Y-g.ball.Radius <= g.paddle.Y+g.paddle.Height &&
		g.ball.X+g.ball.Radius >= g.paddle.X &&
		g.ball.X-g.ball.Radius <= g.paddle.X+g.paddle.Width {
		// パドルのどの位置に当たったかで反射角度を変える
		hitPos := (g.ball.X - g.paddle.X) / g.paddle.Width
		angle := math.Pi * (0.5 + hitPos*0.5) // 45度から135度の範囲
		speed := math.Sqrt(g.ball.VX*g.ball.VX + g.ball.VY*g.ball.VY)
		g.ball.VX = speed * math.Cos(angle)
		g.ball.VY = -speed * math.Sin(angle)
		g.ball.Y = g.paddle.Y - g.ball.Radius
	}

	// ブロックとの衝突判定
	for i := range g.blocks {
		if !g.blocks[i].Alive {
			continue
		}

		block := &g.blocks[i]
		blockCenterX := block.X + float64(blockWidth)/2
		blockCenterY := block.Y + float64(blockHeight)/2
		blockHalfWidth := float64(blockWidth) / 2
		blockHalfHeight := float64(blockHeight) / 2

		// ボールの中心からブロックの中心への距離
		dx := g.ball.X - blockCenterX
		dy := g.ball.Y - blockCenterY

		// AABB衝突判定
		if math.Abs(dx) < blockHalfWidth+g.ball.Radius &&
			math.Abs(dy) < blockHalfHeight+g.ball.Radius {
			// ブロックを破壊
			block.Alive = false
			g.score++

			// 反射方向の決定（より正確な反射）
			if math.Abs(dx/blockHalfWidth) > math.Abs(dy/blockHalfHeight) {
				// 左右からの衝突
				g.ball.VX = -g.ball.VX
			} else {
				// 上下からの衝突
				g.ball.VY = -g.ball.VY
			}

			// ボールをブロックの外に移動
			if dx > 0 {
				g.ball.X = block.X + float64(blockWidth) + g.ball.Radius
			} else {
				g.ball.X = block.X - g.ball.Radius
			}
			if dy > 0 {
				g.ball.Y = block.Y + float64(blockHeight) + g.ball.Radius
			} else {
				g.ball.Y = block.Y - g.ball.Radius
			}
		}
	}

	// すべてのブロックが破壊されたかチェック
	allBlocksDestroyed := true
	for _, block := range g.blocks {
		if block.Alive {
			allBlocksDestroyed = false
			break
		}
	}
	if allBlocksDestroyed && len(g.blocks) > 0 {
		g.gameOver = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を黒で塗りつぶし
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// ブロックの描画
	for _, block := range g.blocks {
		if block.Alive {
			ebitenutil.DrawRect(screen, block.X, block.Y, blockWidth, blockHeight, color.RGBA{100, 200, 255, 255})
			ebitenutil.DrawRect(screen, block.X+2, block.Y+2, blockWidth-4, blockHeight-4, color.RGBA{50, 150, 255, 255})
		}
	}

	// パドルの描画
	ebitenutil.DrawRect(screen, g.paddle.X, g.paddle.Y, g.paddle.Width, g.paddle.Height, color.RGBA{255, 255, 255, 255})

	// ボールの描画（円）
	ebitenutil.DrawCircle(screen, g.ball.X, g.ball.Y, g.ball.Radius, color.RGBA{255, 255, 0, 255})

	// スコアの表示
	scoreText := "Score: " + fmt.Sprintf("%d", g.score)
	ebitenutil.DebugPrint(screen, scoreText)

	// ゲームオーバーの表示
	if g.gameOver {
		gameOverText := "GAME OVER"
		allBlocksDestroyed := true
		for _, block := range g.blocks {
			if block.Alive {
				allBlocksDestroyed = false
				break
			}
		}
		if allBlocksDestroyed {
			gameOverText = "YOU WIN!"
		}
		ebitenutil.DebugPrintAt(screen, gameOverText, screenWidth/2-50, screenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Block Game - ブロック崩し")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
