package game

import (
	"fmt"
	"math"
)

type Building struct {
	Name         string  // Display name
	BaseCost     float64 // Base cost to unlock
	GenerateRate float64 // Money generated per second
	Count        int     // Number of purchases
}

// Cost method: Calculates the cost based on the current number of purchases
func (b *Building) Cost() float64 {
	// Example of cost increasing exponentially with the number of purchases
	return b.BaseCost * math.Pow(1.15, float64(b.Count))
}

func (b *Building) IsUnlocked() bool {
	return b.Count != 0
}

func (b *Building) String() string {
	if b.IsUnlocked() {
		return fmt.Sprintf("%s (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", b.Name, b.Cost(), b.Count, b.GenerateRate)
	}
	return fmt.Sprintf("%s (Locked, Cost: $%.2f)", b.Name, b.Cost())
}

func (b *Building) GenerateIncome(elapsed float64) float64 {
	if b.IsUnlocked() {
		return b.GenerateRate * float64(b.Count) * elapsed
	}
	return 0
}
