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

	for _, block := range state.Blocks {
		if block.Alive {
			ebitenutil.DrawRect(screen, block.X, block.Y, r.layout.BlockW, r.layout.BlockH, color.RGBA{100, 200, 255, 255})
			ebitenutil.DrawRect(screen, block.X+2, block.Y+2, r.layout.BlockW-4, r.layout.BlockH-4, color.RGBA{50, 150, 255, 255})
		}
	}

	ebitenutil.DrawRect(screen, state.Paddle.X, state.Paddle.Y, state.Paddle.Width, state.Paddle.Height, color.RGBA{255, 255, 255, 255})
	ebitenutil.DrawCircle(screen, state.Ball.X, state.Ball.Y, state.Ball.Radius, color.RGBA{255, 255, 0, 255})

	scoreText := "Score: " + fmt.Sprintf("%d", state.Score)
	ebitenutil.DebugPrint(screen, scoreText)

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
