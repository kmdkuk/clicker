package state

import (
	"fmt"
	"time"

	"github.com/kmdkuk/clicker/level"
	"github.com/kmdkuk/clicker/model"
)

type GameState interface {
	ManualWorkAction()                                       // マニュアルワークを実行します
	UpdateMoney(amount float64)                              // お金を更新します
	PurchaseBuildingAction(buildingIndex int) (bool, string) // 建物を購入します
	PurchaseUpgradeAction(upgradeIndex int) (bool, string)   // アップグレードを購入します
	GetTotalGenerateRate() float64                           // 総生成レートを取得します
	UpdateBuildings(now time.Time)
	GetBuildings() []model.Building
	SetBuildingCount(buildingIndex int, count int) error
	GetUpgrades() []model.Upgrade
	SetUpgrades(upgrades []model.Upgrade)
	SetUpgradesIsPurchased(upgradeIndex int, isPurchased bool) error
	GetMoney() float64
	GetManualWork() *model.ManualWork
	SetManualWorkCount(count int) error
}

// GameState はゲームの状態を管理します
type DefaultGameState struct {
	Money      float64          `json:"money"`
	ManualWork model.ManualWork `json:"manual_work"`
	Buildings  []model.Building `json:"buildings"`
	Upgrades   []model.Upgrade  `json:"upgrades"`
	LastUpdate time.Time        `json:"last_update"`
}

func NewGameState() GameState {
	return &DefaultGameState{
		Money:      0,
		ManualWork: model.ManualWork{Name: "Manual Work", BaseValue: 0.1, Count: 0},
		Buildings:  level.NewBuildings(),
		Upgrades:   level.NewUpgrades(),
		LastUpdate: time.Now(),
	}
}

func (g *DefaultGameState) GetBuildings() []model.Building {
	return g.Buildings
}

func (g *DefaultGameState) SetBuildingCount(buildingIndex int, count int) error {
	if buildingIndex < 0 || buildingIndex >= len(g.Buildings) {
		return fmt.Errorf("invalid building index: %d", buildingIndex)
	}
	if count < 0 {
		return fmt.Errorf("invalid building count: %d", count)
	}
	g.Buildings[buildingIndex].Count = count
	return nil
}

func (g *DefaultGameState) GetUpgrades() []model.Upgrade {
	return g.Upgrades
}

func (g *DefaultGameState) SetUpgrades(upgrades []model.Upgrade) {
	g.Upgrades = upgrades
}

func (g *DefaultGameState) SetUpgradesIsPurchased(upgradeIndex int, isPurchased bool) error {
	if upgradeIndex < 0 || upgradeIndex >= len(g.Upgrades) {
		return fmt.Errorf("invalid upgrade index: %d", upgradeIndex)
	}
	g.Upgrades[upgradeIndex].IsPurchased = isPurchased
	return nil
}

func (g *DefaultGameState) GetMoney() float64 {
	return g.Money
}
func (g *DefaultGameState) GetManualWork() *model.ManualWork {
	return &g.ManualWork
}
func (g *DefaultGameState) SetManualWork(manualWork model.ManualWork) {
	g.ManualWork = manualWork
}

func (g *DefaultGameState) SetManualWorkCount(count int) error {
	if count < 0 {
		return fmt.Errorf("invalid manual work count: %d", count)
	}
	g.ManualWork.Count = count
	return nil
}

func (g *DefaultGameState) ManualWorkAction() {
	g.UpdateMoney(g.ManualWork.Work(g.Upgrades))
}

func (g *DefaultGameState) UpdateMoney(amount float64) {
	g.Money += amount
}

func (g *DefaultGameState) PurchaseBuildingAction(buildingIndex int) (bool, string) {
	if buildingIndex < 0 || buildingIndex >= len(g.Buildings) {
		return false, "Invalid building selection!"
	}

	building := &g.Buildings[buildingIndex]
	cost := building.Cost()

	if g.Money < cost {
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
	if upgradeIndex < 0 || upgradeIndex >= len(g.Upgrades) {
		return false, "Invalid upgrade selection!"
	}

	upgrade := &g.Upgrades[upgradeIndex]

	if upgrade.IsPurchased {
		return false, "Upgrade already purchased!"
	}

	if !upgrade.IsReleased(g) {
		return false, "Upgrade not available yet!"
	}

	if g.Money < upgrade.Cost {
		return false, "Not enough money for upgrade!"
	}

	g.UpdateMoney(-upgrade.Cost)
	upgrade.IsPurchased = true
	return true, "Upgrade purchased successfully!"
}

func (g *DefaultGameState) GetTotalGenerateRate() float64 {
	totalRate := 0.0
	for _, building := range g.Buildings {
		if building.IsUnlocked() {
			totalRate += building.TotalGenerateRate(g.Upgrades)
		}
	}
	return totalRate
}

func (g *DefaultGameState) UpdateBuildings(now time.Time) {
	elapsed := now.Sub(g.LastUpdate).Seconds()
	g.LastUpdate = now

	for _, building := range g.Buildings {
		if building.IsUnlocked() {
			g.UpdateMoney(building.GenerateIncome(elapsed, g.Upgrades))
		}
	}
}
