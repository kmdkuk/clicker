package model

import (
	"fmt"
	"math"
)

type Building struct {
	ID               int     // Unique identifier
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
	cost := b.BaseCost * math.Pow(1.15, float64(b.Count))
	return cost
}

func (b *Building) IsUnlocked() bool {
	return b.Count > 0
}

func (b *Building) String(upgrades []Upgrade) string {
	if b.IsUnlocked() {
		return fmt.Sprintf(
			"%s (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
			b.Name,
			b.Cost(),
			b.Count,
			b.TotalGenerateRate(upgrades),
		)
	}
	return fmt.Sprintf(
		"%s (Locked, Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
		b.Name,
		b.Cost(),
		b.Count,
		b.BaseGenerateRate,
	)
}

// TotalGenerateRate メソッドでの丸め処理
func (b *Building) TotalGenerateRate(upgrades []Upgrade) float64 {
	// 計算ロジック
	rate := b.BaseGenerateRate * float64(b.Count)
	// 必要なアップグレード処理
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
