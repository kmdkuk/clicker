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
	GetCursor() int
	GetPage() int
	DebugMessage(message string)
	GetDebugMessage() string
}

type DefaultRenderer struct {
	config       *config.Config
	gameState    model.GameStateReader
	debugMessage string
	decider      input.Decider

	// コンポーネント
	display    *components.Display
	navigation *components.Navigation
	popup      *components.Popup
	manualWork *components.List
	buildings  *components.List
	upgrades   *components.List
	// 必要に応じて他のコンポーネントを追加
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
		buildings: components.NewList(gameState, components.ConvertBuildingToListItems(gameState.GetBuildings()), false, 10, 70),
		upgrades:  components.NewList(gameState, components.ConvertUpgradeToListItems(gameState.GetUpgrades()), false, 10, 70),
	}
}

func (r *DefaultRenderer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // 背景を黒で塗りつぶし

	// デバッグ情報描画
	if r.config.EnableDebug {
		ebitenutil.DebugPrint(screen, r.debugMessage)
	}

	// ゲーム情報描画
	r.display.DrawMoney(screen)

	// ポップアップがアクティブならそれだけ描画して終了
	if r.popup.IsActive() {
		r.popup.Draw(screen)
		return
	}

	r.manualWork.Draw(screen, r.GetCursor())
	r.buildings.Visible = r.navigation.GetPage() == 0
	r.upgrades.Visible = r.navigation.GetPage() == 1
	r.buildings.Draw(screen, r.GetCursor()-1)
	r.upgrades.Draw(screen, r.GetCursor()-1)
}

func (r *DefaultRenderer) HandleInput(keyType input.KeyType) {
	// ポップアップ処理が優先
	if r.popup.IsActive() {
		r.popup.HandleInput(keyType)
		return
	}

	// 通常の入力処理
	r.navigation.HandleNavigation(keyType)

	// 決定ボタン処理
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

// ポップアップ関連メソッド
func (r *DefaultRenderer) ShowPopup(message string) {
	r.popup.Show(message)
}

func (r *DefaultRenderer) IsPopupActive() bool {
	return r.popup.IsActive()
}

func (r *DefaultRenderer) GetPopupMessage() string {
	return r.popup.GetMessage()
}

// 他のインターフェース実装メソッド
func (r *DefaultRenderer) GetCursor() int {
	return r.navigation.GetCursor()
}

func (r *DefaultRenderer) GetPage() int {
	return r.navigation.GetPage()
}

func (r *DefaultRenderer) DebugMessage(message string) {
	r.debugMessage = message
}

func (r *DefaultRenderer) GetDebugMessage() string {
	return r.debugMessage
}
