package model

import (
	"math"

	"github.com/kmdkuk/clicker/config"
)

type Building struct {
	ID               int     // Unique identifier for the building
	Name             string  `json:"name"`
	BaseCost         float64 `json:"base_cost"`
	BaseGenerateRate float64 `json:"base_generate_rate"`
	Count            int     `json:"count"`
}

// Cost method: Calculates the cost based on the current number of purchases
func (b *Building) Cost() float64 {
	if b.Count == 0 {
		return b.BaseCost
	}
	cost := b.BaseCost * math.Pow(config.CostMultiplier, float64(b.Count))
	return cost
}

func (b *Building) IsUnlocked() bool {
	return b.Count > 0
}

// TotalGenerateRate method for calculating rounded values
func (b *Building) TotalGenerateRate(upgrades []Upgrade) float64 {
	// Calculation logic
	rate := b.BaseGenerateRate * float64(b.Count)
	// Apply necessary upgrades
	for _, upgrade := range upgrades {
		if !upgrade.IsTargetManualWork && b.ID == upgrade.TargetBuilding && upgrade.IsPurchased {
			rate = upgrade.Effect(rate)
		}
	}
	return rate
}

func (b *Building) GenerateIncome(elapsed float64, upgrades []Upgrade) float64 {
	if b.IsUnlocked() {
		return b.TotalGenerateRate(upgrades) * elapsed // 丸めを削除
	}
	return 0
}
