package application

import (
	"errors"

	"block-game/pkg/domain"
)

var ErrNilInputPort = errors.New("input port is nil")

type InputPort interface {
	Read() domain.InputState
}

type GameUsecase struct {
	state  *domain.GameState
	layout domain.LayoutConfig
	input  InputPort
}

func NewGameUsecase(layout domain.LayoutConfig, rnd domain.RandomSource, input InputPort) (*GameUsecase, error) {
	if input == nil {
		return nil, ErrNilInputPort
	}

	blocks, err := domain.GenerateBlocks(layout, rnd)
	if err != nil {
		blocks = domain.GenerateGridFallback(layout)
	}
	state := domain.NewGameState(layout, blocks)

	return &GameUsecase{
		state:  state,
		layout: layout,
		input:  input,
	}, nil
}

func (g *GameUsecase) Update() error {
	domain.Advance(g.state, g.input.Read(), g.layout)
	return nil
}

func (g *GameUsecase) State() *domain.GameState {
	return g.state
}

func (g *GameUsecase) Layout() domain.LayoutConfig {
	return g.layout
}
