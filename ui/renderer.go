package ui

import (
	"image/color"

	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/input"
	"github.com/kmdkuk/clicker/model"
	"github.com/kmdkuk/clicker/ui/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Renderer interface {
	Draw(screen *ebiten.Image)
	HandleInput(keyType input.KeyType)
	ShowPopup(message string)
	IsPopupActive() bool
	GetPopupMessage() string
	DebugMessage(message string)
	GetDebugMessage() string
}

type DefaultRenderer struct {
	config       *config.Config
	gameState    model.GameStateReader
	debugMessage string
	decider      input.Decider

	// Components for rendering different parts of the UI
	display    *components.Display
	navigation *components.Navigation
	popup      *components.Popup
	manualWork *components.List
	buildings  *components.List
	upgrades   *components.List
	tabs       *components.Tab
	// Add other components as needed
}

func NewRenderer(config *config.Config, gameState model.GameStateReader, decider input.Decider) Renderer {
	return &DefaultRenderer{
		config:       config,
		gameState:    gameState,
		debugMessage: "",
		decider:      decider,
		display:      components.NewDisplay(gameState),
		navigation:   components.NewNavigation(gameState),
		popup:        components.NewPopup(),
		manualWork: components.NewList(gameState, []components.ListItem{
			gameState.GetManualWork(),
		}, true, 10, 50),
		buildings: components.NewList(gameState, components.ConvertBuildingToListItems(gameState.GetBuildings()), false, 10, 90),
		upgrades:  components.NewList(gameState, components.ConvertUpgradeToListItems(gameState.GetUpgrades()), false, 10, 90),
		tabs:      components.NewTab([]string{"Buildings", "Upgrades"}, 0, 20, 70),
	}
}

func (r *DefaultRenderer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill background with black

	// Draw debug information
	if r.config.EnableDebug {
		ebitenutil.DebugPrint(screen, r.debugMessage)
	}

	// Draw game information
	r.display.DrawMoney(screen)

	// If popup is active, only draw it and return
	if r.popup.IsActive() {
		r.popup.Draw(screen)
		return
	}

	r.manualWork.Draw(screen, r.navigation.GetCursor())

	// Draw tabs
	r.tabs.SetActivePage(r.navigation.GetPage()) // Sync tabs with navigation page
	r.tabs.Draw(screen)

	r.buildings.Visible = r.navigation.GetPage() == 0
	r.upgrades.Visible = r.navigation.GetPage() == 1
	r.buildings.Draw(screen, r.navigation.GetCursor()-1)
	r.upgrades.Draw(screen, r.navigation.GetCursor()-1)

}

func (r *DefaultRenderer) HandleInput(keyType input.KeyType) {
	// Popup handling takes priority
	if r.popup.IsActive() {
		r.popup.HandleInput(keyType)
		return
	}

	// Normal input handling
	r.navigation.HandleNavigation(keyType)

	// Decision button handling
	if keyType == input.KeyTypeDecision {
		r.handleDecision()
	}
}

func (r *DefaultRenderer) handleDecision() {
	_, message := r.decider.Decide(
		r.navigation.GetPage(),
		r.navigation.GetCursor(),
	)

	if message != "" {
		r.ShowPopup(message)
	}
}

// Popup related methods
func (r *DefaultRenderer) ShowPopup(message string) {
	r.popup.Show(message)
}

func (r *DefaultRenderer) IsPopupActive() bool {
	return r.popup.IsActive()
}

func (r *DefaultRenderer) GetPopupMessage() string {
	return r.popup.GetMessage()
}

func (r *DefaultRenderer) DebugMessage(message string) {
	r.debugMessage = message
}

func (r *DefaultRenderer) GetDebugMessage() string {
	return r.debugMessage
}
