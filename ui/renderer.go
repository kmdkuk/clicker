package ui

import (
	"fmt"
	"image/color"

	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/input"
	"github.com/kmdkuk/clicker/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Renderer interface {
	Draw(screen *ebiten.Image)
	HandleInput(keyType input.KeyType)
	HandlePopup(keyType input.KeyType)
	ShowPopup(message string)
	IsPopupActive() bool
	GetPopupMessage() string
	DebugMessage(message string)
	DebugPrint(screen *ebiten.Image)
	GetCursor() int
	GetPage() int
	GetDebugMessage() string
}

type DefaultRenderer struct {
	config       *config.Config
	gameState    model.GameStateReader
	popup        *Popup // popupをRenderer内部に持つ
	debugMessage string
	cursor       int
	page         int
	decider      input.Decider
}

func NewRenderer(config *config.Config, gameState model.GameStateReader, decider input.Decider) Renderer {
	return &DefaultRenderer{
		config:       config,
		gameState:    gameState,
		popup:        NewPopup(),
		debugMessage: "",
		decider:      decider,
	}
}

func (r *DefaultRenderer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the background with black

	r.drawDebug(screen)
	r.drawMoney(screen)
	if r.popup.Active {
		r.drawPopup(screen)
		return
	}
	r.drawManualWork(screen)
	r.drawBuildings(screen)
	r.drawUpgrades(screen)
}

func (r *DefaultRenderer) drawDebug(screen *ebiten.Image) {
	if r.config.EnableDebug {
		ebitenutil.DebugPrint(screen, r.debugMessage)
	}
}

func (r *DefaultRenderer) drawMoney(screen *ebiten.Image) {
	moneyText := fmt.Sprintf("Money: $%.2f (Total Generate Rate: $%.2f/s)", r.gameState.GetMoney(), r.gameState.GetTotalGenerateRate())
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)
}

func (r *DefaultRenderer) drawPopup(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Popup: "+r.popup.Message, 10, 200)
}

func (r *DefaultRenderer) drawManualWork(screen *ebiten.Image) {
	y := 50
	if r.cursor == 0 {
		ebitenutil.DebugPrintAt(screen, "-> "+r.gameState.GetManualWork().String(), 10, y)
		return
	}
	ebitenutil.DebugPrintAt(screen, "   "+r.gameState.GetManualWork().String(), 10, y)
}

func (r *DefaultRenderer) drawBuildings(screen *ebiten.Image) {
	if r.page != 0 {
		return
	}
	for i, building := range r.gameState.GetBuildings() {
		y := 70 + i*20
		if r.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+building.String(r.gameState.GetUpgrades()), 10, y)
			continue
		}
		ebitenutil.DebugPrintAt(screen, "   "+building.String(r.gameState.GetUpgrades()), 10, y)
	}
}

func (r *DefaultRenderer) drawUpgrades(screen *ebiten.Image) {
	if r.page != 1 {
		return
	}
	for i, upgrade := range r.gameState.GetUpgrades() {
		y := 70 + i*20
		if r.cursor == i+1 {
			ebitenutil.DebugPrintAt(screen, "-> "+upgrade.String(r.gameState), 10, y)
			continue
		}
		ebitenutil.DebugPrintAt(screen, "   "+upgrade.String(r.gameState), 10, y)
	}
}

func (r *DefaultRenderer) HandlePopup(keyType input.KeyType) {
	if r.IsPopupActive() && keyType == input.KeyTypeDecision {
		r.popup.Close()
	}
}

// popupに関する操作メソッドを追加
func (r *DefaultRenderer) ShowPopup(message string) {
	r.popup.Show(message)
}

func (r *DefaultRenderer) IsPopupActive() bool {
	return r.popup.Active
}

func (r *DefaultRenderer) GetPopupMessage() string {
	return r.popup.Message
}

func (r *DefaultRenderer) HandleInput(keyType input.KeyType) {
	if r.IsPopupActive() {
		r.HandlePopup(keyType)
		return
	}
	totalPages := 2                                   // Two pages: manual work + buildings, upgrades
	totalItems := len(r.gameState.GetBuildings()) + 1 // manualWork + buildings
	if r.page == 1 {
		totalItems = len(r.gameState.GetUpgrades()) + 1 // manualWork + upgrades
	}

	switch keyType {
	case input.KeyTypeUp:
		r.cursor = (r.cursor - 1 + totalItems) % totalItems
	case input.KeyTypeDown:
		r.cursor = (r.cursor + 1) % totalItems
	case input.KeyTypeLeft:
		r.page = (r.page - 1 + totalPages) % totalPages // Toggle between pages
	case input.KeyTypeRight:
		r.page = (r.page + 1) % totalPages
	case input.KeyTypeDecision:
		r.handleDecision()
	}
	r.validateCursorPosition()
}

// validateCursorPosition はカーソル位置が有効範囲内にあることを確保します
func (r *DefaultRenderer) validateCursorPosition() {
	totalItems := len(r.gameState.GetBuildings()) + 1 // Manual Work + Buildings
	if r.page == 1 {
		totalItems = len(r.gameState.GetUpgrades()) + 1 // Manual Work + Upgrades
	}

	// カーソルが範囲外の場合、安全な値に設定
	if r.cursor < 0 {
		r.cursor = 0
	} else if r.cursor >= totalItems {
		r.cursor = totalItems - 1
	}
}

func (r *DefaultRenderer) handleDecision() {
	_, message := r.decider.Decide(r.page, r.cursor)
	if message != "" {
		r.ShowPopup(message)
	}
}

func (r *DefaultRenderer) DebugMessage(message string) {
	r.debugMessage = message
}

func (r *DefaultRenderer) DebugPrint(screen *ebiten.Image) {
	if r.config.EnableDebug {
		ebitenutil.DebugPrint(screen, r.debugMessage)
	}
}

func (r *DefaultRenderer) GetCursor() int {
	return r.cursor
}
func (r *DefaultRenderer) GetPage() int {
	return r.page
}
func (r *DefaultRenderer) SetCursor(cursor int) {
	r.cursor = cursor
}
func (r *DefaultRenderer) SetPage(page int) {
	r.page = page
}
func (r *DefaultRenderer) GetDebugMessage() string {
	return r.debugMessage
}
