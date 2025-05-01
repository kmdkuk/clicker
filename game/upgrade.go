package game

import "fmt"

type Upgrade struct {
	name               string                // Upgrade name
	cost               float64               // Cost to purchase the upgrade
	effect             func(float64) float64 // Effect applied when the upgrade is purchased
	isPurchased        bool                  // Whether the upgrade has been purchased
	isTargetManualWork bool                  // Whether the upgrade is for manual work
	targetBuilding     int                   // Target building index (if applicable)
	isReleased         func(*Game) bool      // Whether the upgrade is released
}

func (u *Upgrade) String(g *Game) string {
	if u.isPurchased {
		return u.name + " (Purchased)"
	}
	if u.isReleased(g) {
		return u.name + " (Selling Cost: $" + fmt.Sprintf("%.2f", u.cost) + ")"
	}
	return u.name + " (Locked Cost: $" + fmt.Sprintf("%.2f", u.cost) + ")"
}
