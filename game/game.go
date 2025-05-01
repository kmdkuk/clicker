package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	config       *Config      // Game configuration
	money        float64      // Player's money
	cursor       int          // Cursor position
	manualWork   ManualWork   // Manual work option
	buildings    []Building   // List of buildings
	inputHandler InputHandler // Handler to manage input processing
	lastUpdate   time.Time    // Last update time
	popup        Popup        // Popup message
	debugMessage string       // Debug message
}

func NewGame(config *Config) *Game {
	return &Game{
		config:     config,
		money:      0, // Initial money
		cursor:     0, // Initial cursor position
		manualWork: ManualWork{Name: "Manual Work: $0.1", Value: 0.1},
		buildings: []Building{
			{name: "Building 1", baseCost: 1.0, baseGenerateRate: 0.01, count: 0},
			{name: "Building 2", baseCost: 10.0, baseGenerateRate: 0.1, count: 0},
			{name: "Building 3", baseCost: 100.0, baseGenerateRate: 1, count: 0},
			{name: "Building 4", baseCost: 1000.0, baseGenerateRate: 10, count: 0},
			{name: "Building 5", baseCost: 10000.0, baseGenerateRate: 100, count: 0},
			{name: "Building 6", baseCost: 100000.0, baseGenerateRate: 1000, count: 0},
			{name: "Building 7", baseCost: 1000000.0, baseGenerateRate: 10000, count: 0},
		},
		inputHandler: &DefaultInputHandler{}, // Use the default implementation
		lastUpdate:   time.Now(),
	}
}

func (g *Game) Update() error {
	g.inputHandler.Update() // Update input handler

	g.updateBuildings()

	// Handle popup
	if g.popup.Active {
		g.popup.HandleInput(g.inputHandler.GetPressedKey())
		return nil
	}

	// Update game state
	g.handleInput()

	return nil
}

func (g *Game) updateBuildings() {
	now := time.Now()
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	for _, building := range g.buildings {
		g.UpdateMoney(building.GenerateIncome(elapsed))
	}
}

func (g *Game) handleInput() {
	keyType := g.inputHandler.GetPressedKey()
	totalItems := len(g.buildings) + 1 // manualWork + buildings

	switch keyType {
	case KeyTypeUp:
		g.cursor = (g.cursor - 1 + totalItems) % totalItems
	case KeyTypeDown:
		g.cursor = (g.cursor + 1) % totalItems
	case KeyTypeDecision:
		g.handleDecision()
	}
}

func (g *Game) handleDecision() {
	if g.cursor == 0 {
		// Manual work processing
		g.UpdateMoney(g.manualWork.Value)
	} else {
		// Building processing
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
	for i, building := range g.buildings {
		y := 70 + i*20
		if g.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+building.String(), 10, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "   "+building.String(), 10, y)
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
			totalRate += building.totalGenerateRate()
		}
	}
	return totalRate
}

func (g *Game) UpdateMoney(amount float64) {
	g.money += amount
	// Round to avoid floating-point errors
	g.money = math.Round(g.money*10000) / 10000
}
