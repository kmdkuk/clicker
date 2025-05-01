package game

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	money        float64       // Player's money
	cursor       int           // Cursor position (0: top, 1: middle, 2: bottom)
	buildings    []Building    // 建物のリスト
	inputHandler *InputHandler // Handler to manage input processing
}

type Building struct {
	Label string  // 表示するラベル
	Value float64 // 増加する金額
}

func NewGame() *Game {
	return &Game{
		money:  0.0, // Initial money
		cursor: 0,   // Initial cursor position
		buildings: []Building{
			{Label: "Building 1: $0.05", Value: 0.05},
			{Label: "Building 2: $0.10", Value: 0.10},
			{Label: "Building 3: $0.20", Value: 0.20},
		},
		inputHandler: &InputHandler{},
	}
}

func (g *Game) Update() error {
	// InputHandler を更新
	g.inputHandler.Update()

	// 押されたキーを分類して取得
	keyType := g.inputHandler.GetPressedKey()
	switch keyType {
	case KeyTypeUp:
		if g.cursor > 0 {
			g.cursor--
		} else {
			// 一番上にいる場合は一番下に循環
			g.cursor = len(g.buildings) - 1
		}
	case KeyTypeDown:
		if g.cursor < len(g.buildings)-1 {
			g.cursor++
		} else {
			// 一番下にいる場合は一番上に循環
			g.cursor = 0
		}
	case KeyTypeDecision:
		// 決定キーの処理
		g.money += g.buildings[g.cursor].Value
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the background with black

	// Draw the player's money
	moneyText := "Money: $" + strconv.FormatFloat(g.money, 'f', 2, 64)
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)

	// Draw menu options
	for i, building := range g.buildings {
		y := 50 + i*20
		if i == g.cursor {
			// Highlight the position of the cursor
			ebitenutil.DebugPrintAt(screen, "-> "+building.Label, 10, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "   "+building.Label, 10, y)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game screen size
}
