package dto

import (
	"fmt"

	"github.com/kmdkuk/clicker/presentation/formatter"
)

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
		"%s (%s, Cost: %s, Count: %d, Generate Rate: %s/s)",
		b.Name,
		locked,
		formatter.FormatCurrency(b.Cost, "$"),
		b.Count,
		formatter.FormatCurrency(b.TotalGenerateRate, "$"),
	)
}

func (b *Building) GetName() string {
	return b.Name
}
