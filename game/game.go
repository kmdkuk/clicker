package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	config       *Config      // Game configuration
	money        float64      // Player's money
	cursor       int          // Cursor position
	page         int          // Page position
	manualWork   ManualWork   // Manual work option
	buildings    []Building   // List of buildings
	upgrades     []Upgrade    // List of upgrades
	inputHandler InputHandler // Handler to manage input processing
	lastUpdate   time.Time    // Last update time
	popup        Popup        // Popup message
	debugMessage string       // Debug message
}

func NewGame(config *Config) *Game {
	return &Game{
		config:       config,
		money:        0, // Initial money
		cursor:       0, // Initial cursor position
		page:         0, // Initial page position,
		manualWork:   ManualWork{name: "Manual Work: $0.1", value: 0.1, count: 0},
		buildings:    newBuildings(),
		upgrades:     newUpgrades(),
		inputHandler: &DefaultInputHandler{}, // Use the default implementation
		lastUpdate:   time.Now(),
	}
}

func (g *Game) Update() error {
	g.inputHandler.Update() // Update input handler

	g.updateBuildings(time.Now())

	// Handle popup
	if g.popup.Active {
		g.popup.HandleInput(g.inputHandler.GetPressedKey())
		return nil
	}

	// Update game state
	g.handleInput()

	return nil
}

func (g *Game) updateBuildings(now time.Time) {
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	for _, building := range g.buildings {
		if building.IsUnlocked() {
			g.UpdateMoney(building.GenerateIncome(elapsed, g.upgrades))
		}
	}
}

func (g *Game) handleInput() {
	keyType := g.inputHandler.GetPressedKey()
	g.debugMessage = fmt.Sprintf("Key Pressed: %v", keyType)
	totalPages := 2                    // Two pages: manual work + buildings, upgrades
	totalItems := len(g.buildings) + 1 // manualWork + buildings
	if g.page == 1 {
		totalItems = len(g.upgrades) + 1 // manualWork + upgrades
	}

	switch keyType {
	case KeyTypeUp:
		g.cursor = (g.cursor - 1 + totalItems) % totalItems
	case KeyTypeDown:
		g.cursor = (g.cursor + 1) % totalItems
	case KeyTypeRight:
		g.cursor = 0
		g.page = (g.page + 1) % totalPages // Toggle between pages
	case KeyTypeLeft:
		g.cursor = 0
		g.page = (g.page - 1 + totalPages) % totalPages // Toggle between pages
	case KeyTypeDecision:
		g.handleDecision()
	}
}

func (g *Game) handleDecision() {
	if (g.page == 0 && g.cursor < 0 || g.cursor >= len(g.buildings)+1) ||
		(g.page == 1 && g.cursor < 0 || g.cursor >= len(g.upgrades)+1) { // 無効なカーソル位置をチェック
		return
	}
	if g.cursor == 0 {
		// Manual work processing
		g.manualWork.Work()
		g.UpdateMoney(g.manualWork.Value(g.upgrades))
		return
	}
	switch g.page { // Building processing
	case 0:
		building := &g.buildings[g.cursor-1]
		cost := building.Cost()
		if g.money >= cost {
			g.UpdateMoney(-cost)
			building.count++
		} else {
			if building.IsUnlocked() {
				g.popup.Show("Not enough money to purchase!")
			} else {
				g.popup.Show("Not enough money to unlock!")
			}
			g.DebugMessage(fmt.Sprintf("Cost: %.2f > Money: %.2f", cost, g.money))
		}
	case 1:
		// Upgrade processing
		upgrade := &g.upgrades[g.cursor-1]
		if upgrade.isReleased(g) {
			cost := upgrade.cost
			if upgrade.isPurchased {
				g.popup.Show("Upgrade already purchased!")
				return
			}
			if g.money >= cost {
				g.UpdateMoney(-cost)
				upgrade.isPurchased = true
				g.popup.Show("Upgrade purchased successfully!")
			} else {
				g.popup.Show("Not enough money for upgrade!")
			}
		} else {
			g.popup.Show("Upgrade not available yet!")
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the background with black

	g.drawDebug(screen)
	g.drawMoney(screen)
	if g.popup.Active {
		g.drawPopup(screen)
		return
	}
	g.drawManualWork(screen)
	g.drawBuildings(screen)
	g.drawUpgrades(screen)
}

func (g *Game) drawDebug(screen *ebiten.Image) {
	if g.config.EnableDebug {
		ebitenutil.DebugPrint(screen, g.debugMessage)
	}
}

func (g *Game) drawMoney(screen *ebiten.Image) {
	moneyText := fmt.Sprintf("Money: $%.2f (Total Generate Rate: $%.2f/s)", g.money, g.GetTotalGenerateRate())
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)
}

func (g *Game) drawPopup(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Popup: "+g.popup.Message, 10, 200)
}

func (g *Game) drawManualWork(screen *ebiten.Image) {
	y := 50
	if g.cursor == 0 {
		ebitenutil.DebugPrintAt(screen, "-> "+g.manualWork.String(), 10, y)
	} else {
		ebitenutil.DebugPrintAt(screen, "   "+g.manualWork.String(), 10, y)
	}
}

func (g *Game) drawBuildings(screen *ebiten.Image) {
	if g.page != 0 {
		return
	}
	for i, building := range g.buildings {
		y := 70 + i*20
		if g.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+building.String(g.upgrades), 10, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "   "+building.String(g.upgrades), 10, y)
		}
	}
}

func (g *Game) drawUpgrades(screen *ebiten.Image) {
	if g.page != 1 {
		return
	}
	for i, upgrade := range g.upgrades {
		y := 70 + i*20
		if g.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+upgrade.String(g), 10, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "   "+upgrade.String(g), 10, y)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game screen size
}

func (g *Game) DebugMessage(message string) {
	g.debugMessage = message
}

func (g *Game) DebugPrint(screen *ebiten.Image) {
	if g.config.EnableDebug {
		ebitenutil.DebugPrint(screen, g.debugMessage)
	}
}

// GetTotalGenerateRate calculates the total money generation rate from all unlocked buildings
func (g *Game) GetTotalGenerateRate() float64 {
	totalRate := 0.0
	for _, building := range g.buildings {
		if building.IsUnlocked() {
			totalRate += building.totalGenerateRate(g.upgrades)
		}
	}
	return round(totalRate)
}

func (g *Game) UpdateMoney(amount float64) {
	g.money += amount
	g.money = round(g.money)
}
