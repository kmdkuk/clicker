package game

import "time"

// GameState はゲームの状態を管理します
type GameState struct {
	money      float64
	manualWork ManualWork
	buildings  []Building
	upgrades   []Upgrade
	lastUpdate time.Time
}

func NewGameState() *GameState {
	return &GameState{
		money:      0,
		manualWork: ManualWork{name: "Manual Work: $0.1", value: 0.1, count: 0},
		buildings:  newBuildings(),
		upgrades:   newUpgrades(),
		lastUpdate: time.Now(),
	}
}

func (g *GameState) Work() {
	g.UpdateMoney(g.manualWork.Work(g.upgrades))
}

func (g *GameState) UpdateMoney(amount float64) {
	g.money += amount
}

func (g *GameState) GetTotalGenerateRate() float64 {
	totalRate := 0.0
	for _, building := range g.buildings {
		if building.IsUnlocked() {
			totalRate += building.totalGenerateRate(g.upgrades)
		}
	}
	return totalRate
}

func (g *GameState) updateBuildings(now time.Time) {
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	for _, building := range g.buildings {
		if building.IsUnlocked() {
			g.UpdateMoney(building.GenerateIncome(elapsed, g.upgrades))
		}
	}
}
