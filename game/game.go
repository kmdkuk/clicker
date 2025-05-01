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
	cursor       int          // Cursor position
	page         int          // Page position
	gameState    GameState    // Game state
	decider      Decider      // Decision maker
	inputHandler InputHandler // Handler to manage input processing
	popup        Popup        // Popup message
	debugMessage string       // Debug message
}

func NewGame(config *Config) *Game {
	gameState := NewGameState() // Initialize game state
	game := &Game{
		config:    config,
		cursor:    0, // Initial cursor position
		page:      0, // Initial page position,
		gameState: gameState,
		decider:   NewDefaultDecider(gameState),
	}
	game.validateCursorPosition()
	return game
}

func (g *Game) Update() error {
	g.inputHandler.Update() // Update input handler

	g.gameState.UpdateBuildings(time.Now())

	// Handle popup
	if g.popup.Active {
		g.popup.HandleInput(g.inputHandler.GetPressedKey())
		return nil
	}

	// Update game state
	g.handleInput()

	return nil
}

func (g *Game) handleInput() {
	keyType := g.inputHandler.GetPressedKey()
	g.DebugMessage(fmt.Sprintf("Key Pressed: %v", keyType))
	totalPages := 2                                   // Two pages: manual work + buildings, upgrades
	totalItems := len(g.gameState.GetBuildings()) + 1 // manualWork + buildings
	if g.page == 1 {
		totalItems = len(g.gameState.GetUpgrades()) + 1 // manualWork + upgrades
	}

	switch keyType {
	case KeyTypeUp:
		g.cursor = (g.cursor - 1 + totalItems) % totalItems
	case KeyTypeDown:
		g.cursor = (g.cursor + 1) % totalItems
	case KeyTypeLeft:
		g.page = (g.page - 1 + totalPages) % totalPages // Toggle between pages
	case KeyTypeRight:
		g.page = (g.page + 1) % totalPages
	case KeyTypeDecision:
		g.handleDecision()
	}
	g.validateCursorPosition()
}

// validateCursorPosition はカーソル位置が有効範囲内にあることを確保します
func (g *Game) validateCursorPosition() {
	totalItems := len(g.gameState.GetBuildings()) + 1 // Manual Work + Buildings
	if g.page == 1 {
		totalItems = len(g.gameState.GetUpgrades()) + 1 // Manual Work + Upgrades
	}

	// カーソルが範囲外の場合、安全な値に設定
	if g.cursor < 0 {
		g.cursor = 0
	} else if g.cursor >= totalItems {
		g.cursor = totalItems - 1
	}
}

func (g *Game) handleDecision() {
	_, message := g.decider.Decide(g.page, g.cursor)
	if message != "" {
		g.popup.Show(message)
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
	moneyText := fmt.Sprintf("Money: $%.2f (Total Generate Rate: $%.2f/s)", g.gameState.GetMoney(), g.GetTotalGenerateRate())
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)
}

func (g *Game) drawPopup(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Popup: "+g.popup.Message, 10, 200)
}

func (g *Game) drawManualWork(screen *ebiten.Image) {
	y := 50
	if g.cursor == 0 {
		ebitenutil.DebugPrintAt(screen, "-> "+g.gameState.GetManualWork().String(), 10, y)
		return
	}
	ebitenutil.DebugPrintAt(screen, "   "+g.gameState.GetManualWork().String(), 10, y)
}

func (g *Game) drawBuildings(screen *ebiten.Image) {
	if g.page != 0 {
		return
	}
	for i, building := range g.gameState.GetBuildings() {
		y := 70 + i*20
		if g.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+building.String(g.gameState.GetUpgrades()), 10, y)
			continue
		}
		ebitenutil.DebugPrintAt(screen, "   "+building.String(g.gameState.GetUpgrades()), 10, y)
	}
}

func (g *Game) drawUpgrades(screen *ebiten.Image) {
	if g.page != 1 {
		return
	}
	for i, upgrade := range g.gameState.GetUpgrades() {
		y := 70 + i*20
		if g.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+upgrade.String(g.gameState), 10, y)
			continue
		}
		ebitenutil.DebugPrintAt(screen, "   "+upgrade.String(g.gameState), 10, y)
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
	for _, building := range g.gameState.GetBuildings() {
		if building.IsUnlocked() {
			totalRate += building.totalGenerateRate(g.gameState.GetUpgrades())
		}
	}
	return totalRate
}
