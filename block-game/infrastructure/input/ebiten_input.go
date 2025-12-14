package input

import (
	"block-game/domain"

	"github.com/hajimehoshi/ebiten/v2"
)

type EbitenInputAdapter struct{}

func NewEbitenInputAdapter() *EbitenInputAdapter {
	return &EbitenInputAdapter{}
}

func (e *EbitenInputAdapter) Read() domain.InputState {
	return domain.InputState{
		MoveLeft:  ebiten.IsKeyPressed(ebiten.KeyLeft),
		MoveRight: ebiten.IsKeyPressed(ebiten.KeyRight),
	}
}
