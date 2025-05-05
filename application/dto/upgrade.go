package dto

import (
	"github.com/kmdkuk/clicker/presentation/formatter"
)

type Upgrade struct {
	ID          string
	Name        string
	IsPurchased bool
	IsReleased  bool
	Cost        float64
}

func (u *Upgrade) String() string {
	if u.IsPurchased {
		return u.Name + " (Purchased)"
	}
	if u.IsReleased {
		return u.Name + " (Selling Cost: " + formatter.FormatCurrency(u.Cost, "$") + ")"
	}
	return u.Name + " (Locked Cost: " + formatter.FormatCurrency(u.Cost, "$") + ")"
}

func (u *Upgrade) GetName() string {
	return u.Name
}
