package level

import "github.com/kmdkuk/clicker/domain/model"

func NewBuildings() []model.Building {
	return []model.Building{
		{ID: 0, Name: "Building 1", BaseCost: 1.0, BaseGenerateRate: 0.01, Count: 0},
		{ID: 1, Name: "Building 2", BaseCost: 10.0, BaseGenerateRate: 0.1, Count: 0},
		{ID: 2, Name: "Building 3", BaseCost: 100.0, BaseGenerateRate: 1, Count: 0},
		{ID: 3, Name: "Building 4", BaseCost: 1000.0, BaseGenerateRate: 10, Count: 0},
		{ID: 4, Name: "Building 5", BaseCost: 10000.0, BaseGenerateRate: 100, Count: 0},
		{ID: 5, Name: "Building 6", BaseCost: 100000.0, BaseGenerateRate: 1000, Count: 0},
		{ID: 6, Name: "Building 7", BaseCost: 1000000.0, BaseGenerateRate: 10000, Count: 0},
		{ID: 7, Name: "Building 8", BaseCost: 10000000.0, BaseGenerateRate: 100000, Count: 0},
	}
}

func NewUpgrades() []model.Upgrade {
	return []model.Upgrade{
		{Name: "Manual Work Upgrade 1", Cost: 10,
			Effect: func(value float64) float64 {
				return value * 1.1
			},
			IsPurchased: false, IsTargetManualWork: true, TargetBuilding: -1,
			IsReleased: func(g model.GameStateReader) bool {
				return g.GetManualWork().Count >= 10
			},
		},
	}
}
