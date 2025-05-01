package game

import (
	"fmt"
	"math"
)

const roundToTwoDecimalPlaces = 1e2 // 小数点以下2桁に丸めるための定数

func round(value float64) float64 {
	return math.Round(value*roundToTwoDecimalPlaces) / roundToTwoDecimalPlaces
}

type Building struct {
	id               int     // Unique identifier
	name             string  // Display name
	baseCost         float64 // Base cost to unlock
	baseGenerateRate float64 // Money generated per second
	count            int     // Number of purchases
}

// Cost method: Calculates the cost based on the current number of purchases
func (b *Building) Cost() float64 {
	if b.count == 0 {
		return b.baseCost
	}
	return b.baseCost * math.Pow(1.15, float64(b.count))
}

func (b *Building) IsUnlocked() bool {
	return b.count > 0
}

func (b *Building) String(upgrades []Upgrade) string {
	if b.IsUnlocked() {
		return fmt.Sprintf(
			"%s (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
			b.name,
			round(b.Cost()), // 表示時に丸める
			b.count,
			round(b.totalGenerateRate(upgrades)), // 表示時に丸める
		)
	}
	return fmt.Sprintf(
		"%s (Locked, Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
		b.name,
		round(b.Cost()), // 表示時に丸める
		b.count,
		b.baseGenerateRate,
	)
}

func (b *Building) totalGenerateRate(upgrades []Upgrade) float64 {
	totalGenerateRate := b.baseGenerateRate
	for _, upgrade := range upgrades {
		if !upgrade.isTargetManualWork && b.id == upgrade.targetBuilding && upgrade.isPurchased {
			totalGenerateRate = upgrade.effect(totalGenerateRate)
		}
	}
	return totalGenerateRate * float64(b.count) // 丸めを削除
}

func (b *Building) GenerateIncome(elapsed float64, upgrades []Upgrade) float64 {
	if b.IsUnlocked() {
		return b.totalGenerateRate(upgrades) * elapsed // 丸めを削除
	}
	return 0
}
