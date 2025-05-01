package state

import (
	"time"

	"github.com/kmdkuk/clicker/level"
	"github.com/kmdkuk/clicker/model"
)

type GameState interface {
	ManualWork()                                             // マニュアルワークを実行します
	UpdateMoney(amount float64)                              // お金を更新します
	PurchaseBuildingAction(buildingIndex int) (bool, string) // 建物を購入します
	PurchaseUpgradeAction(upgradeIndex int) (bool, string)   // アップグレードを購入します
	GetTotalGenerateRate() float64                           // 総生成レートを取得します
	UpdateBuildings(now time.Time)
	GetBuildings() []model.Building
	GetUpgrades() []model.Upgrade
	SetUpgrades(upgrades []model.Upgrade)
	GetMoney() float64
	GetManualWork() *model.ManualWork
}

// GameState はゲームの状態を管理します
type DefaultGameState struct {
	money      float64
	manualWork model.ManualWork
	buildings  []model.Building
	upgrades   []model.Upgrade
	lastUpdate time.Time
}

func NewGameState() GameState {
	return &DefaultGameState{
		money:      0,
		manualWork: model.ManualWork{Name: "Manual Work: $0.1", Value: 0.1, Count: 0},
		buildings:  level.NewBuildings(),
		upgrades:   level.NewUpgrades(),
		lastUpdate: time.Now(),
	}
}

func (g *DefaultGameState) GetBuildings() []model.Building {
	return g.buildings
}
func (g *DefaultGameState) GetUpgrades() []model.Upgrade {
	return g.upgrades
}
func (g *DefaultGameState) SetUpgrades(upgrades []model.Upgrade) {
	g.upgrades = upgrades
}
func (g *DefaultGameState) GetMoney() float64 {
	return g.money
}
func (g *DefaultGameState) GetManualWork() *model.ManualWork {
	return &g.manualWork
}
func (g *DefaultGameState) SetManualWork(manualWork model.ManualWork) {
	g.manualWork = manualWork
}

func (g *DefaultGameState) ManualWork() {
	g.UpdateMoney(g.manualWork.Work(g.upgrades))
}

func (g *DefaultGameState) UpdateMoney(amount float64) {
	g.money += amount
}

func (g *DefaultGameState) PurchaseBuildingAction(buildingIndex int) (bool, string) {
	if buildingIndex < 0 || buildingIndex >= len(g.buildings) {
		return false, "Invalid building selection!"
	}

	building := &g.buildings[buildingIndex]
	cost := building.Cost()

	if g.money < cost {
		if building.IsUnlocked() {
			return false, "Not enough money to purchase!"
		}
		return false, "Not enough money to unlock!"
	}

	g.UpdateMoney(-cost)
	building.Count++
	return true, "Building purchased successfully!"
}

// PurchaseUpgradeAction はアップグレードの購入を試みて結果を返します
func (g *DefaultGameState) PurchaseUpgradeAction(upgradeIndex int) (bool, string) {
	if upgradeIndex < 0 || upgradeIndex >= len(g.upgrades) {
		return false, "Invalid upgrade selection!"
	}

	upgrade := &g.upgrades[upgradeIndex]

	if upgrade.IsPurchased {
		return false, "Upgrade already purchased!"
	}

	if !upgrade.IsReleased(g) {
		return false, "Upgrade not available yet!"
	}

	if g.money < upgrade.Cost {
		return false, "Not enough money for upgrade!"
	}

	g.UpdateMoney(-upgrade.Cost)
	upgrade.IsPurchased = true
	return true, "Upgrade purchased successfully!"
}

func (g *DefaultGameState) GetTotalGenerateRate() float64 {
	totalRate := 0.0
	for _, building := range g.buildings {
		if building.IsUnlocked() {
			totalRate += building.TotalGenerateRate(g.upgrades)
		}
	}
	return totalRate
}

func (g *DefaultGameState) UpdateBuildings(now time.Time) {
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	for _, building := range g.buildings {
		if building.IsUnlocked() {
			g.UpdateMoney(building.GenerateIncome(elapsed, g.upgrades))
		}
	}
}
