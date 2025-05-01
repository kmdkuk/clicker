package model

import (
	"fmt"
)

type Upgrade struct {
	Name               string                     // Upgrade name
	Cost               float64                    // Cost to purchase the upgrade
	Effect             func(float64) float64      // Effect applied when the upgrade is purchased
	IsPurchased        bool                       // Whether the upgrade has been purchased
	IsTargetManualWork bool                       // Whether the upgrade is for manual work
	TargetBuilding     int                        // Target building index (if applicable)
	IsReleased         func(GameStateReader) bool // Whether the upgrade is released
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
