package model

import (
	"fmt"
)

type Upgrade struct {
	Name               string                     `json:"name"`
	Cost               float64                    `json:"cost"`
	Effect             func(float64) float64      `json:"-"` // Exclude from JSON encoding
	IsPurchased        bool                       `json:"is_purchased"`
	IsTargetManualWork bool                       `json:"is_target_manual_work"`
	TargetBuilding     int                        `json:"target_building"`
	IsReleased         func(GameStateReader) bool `json:"-"` // Exclude from JSON encoding
}

func (u *Upgrade) String(g GameStateReader) string {
	if u.IsPurchased {
		return u.Name + " (Purchased)"
	}
	if u.IsReleased(g) {
		return u.Name + " (Selling Cost: $" + fmt.Sprintf("%.2f", u.Cost) + ")"
	}
	return u.Name + " (Locked Cost: $" + fmt.Sprintf("%.2f", u.Cost) + ")"
}
