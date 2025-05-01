package game

import (
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	money        float64       // Player's money
	cursor       int           // Cursor position
	manualWork   ManualWork    // Manual work option
	buildings    []Building    // 建物のリスト
	inputHandler *InputHandler // Handler to manage input processing
	lastUpdate   time.Time     // 最後に更新された時間
}

type Building struct {
	Label        string  // 表示するラベル
	Unlocked     bool    // アンロックされているか
	UnlockCost   float64 // アンロックに必要な金額
	GenerateRate float64 // 毎秒生成する金額
	Count        int     // 購入数
}

func (b *Building) toString() string {
	if b.Unlocked {
		return b.Label + " (Count: " + strconv.Itoa(b.Count) + ")"
	}
	return "Locked ($" + strconv.FormatFloat(b.UnlockCost, 'f', 2, 64) + ")"
}

type ManualWork struct {
	Label string  // 表示するラベル
	Value float64 // 手動で得られる金額
}

func (m *ManualWork) toString() string {
	return m.Label
}

func NewGame() *Game {
	return &Game{
		money:      0.0, // Initial money
		cursor:     0,   // Initial cursor position
		manualWork: ManualWork{Label: "Manual Work: $0.10", Value: 0.10},
		buildings: []Building{
			{Label: "Building 1: $0.05", Unlocked: false, UnlockCost: 1.0, GenerateRate: 0.01, Count: 0},
			{Label: "Building 2: $0.10", Unlocked: false, UnlockCost: 10.0, GenerateRate: 0.05, Count: 0},
			{Label: "Building 3: $0.20", Unlocked: false, UnlockCost: 100.0, GenerateRate: 0.10, Count: 0},
		},
		inputHandler: &InputHandler{},
		lastUpdate:   time.Now(),
	}
}

func (g *Game) Update() error {
	// 経過時間を計算
	now := time.Now()
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// アンロック済みの建物から金を生成
	for _, building := range g.buildings {
		if building.Unlocked {
			g.money += building.GenerateRate * float64(building.Count) * elapsed
		}
	}

	// InputHandler を更新
	g.inputHandler.Update()

	// 押されたキーを分類して取得
	keyType := g.inputHandler.GetPressedKey()
	totalItems := len(g.buildings) + 1 // manualWork + buildings
	switch keyType {
	case KeyTypeUp:
		if g.cursor > 0 {
			g.cursor--
		} else {
			// 一番上にいる場合は一番下に循環
			g.cursor = totalItems - 1
		}
	case KeyTypeDown:
		if g.cursor < (totalItems - 1) {
			g.cursor++
		} else {
			// 一番下にいる場合は一番上に循環
			g.cursor = 0
		}
	case KeyTypeDecision:
		// 決定キーの処理
		if g.cursor == 0 {
			// manualWork の処理
			g.money += g.manualWork.Value
		} else {
			// buildings の処理
			building := &g.buildings[g.cursor-1]
			if building.Unlocked {
				// アンロック済みなら追加購入
				if g.money >= building.UnlockCost {
					g.money -= building.UnlockCost
					building.Count++
				}
			} else if g.money >= building.UnlockCost {
				// アンロック条件を満たしている場合はアンロック
				g.money -= building.UnlockCost
				building.Unlocked = true
				building.Count = 1 // 初回購入
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the background with black

	// Draw the player's money
	moneyText := "Money: $" + strconv.FormatFloat(g.money, 'f', 2, 64)
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)

	// Draw manualWork
	y := 50
	if g.cursor == 0 {
		ebitenutil.DebugPrintAt(screen, "-> "+g.manualWork.toString(), 10, y)
	} else {
		ebitenutil.DebugPrintAt(screen, "   "+g.manualWork.toString(), 10, y)
	}

	// Draw buildings
	for i, building := range g.buildings {
		y := 70 + i*20
		if g.cursor == i+1 {
			// Highlight the position of the cursor
			ebitenutil.DebugPrintAt(screen, "-> "+building.toString(), 10, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "   "+building.toString(), 10, y)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game screen size
}
