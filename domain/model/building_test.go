package model

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type GameStateMock struct {
	Upgrades []Upgrade
	Money    float64
}

func (g GameStateMock) GetMoney() float64 {
	return g.Money
}

func (g GameStateMock) GetTotalGenerateRate() float64 {
	return 0.0
}

func (g GameStateMock) GetManualWork() *ManualWork {
	return nil
}

func (g GameStateMock) GetBuildings() []Building {
	return nil
}

func (g GameStateMock) GetUpgrades() []Upgrade {
	return g.Upgrades
}

func (g *GameStateMock) UpdateMoney(amount float64) {
}

func (g *GameStateMock) SetManualWorkCount(count int) {
}

func (g *GameStateMock) ManualWorkAction() (float64, string) {
	return 0.0, ""
}

func (g *GameStateMock) PurchaseBuildingAction(buildingID int) (bool, string) {
	return false, ""
}

func (g *GameStateMock) PurchaseUpgradeAction(upgradeID int) (bool, string) {
	return false, ""
}

func (g *GameStateMock) SetBuildingCount(buildingID, count int) {
}

func (g *GameStateMock) SetUpgradesIsPurchased(upgradeID int, isPurchased bool) {
}

func (g *GameStateMock) UpdateBuildings(currentTime time.Time) {
}

// GameStateMock is a mock implementation of GameStateReader and GameState interfaces
func NewGameStateMock() *GameStateMock {
	return &GameStateMock{
		Upgrades: []Upgrade{
			{
				Name:               "Upgrade 1",
				Cost:               50.0,
				TargetBuilding:     1,
				IsTargetManualWork: false,
				IsPurchased:        false,
				Effect: func(rate float64) float64 {
					return rate * 2.0
				},
			},
		},
	}
}

func newBuilding() *Building {
	return &Building{
		ID:               1,
		Name:             "Test Building",
		BaseCost:         10.0,
		BaseGenerateRate: 0.5,
		Count:            0,
	}
}

var _ = Describe("Building", func() {
	building := newBuilding()

	BeforeEach(func() {
	})

	Describe("Cost", func() {
		It("should calculate the correct cost for 0 purchases", func() {
			Expect(building.Cost()).To(Equal(10.0))
		})

		It("should calculate the correct cost for 1 purchase", func() {
			building.Count = 1
			Expect(building.Cost()).To(Equal(10.0 * 1.15))
		})

		It("should calculate the correct cost for multiple purchases", func() {
			building.Count = 3
			expectedCost := 10.0 * 1.15 * 1.15 * 1.15
			Expect(building.Cost()).To(BeNumerically("~", expectedCost, 0.00001))
		})
	})

	Describe("IsUnlocked", func() {
		It("should return false when the building is locked", func() {
			building.Count = 0
			Expect(building.IsUnlocked()).To(BeFalse())
		})

		It("should return true when the building is unlocked", func() {
			building.Count = 1
			Expect(building.IsUnlocked()).To(BeTrue())
		})
	})

	Describe("String", func() {
		It("should return the correct string when locked", func() {
			building.Count = 0
			expected := "Test Building (Locked, Cost: $10.00, Count: 0, Generate Rate: $0.50/s)"
			Expect(building.String(nil)).To(Equal(expected))
		})

		It("should return the correct string when unlocked", func() {
			building.Count = 1
			expected := "Test Building (Next Cost: $11.50, Count: 1, Generate Rate: $0.50/s)"
			gameState := NewGameStateMock()
			Expect(building.String(gameState)).To(Equal(expected))
		})

		It("should return the correct string when unlocked with multiple purchases", func() {
			building.Count = 3
			expectedCost := 10.0
			for i := 0; i < building.Count; i++ {
				expectedCost *= 1.15
			}

			expected := fmt.Sprintf("Test Building (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", expectedCost, building.Count, building.BaseGenerateRate*float64(building.Count))
			gameState := NewGameStateMock()
			Expect(building.String(gameState)).To(Equal(expected))
		})
	})

	Describe("GenerateIncome", func() {
		It("should return 0 when the building is locked", func() {
			building.Count = 0
			Expect(building.GenerateIncome(10.0, nil)).To(Equal(0.0))
		})

		It("should calculate the correct income when the building is unlocked", func() {
			building.Count = 2
			expectedIncome := 0.5 * 2 * 10.0
			Expect(building.GenerateIncome(10.0, nil)).To(BeNumerically("~", expectedIncome, 0.001))
		})
	})

	Describe("totalGenerateRate", func() {
		It("should calculate the correct total generate rate without upgrades", func() {
			building.Count = 2
			Expect(building.TotalGenerateRate(nil)).To(Equal(0.5 * 2))
		})

		It("should calculate the correct total generate rate with upgrades", func() {
			building.Count = 2
			upgrades := []Upgrade{
				{IsTargetManualWork: false, TargetBuilding: 1, IsPurchased: true, Effect: func(rate float64) float64 {
					return rate * 1.5
				}},
			}
			Expect(building.TotalGenerateRate(upgrades)).To(BeNumerically("~", 0.5*1.5*2, 0.00001))
		})
	})
})
