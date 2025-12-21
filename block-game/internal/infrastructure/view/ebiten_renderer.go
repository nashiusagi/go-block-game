package view

import (
	"fmt"
	"image/color"

	"block-game/pkg/domain"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Renderer struct {
	layout domain.LayoutConfig
}

func NewRenderer(layout domain.LayoutConfig) *Renderer {
	return &Renderer{layout: layout}
}

func (r *Renderer) Render(screen *ebiten.Image, state *domain.GameState) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	diffText := fmt.Sprintf("Difficulty: %s", r.layout.Difficulty)
	ebitenutil.DebugPrintAt(screen, diffText, 0, 0)

	for _, item := range state.Items {
		if item.Active {
			var itemColor color.RGBA
			switch item.Type {
			case domain.ItemTypeMultiball:
				itemColor = color.RGBA{255, 200, 50, 255} // yellow/orange
			case domain.ItemTypePaddleEnlarge:
				itemColor = color.RGBA{50, 255, 100, 255} // green
			default:
				itemColor = color.RGBA{255, 200, 50, 255}
			}
			ebitenutil.DrawRect(screen, item.X, item.Y, item.Width, item.Height, itemColor)
		}
	}

	for _, block := range state.Blocks {
		if block.Alive {
			ebitenutil.DrawRect(screen, block.X, block.Y, r.layout.BlockW, r.layout.BlockH, color.RGBA{100, 200, 255, 255})
			ebitenutil.DrawRect(screen, block.X+2, block.Y+2, r.layout.BlockW-4, r.layout.BlockH-4, color.RGBA{50, 150, 255, 255})
		}
	}

	// Draw paddle with color change when effect is active
	paddleColor := color.RGBA{255, 255, 255, 255} // white (normal)
	if state.PaddleEffect.Active {
		paddleColor = color.RGBA{0, 255, 255, 255} // cyan (enlarged)
	}
	ebitenutil.DrawRect(screen, state.Paddle.X, state.Paddle.Y, state.Paddle.Width, state.Paddle.Height, paddleColor)

	for _, ball := range state.Balls {
		ebitenutil.DrawCircle(screen, ball.X, ball.Y, ball.Radius, color.RGBA{255, 255, 0, 255})
	}

	scoreText := "Score: " + fmt.Sprintf("%d", state.Score)
	ebitenutil.DebugPrintAt(screen, scoreText, 0, 16)

	// Show paddle effect indicator
	if state.PaddleEffect.Active {
		remainingSec := float64(state.PaddleEffect.RemainingTicks) / 60.0
		effectText := fmt.Sprintf("PADDLE x%.0f (%.1fs)", state.PaddleEffect.Multiplier, remainingSec)
		ebitenutil.DebugPrintAt(screen, effectText, 0, 32)
	}

	if state.GameOver {
		gameOverText := "GAME OVER"
		allBlocksDestroyed := true
		for _, block := range state.Blocks {
			if block.Alive {
				allBlocksDestroyed = false
				break
			}
		}
		if allBlocksDestroyed {
			gameOverText = "YOU WIN!"
		}
		ebitenutil.DebugPrintAt(screen, gameOverText, int(r.layout.ScreenW)/2-50, int(r.layout.ScreenH)/2)
	}
}
