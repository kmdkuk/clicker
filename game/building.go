package game

import (
	"fmt"
	"math"
)

type Building struct {
	name             string  // Display name
	baseCost         float64 // Base cost to unlock
	baseGenerateRate float64 // Money generated per second
	count            int     // Number of purchases
}

// Cost method: Calculates the cost based on the current number of purchases
func (b *Building) Cost() float64 {
	// Example of cost increasing exponentially with the number of purchases
	return b.baseCost * math.Pow(1.15, float64(b.count))
}

func (b *Building) IsUnlocked() bool {
	return b.count != 0
}

func (b *Building) String() string {
	if b.IsUnlocked() {
		return fmt.Sprintf("%s (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", b.name, b.Cost(), b.count, b.totalGenerateRate())
	}
	return fmt.Sprintf("%s (Locked, Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", b.name, b.Cost(), b.count, b.baseGenerateRate)
}

func (b *Building) totalGenerateRate() float64 {
	return b.baseGenerateRate * float64(b.count)
}

func (b *Building) GenerateIncome(elapsed float64) float64 {
	if b.IsUnlocked() {
		return b.totalGenerateRate() * elapsed
	}
	return 0
}
