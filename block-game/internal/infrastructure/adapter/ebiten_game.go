package adapter

import (
	"fmt"
	"image/color"

	"block-game/internal/application"
	"block-game/internal/infrastructure/view"
	"block-game/pkg/domain"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type gameScene int

const (
	sceneTitle gameScene = iota
	scenePlaying
)

type EbitenGame struct {
	usecase      *application.GameUsecase
	renderer     *view.Renderer
	scene        gameScene
	selectedDiff domain.Difficulty
	selectedIdx  int
	options      []domain.Difficulty
	descriptions map[domain.Difficulty]string
	prevUp       bool
	prevDown     bool
	prevLeft     bool
	prevRight    bool
}

func NewEbitenGame(usecase *application.GameUsecase, renderer *view.Renderer) *EbitenGame {
	return &EbitenGame{
		usecase:      usecase,
		renderer:     renderer,
		scene:        sceneTitle,
		selectedDiff: domain.DifficultyNormal, // デフォルト表示
		selectedIdx:  1,
		options: []domain.Difficulty{
			domain.DifficultyEasy,
			domain.DifficultyNormal,
			domain.DifficultyHard,
		},
		descriptions: map[domain.Difficulty]string{
			domain.DifficultyEasy:   "球が少し遅くパドルが大きい（やさしめ）",
			domain.DifficultyNormal: "標準の設定",
			domain.DifficultyHard:   "球が速くブロックが多い（チャレンジ）",
		},
	}
}

func (g *EbitenGame) Update() error {
	switch g.scene {
	case sceneTitle:
		g.handleTitleInput()
		g.handleTitleMouse()
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.scene = scenePlaying
		}
		return nil
	case scenePlaying:
		return g.usecase.Update()
	default:
		return nil
	}
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	switch g.scene {
	case sceneTitle:
		g.renderTitle(screen)
	case scenePlaying:
		g.renderer.Render(screen, g.usecase.State())
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	layout := g.usecase.Layout()
	return int(layout.ScreenW), int(layout.ScreenH)
}

func (g *EbitenGame) renderTitle(screen *ebiten.Image) {
	layout := g.usecase.Layout()
	screen.Fill(color.RGBA{0, 0, 0, 255})

	title := "BLOCK GAME"
	prompt := "Enter/Space: Start"

	startX := int(layout.ScreenW)/2 - 120
	startY := int(layout.ScreenH)/2 - 40

	ebitenutil.DebugPrintAt(screen, title, startX+60, startY-24)

	ebitenutil.DebugPrintAt(screen, "Select Difficulty:", startX, startY)
	for i, diff := range g.options {
		lineY := startY + 16*(i+1)
		marker := "  "
		if diff == g.selectedDiff {
			marker = "->"
		}
		desc := g.descriptions[diff]
		text := fmt.Sprintf("%s %s : %s", marker, diff, desc)
		ebitenutil.DebugPrintAt(screen, text, startX, lineY)
	}

	ebitenutil.DebugPrintAt(screen, prompt, startX, startY+80)
}

// handleTitleInput handles keyboard selection (up/down/left/right) on the title screen.
func (g *EbitenGame) handleTitleInput() {
	up := ebiten.IsKeyPressed(ebiten.KeyUp)
	down := ebiten.IsKeyPressed(ebiten.KeyDown)
	left := ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyRight)

	// edge-trigger to avoid rapid repeat while key held down
	if up && !g.prevUp {
		g.moveSelection(-1)
	}
	if down && !g.prevDown {
		g.moveSelection(1)
	}
	if left && !g.prevLeft {
		g.moveSelection(-1)
	}
	if right && !g.prevRight {
		g.moveSelection(1)
	}

	g.prevUp = up
	g.prevDown = down
	g.prevLeft = left
	g.prevRight = right
}

// handleTitleMouse handles mouse hover and click selection on the title screen.
func (g *EbitenGame) handleTitleMouse() {
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}

	x, y := ebiten.CursorPosition()
	layout := g.usecase.Layout()

	startX := int(layout.ScreenW)/2 - 120
	startY := int(layout.ScreenH)/2 - 40

	for i := range g.options {
		lineY := startY + 16*(i+1)
		// ヒットボックス: 行の左端〜右端、行高16px程度
		if x >= startX && x <= startX+240 && y >= lineY-2 && y <= lineY+12 {
			g.selectedIdx = i
			g.selectedDiff = g.options[i]
			break
		}
	}
}

func (g *EbitenGame) moveSelection(delta int) {
	count := len(g.options)
	if count == 0 {
		return
	}
	g.selectedIdx = (g.selectedIdx + delta + count) % count
	g.selectedDiff = g.options[g.selectedIdx]
}
