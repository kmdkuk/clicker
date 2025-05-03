package model

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockGameState struct {
	money                        float64
	manualWork                   ManualWork
	buildings                    []Building
	upgrades                     []Upgrade
	manualWorkCalled             bool
	updateBuildingsCalled        bool
	getTotalGenerateRateCalled   bool
	purchaseBuildingActionCalled bool
	purchaseUpgradeActionCalled  bool
}

func (m *MockGameState) UpdateMoney(amount float64) {
	m.money += amount
}

func (m *MockGameState) GetMoney() float64 {
	return m.money
}
func (m *MockGameState) GetManualWork() *ManualWork {
	return &m.manualWork
}
func (m *MockGameState) SetManualWork(manualWork ManualWork) {
	m.manualWork = manualWork
}
func (m *MockGameState) GetBuildings() []Building {
	return m.buildings
}
func (m *MockGameState) GetUpgrades() []Upgrade {
	return m.upgrades
}
func (m *MockGameState) SetUpgrades(upgrades []Upgrade) {
	m.upgrades = upgrades
}
func (m *MockGameState) ManualWorkAction() {
	m.manualWorkCalled = true
	m.UpdateMoney(m.manualWork.Work(m.upgrades))
}
func (m *MockGameState) UpdateBuildings(now time.Time) {
	m.updateBuildingsCalled = true
}
func (m *MockGameState) GetTotalGenerateRate() float64 {
	m.getTotalGenerateRateCalled = true
	return 0.0
}

func (m *MockGameState) PurchaseBuildingAction(buildingIndex int) (bool, string) {
	m.purchaseBuildingActionCalled = true
	return true, ""
}

func (m *MockGameState) PurchaseUpgradeAction(buildingIndex int) (bool, string) {
	m.purchaseUpgradeActionCalled = true
	return true, ""
}

var _ = Describe("ManualWork", func() {
	var (
		manualWork *ManualWork
		upgrades   []Upgrade
		gameState  *MockGameState
	)

	BeforeEach(func() {
		manualWork = &ManualWork{
			Name:      "Manual Click",
			BaseValue: 1.0,
			Count:     0,
		}

		upgrades = []Upgrade{
			{
				Name:               "Manual Boost",
				IsTargetManualWork: true,
				IsPurchased:        true,
				Effect: func(value float64) float64 {
					return value * 2.0
				},
			},
			{
				Name:               "Another Upgrade",
				IsTargetManualWork: true,
				IsPurchased:        false, // Not purchased yet
				Effect: func(value float64) float64 {
					return value * 1.5
				},
			},
		}

		gameState = &MockGameState{
			money: 0,
			buildings: []Building{
				{Name: "Building 1", BaseGenerateRate: 1.0},
				{Name: "Building 2", BaseGenerateRate: 2.0},
			},
			upgrades:                     upgrades,
			manualWorkCalled:             false,
			updateBuildingsCalled:        false,
			getTotalGenerateRateCalled:   false,
			purchaseBuildingActionCalled: false,
			purchaseUpgradeActionCalled:  false,
		}
	})

	Describe("Work", func() {
		It("should increase count when work is performed", func() {
			initialCount := manualWork.Count
			manualWork.Work(upgrades)
			Expect(manualWork.Count).To(Equal(initialCount + 1))
		})

		It("should return the correct value with upgrades applied", func() {
			value := manualWork.Work(upgrades)
			// Only the purchased upgrades should apply
			Expect(value).To(Equal(2.0)) // 1.0 * 2.0
		})
	})

	Describe("UpdateValue", func() {
		It("should apply purchased upgrades to the base value", func() {
			value := manualWork.GetValue(upgrades)
			Expect(value).To(Equal(2.0)) // 1.0 * 2.0
		})

		It("should not apply unpurchased upgrades", func() {
			// Change first upgrade to unpurchased
			upgrades[0].IsPurchased = false
			value := manualWork.GetValue(upgrades)
			Expect(value).To(Equal(1.0)) // No upgrades applied
		})

		It("should apply multiple purchased upgrades cumulatively", func() {
			// Make both upgrades purchased
			upgrades[0].IsPurchased = true
			upgrades[1].IsPurchased = true
			value := manualWork.GetValue(upgrades)
			Expect(value).To(Equal(3.0)) // 1.0 * 2.0 * 1.5
		})

		It("should ignore upgrades not targeting manual work", func() {
			// Add a building upgrade that should be ignored
			buildingUpgrade := Upgrade{
				Name:               "Building Boost",
				IsTargetManualWork: false,
				TargetBuilding:     1,
				IsPurchased:        true,
				Effect: func(value float64) float64 {
					return value * 10.0 // This should be ignored
				},
			}

			upgrades = append(upgrades, buildingUpgrade)
			value := manualWork.GetValue(upgrades)
			Expect(value).To(Equal(2.0)) // Only the manual work upgrade should apply
		})
	})

	Describe("String", func() {
		It("should return formatted string with purchased upgrades applied", func() {
			result := manualWork.String(gameState)
			Expect(result).To(Equal("Manual Click: $2.00"))
		})

		It("should show base value when no upgrades are purchased", func() {
			// Mark all upgrades as not purchased
			for i := range upgrades {
				upgrades[i].IsPurchased = false
			}

			result := manualWork.String(gameState)
			Expect(result).To(Equal("Manual Click: $1.00"))
		})

		It("should apply multiple purchased upgrades correctly", func() {
			// Make both upgrades purchased
			upgrades[0].IsPurchased = true
			upgrades[1].IsPurchased = true

			result := manualWork.String(gameState)
			Expect(result).To(Equal("Manual Click: $3.00"))
		})
	})
})
