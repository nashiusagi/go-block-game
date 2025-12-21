package adapter

import (
	"fmt"
	"image/color"
	"log"

	"block-game/internal/application"
	"block-game/internal/infrastructure/view"
	"block-game/pkg/config"
	"block-game/pkg/domain"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type gameScene int

const (
	sceneTitle gameScene = iota
	scenePlaying
	scenePaused
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
	prevEscape   bool
	baseLayout   domain.LayoutConfig
	input        application.InputPort
	statusMsg    string
}

func NewEbitenGame(input application.InputPort) *EbitenGame {
	base := config.DefaultLayoutConfig()
	return &EbitenGame{
		usecase:      nil,
		renderer:     nil,
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
		baseLayout: base,
		input:      input,
	}
}

func (g *EbitenGame) Update() error {
	switch g.scene {
	case sceneTitle:
		g.handleTitleInput()
		g.handleTitleMouse()
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			if err := g.startGame(); err != nil {
				log.Printf("failed to start game with difficulty %s: %v", g.selectedDiff, err)
				return nil
			}
			g.scene = scenePlaying
		}
		return nil
	case scenePlaying:
		if g.edgeEscape() {
			g.scene = scenePaused
			return nil
		}
		return g.usecase.Update()
	case scenePaused:
		if g.edgeEscape() {
			g.scene = scenePlaying
		}
		return nil
	default:
		return nil
	}
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	switch g.scene {
	case sceneTitle:
		g.renderTitle(screen)
	case scenePlaying:
		if g.renderer == nil || g.usecase == nil {
			return
		}
		g.renderer.Render(screen, g.usecase.State())
	case scenePaused:
		if g.renderer == nil || g.usecase == nil {
			return
		}
		g.renderer.Render(screen, g.usecase.State())
		g.renderPauseOverlay(screen)
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	layout := g.currentLayout()
	return int(layout.ScreenW), int(layout.ScreenH)
}

func (g *EbitenGame) renderTitle(screen *ebiten.Image) {
	layout := g.currentLayout()
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
	if g.statusMsg != "" {
		ebitenutil.DebugPrintAt(screen, g.statusMsg, startX, startY+96)
	}
}

func (g *EbitenGame) renderPauseOverlay(screen *ebiten.Image) {
	layout := g.currentLayout()
	// 半透明オーバーレイは簡略化のため省略し、テキストのみ表示
	pauseTitle := "PAUSED"
	diffLine := fmt.Sprintf("Difficulty: %s", layout.Difficulty)

	startX := int(layout.ScreenW)/2 - 60
	startY := int(layout.ScreenH)/2 - 20

	ebitenutil.DebugPrintAt(screen, pauseTitle, startX, startY)
	ebitenutil.DebugPrintAt(screen, diffLine, startX-24, startY+16)
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

func (g *EbitenGame) edgeEscape() bool {
	esc := ebiten.IsKeyPressed(ebiten.KeyEscape)
	defer func() { g.prevEscape = esc }()
	return esc && !g.prevEscape
}

func (g *EbitenGame) startGame() error {
	layout, applied, err := config.LayoutWithDifficulty(string(g.selectedDiff))
	if err != nil {
		msg := fmt.Sprintf("fallback to %s (invalid: %s)", applied, g.selectedDiff)
		log.Printf("difficulty selection error: requested=%q fallback=%s err=%v", g.selectedDiff, applied, err)
		g.statusMsg = msg
	}
	if applied != g.selectedDiff {
		g.statusMsg = fmt.Sprintf("fallback to %s (invalid: %s)", applied, g.selectedDiff)
	}
	rnd := domain.NewRandomSource(layout.Seed)

	usecase, err := application.NewGameUsecase(layout, rnd, g.input)
	if err != nil {
		return err
	}

	g.usecase = usecase
	g.renderer = view.NewRenderer(layout)
	g.selectedDiff = applied
	return nil
}

func (g *EbitenGame) currentLayout() domain.LayoutConfig {
	if g.usecase != nil {
		return g.usecase.Layout()
	}
	return g.baseLayout
}
