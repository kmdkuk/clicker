package level

import (
	"fmt"

	"github.com/kmdkuk/clicker/domain/model"
)

var building_names = []string{
	"CPU Miner",
	"GPU Rig",
	"ASIC Miner",
	"Mining Farm",
	"Staking Pool",
	"DEX Platform",
	"Layer-2 Network",
	"Blockchain Validator",
	"Quantum Mining Cluster",
	"AI Trading Algorithm",
}

var building_base_costs = []float64{
	0.15,
	1.00,
	11.00,
	120.00,
	1300.00,
	14000.00,
	200000.00,
	3300000.00,
	51000000.00,
	750000000.00,
}

var building_base_generate_rates = []float64{
	0.01,
	0.1,
	0.8,
	4.7,
	26.0,
	140.0,
	780.0,
	4400.0,
	26000.0,
	160000.0,
}

func NewBuildings() []model.Building {
	buildings := make([]model.Building, len(building_names))
	for i := 0; i < len(building_names); i++ {
		buildings[i] = model.Building{
			ID:               i,
			Name:             building_names[i],
			BaseCost:         building_base_costs[i],
			BaseGenerateRate: building_base_generate_rates[i],
			Count:            0,
		}
	}
	return buildings
}

const upgrtade_count_per_unit = 15

var upgrade_unlock_count = []int{
	1,
	5,
	25,
	50,
	100,
	150,
	200,
	250,
	300,
	350,
	400,
	450,
	500,
	550,
	600,
}

var upgrade_base_cost_multiplier = []float64{
	10,
	50,
	500,
	5000,
	50000,
	500000,
	5000000,
	50000000,
	500000000,
	5000000000,
	50000000000,
	500000000000,
	5000000000000,
	50000000000000,
	500000000000000,
}

func newBuildingUpgrade() []model.Upgrade {
	var upgrades []model.Upgrade
	for i := 0; i < len(building_names); i++ {
		for j := 0; j < upgrtade_count_per_unit; j++ {
			upgrades = append(upgrades, model.Upgrade{
				ID:                 fmt.Sprintf("%d_%d", i, j),
				Name:               fmt.Sprintf("%s Upgrade %d", building_names[i], j+1),
				Cost:               building_base_costs[i] * upgrade_base_cost_multiplier[j],
				TargetBuilding:     i,
				IsTargetManualWork: false,
				IsPurchased:        false,
				Effect: func(rate float64) float64 {
					return rate * 2.0
				},
				IsReleased: func(g model.GameStateReader) bool {
					return g.GetBuildings()[i].Count >= upgrade_unlock_count[j]
				},
			})
		}
	}
	return upgrades
}

func newManualWorkUpgrade() []model.Upgrade {
	var upgrades []model.Upgrade
	for i := 0; i < upgrtade_count_per_unit; i++ {
		upgrades = append(upgrades, model.Upgrade{
			ID:                 fmt.Sprintf("manual_work_%d", i),
			Name:               fmt.Sprintf("Manual Work Upgrade %d", i+1),
			Cost:               upgrade_base_cost_multiplier[i],
			TargetBuilding:     -1,
			IsTargetManualWork: true,
			IsPurchased:        false,
			Effect: func(rate float64) float64 {
				return rate * 2.0
			},
			IsReleased: func(g model.GameStateReader) bool {
				return g.GetManualWork().Count >= upgrade_unlock_count[i]
			},
		})
	}
	return upgrades
}

func NewUpgrades() []model.Upgrade {
	upgrades := newBuildingUpgrade()
	upgrades = append(upgrades, newManualWorkUpgrade()...)
	return upgrades
}
