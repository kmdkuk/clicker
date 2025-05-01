package game

import (
	"image/color"
	"math"
	"strconv"
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
	str := b.Name + " ("
	if b.Unlocked() {
		str += "Next "
	} else {
		str += "Locked,"
	}
	str += "Cost: " + strconv.FormatFloat(b.Cost(), 'f', 2, 64) + ", Count: " + strconv.Itoa(b.Count) + ", GenerateRate: $" + strconv.FormatFloat(b.GenerateRate, 'f', 2, 64) + "/s)"
	return str
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
	// Update the InputHandler
	g.inputHandler.Update()

	// If a popup is active
	if g.popup.Active {
		keyType := g.inputHandler.GetPressedKey()
		if keyType == KeyTypeDecision {
			g.popup.Active = false // Close the popup
		}
		return nil
	}

	// Regular update process
	now := time.Now()
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// Generate money from unlocked buildings
	for _, building := range g.buildings {
		if building.Unlocked() {
			g.UpdateMoney(building.GenerateRate * float64(building.Count) * elapsed)
		}
	}

	// Get the pressed key
	keyType := g.inputHandler.GetPressedKey()
	totalItems := len(g.buildings) + 1 // manualWork + buildings
	switch keyType {
	case KeyTypeUp:
		if g.cursor > 0 {
			g.cursor--
		} else {
			// If at the top, loop to the bottom
			g.cursor = totalItems - 1
		}
	case KeyTypeDown:
		if g.cursor < (totalItems - 1) {
			g.cursor++
		} else {
			// If at the bottom, loop to the top
			g.cursor = 0
		}
	case KeyTypeDecision:
		// Decision key processing
		if g.cursor == 0 {
			// Manual work processing
			g.UpdateMoney(g.manualWork.Value)
		} else {
			// Building processing
			building := &g.buildings[g.cursor-1]
			cost := building.Cost() // Get the current purchase cost
			if g.money >= cost {
				g.UpdateMoney(-cost)
				building.Count++
			} else {
				if building.Unlocked() {
					g.ShowPopup("Not enough money to purchase!")
				} else {
					g.ShowPopup("Not enough money to unlock!")
				}
				g.DebugMessage("Cost: " + strconv.FormatFloat(cost, 'f', 2, 64) + " > money: " + strconv.FormatFloat(g.money, 'f', 2, 64))
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the background with black
	// Draw the debug message if in debug mode
	g.DebugPrint(screen)

	// Draw the player's money
	moneyText := "Money: $" + strconv.FormatFloat(g.money, 'f', 2, 64)
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)

	// If a popup is active
	if g.popup.Active {
		ebitenutil.DebugPrintAt(screen, "Popup: "+g.popup.Message, 10, 200)
		return
	}

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

func (g *Game) UpdateMoney(amount float64) {
	g.money += amount
	// Round to avoid floating-point errors
	g.money = math.Round(g.money*10000) / 10000
}
