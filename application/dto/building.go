package dto

import "fmt"

type Building struct {
	Name              string
	IsUnlocked        bool
	Cost              float64
	Count             int
	TotalGenerateRate float64
}

func (b *Building) String() string {
	locked := "Locked"
	if b.IsUnlocked {
		locked = "Next"
	}
	return fmt.Sprintf(
		"%s (%s, Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)",
		b.Name,
		locked,
		b.Cost,
		b.Count,
		b.TotalGenerateRate,
	)
}

func (b *Building) GetName() string {
	return b.Name
}
