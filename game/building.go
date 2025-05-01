package game

import (
	"fmt"
	"math"
)

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
	cost := b.baseCost * math.Pow(1.15, float64(b.count))
	return cost
}

func (b *Building) IsUnlocked() bool {
	return b.count > 0
}

func (b *Building) String(upgrades []Upgrade) string {
	if b.IsUnlocked() {
		return fmt.Sprintf(
			"%s (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
			b.name,
			b.Cost(),
			b.count,
			b.totalGenerateRate(upgrades),
		)
	}
	return fmt.Sprintf(
		"%s (Locked, Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
		b.name,
		b.Cost(),
		b.count,
		b.baseGenerateRate,
	)
}

// totalGenerateRate メソッドでの丸め処理
func (b *Building) totalGenerateRate(upgrades []Upgrade) float64 {
	// 計算ロジック
	rate := b.baseGenerateRate * float64(b.count)
	// 必要なアップグレード処理
	for _, upgrade := range upgrades {
		if !upgrade.isTargetManualWork && b.id == upgrade.targetBuilding && upgrade.isPurchased {
			rate = upgrade.effect(rate)
		}
	}
	return rate
}

func (b *Building) GenerateIncome(elapsed float64, upgrades []Upgrade) float64 {
	if b.IsUnlocked() {
		return b.totalGenerateRate(upgrades) * elapsed // 丸めを削除
	}
	return 0
}
