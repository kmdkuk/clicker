package game

func newBuildings() []Building {
	return []Building{
		{name: "Building 1", baseCost: 1.0, baseGenerateRate: 0.01, count: 0},
		{name: "Building 2", baseCost: 10.0, baseGenerateRate: 0.1, count: 0},
		{name: "Building 3", baseCost: 100.0, baseGenerateRate: 1, count: 0},
		{name: "Building 4", baseCost: 1000.0, baseGenerateRate: 10, count: 0},
		{name: "Building 5", baseCost: 10000.0, baseGenerateRate: 100, count: 0},
		{name: "Building 6", baseCost: 100000.0, baseGenerateRate: 1000, count: 0},
		{name: "Building 7", baseCost: 1000000.0, baseGenerateRate: 10000, count: 0},
	}
}

func newUpgrades() []Upgrade {
	return []Upgrade{
		{name: "Manual Work Upgrade 1", cost: 10,
			effect: func(value float64) float64 {
				return value * 1.1
			},
			isPurchased: false, isTargetManualWork: true, targetBuilding: -1,
			isReleased: func(g GameState) bool {
				return g.GetManualWork().count >= 10
			},
		},
	}
}
