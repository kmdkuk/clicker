package presentation

import (
	"image/color"

	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/presentation/components"
	"github.com/kmdkuk/clicker/presentation/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Renderer interface {
	Update()
	Draw(screen *ebiten.Image)
	HandleInput(keyType input.KeyType)
	ShowPopup(message string)
	IsPopupActive() bool
	GetPopupMessage() string
	DebugMessage(message string)
	GetDebugMessage() string
}

type PlayerUseCase interface {
	GetPlayer() *dto.Player
}

type ManualWorkUseCase interface {
	ManualWorkAction()
	GetManualWork() *dto.ManualWork
}

type BuildingUseCase interface {
	PurchaseBuildingAction(cursor int) (bool, string)
	GetBuildings() []dto.Building
	GetBuildingsIsUnlockedWithMaskedNextLock() []dto.Building
}

type UpgradeUseCase interface {
	PurchaseUpgradeAction(cursor int) (bool, string)
	GetUpgrades() []dto.Upgrade
	GetUpgradesIsReleasedCostSorted() []dto.Upgrade
}

type DefaultRenderer struct {
	config            *config.Config
	playerUseCase     PlayerUseCase
	manualWorkUseCase ManualWorkUseCase
	buildingUseCase   BuildingUseCase
	upgradeUseCase    UpgradeUseCase
	debugMessage      string
	decider           Decider
	navigation        *Navigation

	// Components for rendering different parts of the UI
	display    *components.Display
	popup      *components.Popup
	manualWork *components.List
	buildings  *components.List
	upgrades   *components.List
	tabs       *components.Tab
	// Add other components as needed
}

func NewRenderer(config *config.Config, playerUseCase PlayerUseCase, manualWorkUseCase ManualWorkUseCase, buildingUseCase BuildingUseCase, upgradeUseCase UpgradeUseCase) Renderer {

	return &DefaultRenderer{
		config:            config,
		playerUseCase:     playerUseCase,
		manualWorkUseCase: manualWorkUseCase,
		buildingUseCase:   buildingUseCase,
		upgradeUseCase:    upgradeUseCase,
		debugMessage:      "",
		decider:           NewDecider(manualWorkUseCase, buildingUseCase, upgradeUseCase),
		navigation:        NewNavigation([]int{len(buildingUseCase.GetBuildings()), len(upgradeUseCase.GetUpgrades())}),
		display:           components.NewDisplay(10, 10),
		popup:             components.NewPopup(),
		manualWork:        components.NewList(true, 10, 50),
		tabs:              components.NewTab([]string{"Buildings", "Upgrades"}, 0, 10, 90),
		buildings:         components.NewList(false, 10, 130),
		upgrades:          components.NewList(false, 10, 130),
	}
}

func (r *DefaultRenderer) Update() {
	r.manualWork.Items = []components.ListItem{
		r.manualWorkUseCase.GetManualWork(),
	}
	r.buildings.Items = components.ConvertBuildingToListItems(r.buildingUseCase.GetBuildingsIsUnlockedWithMaskedNextLock())
	r.upgrades.Items = components.ConvertUpgradeToListItems(r.upgradeUseCase.GetUpgradesIsReleasedCostSorted())

	r.navigation.totalItems = []int{
		len(r.buildings.Items),
		len(r.upgrades.Items),
	}
}

func (r *DefaultRenderer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill background with black

	// Draw debug information
	if r.config.EnableDebug {
		ebitenutil.DebugPrint(screen, r.debugMessage)
	}

	// Draw game information
	r.display.DrawMoney(screen, r.playerUseCase.GetPlayer())

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
