package adapter

import (
	"block-game/application"
	"block-game/infrastructure/view"
	"github.com/hajimehoshi/ebiten/v2"
)

type EbitenGame struct {
	usecase  *application.GameUsecase
	renderer *view.Renderer
}

func NewEbitenGame(usecase *application.GameUsecase, renderer *view.Renderer) *EbitenGame {
	return &EbitenGame{
		usecase:  usecase,
		renderer: renderer,
	}
}

func (g *EbitenGame) Update() error {
	return g.usecase.Update()
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.renderer.Render(screen, g.usecase.State())
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	layout := g.usecase.Layout()
	return int(layout.ScreenW), int(layout.ScreenH)
}
