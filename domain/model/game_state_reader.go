package model

type GameStateReader interface {
	GetBuildings() []Building
	GetUpgrades() []Upgrade
	GetMoney() float64
	GetManualWork() *ManualWork
	GetTotalGenerateRate() float64
}
