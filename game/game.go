package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kmdkuk/clicker/config"
)

type Game struct {
	config       *config.Config // Game configuration
	money        float64        // Player's money
	cursor       int            // Cursor position
	manualWork   ManualWork     // Manual work option
	buildings    []Building     // List of buildings
	inputHandler *InputHandler  // Handler to manage input processing
	lastUpdate   time.Time      // Last update time
	popup        Popup          // Popup message
	debugMessage string         // Debug message
}

type Building struct {
	Name         string  // Display name
	BaseCost     float64 // Base cost to unlock
	GenerateRate float64 // Money generated per second
	Count        int     // Number of purchases
}

// Cost method: Calculates the cost based on the current number of purchases
func (b *Building) Cost() float64 {
	// Example of cost increasing exponentially with the number of purchases
	return b.BaseCost * math.Pow(1.15, float64(b.Count))
}

func (b *Building) Unlocked() bool {
	return b.Count != 0
}

func (b *Building) toString() string {
	if b.Unlocked() {
		return fmt.Sprintf("%s (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", b.Name, b.Cost(), b.Count, b.GenerateRate)
	}
	return fmt.Sprintf("%s (Locked, Cost: $%.2f)", b.Name, b.Cost())
}

type ManualWork struct {
	Name  string  // Display name
	Value float64 // Money earned manually
}

func (m *ManualWork) toString() string {
	return m.Name
}

type Popup struct {
	Message string // Message to display
	Active  bool   // Whether the popup is active
}

func (p *Popup) Show(message string) {
	p.Message = message
	p.Active = true
}

func (p *Popup) Close() {
	p.Active = false
}

func NewGame(config *config.Config) *Game {
	return &Game{
		config:     config,
		money:      0, // Initial money
		cursor:     0, // Initial cursor position
		manualWork: ManualWork{Name: "Manual Work: $0.1", Value: 0.1},
		buildings: []Building{
			{Name: "Building 1", BaseCost: 1.0, GenerateRate: 0.01, Count: 0},
			{Name: "Building 2", BaseCost: 10.0, GenerateRate: 0.05, Count: 0},
			{Name: "Building 3", BaseCost: 100.0, GenerateRate: 0.10, Count: 0},
		},
		inputHandler: &InputHandler{},
		lastUpdate:   time.Now(),
	}
}

func (g *Game) Update() error {
	g.inputHandler.Update() // Update input handler
	g.updateBuildings()

	// Handle popup
	if g.popup.Active {
		g.handlePopup()
		return nil
	}

	// Update game state
	g.handleInput()

	return nil
}

func (g *Game) handlePopup() {
	keyType := g.inputHandler.GetPressedKey()
	if keyType == KeyTypeDecision {
		g.popup.Close()
	}
}

func (g *Game) updateBuildings() {
	now := time.Now()
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	for _, building := range g.buildings {
		if building.Unlocked() {
			g.UpdateMoney(building.GenerateRate * float64(building.Count) * elapsed)
		}
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
			building.Count++
		} else {
			if building.Unlocked() {
				g.ShowPopup("Not enough money to purchase!")
			} else {
				g.ShowPopup("Not enough money to unlock!")
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
		ebitenutil.DebugPrintAt(screen, "-> "+g.manualWork.toString(), 10, y)
	} else {
		ebitenutil.DebugPrintAt(screen, "   "+g.manualWork.toString(), 10, y)
	}
}

func (g *Game) drawBuildings(screen *ebiten.Image) {
	for i, building := range g.buildings {
		y := 70 + i*20
		if g.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+building.toString(), 10, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "   "+building.toString(), 10, y)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game screen size
}

func (g *Game) ShowPopup(message string) {
	g.popup.Message = message
	g.popup.Active = true
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
		if building.Unlocked() {
			totalRate += building.GenerateRate * float64(building.Count)
		}
	}
	return totalRate
}

func (g *Game) UpdateMoney(amount float64) {
	g.money += amount
	// Round to avoid floating-point errors
	g.money = math.Round(g.money*10000) / 10000
}
