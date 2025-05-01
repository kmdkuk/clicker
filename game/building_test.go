package game

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func newBuilding() *Building {
	return &Building{
		id:               1,
		name:             "Test Building",
		baseCost:         10.0,
		baseGenerateRate: 0.5,
		count:            0,
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
			building.count = 1
			Expect(building.Cost()).To(Equal(10.0 * 1.15))
		})

		It("should calculate the correct cost for multiple purchases", func() {
			building.count = 3
			expectedCost := 10.0 * 1.15 * 1.15 * 1.15
			Expect(building.Cost()).To(BeNumerically("~", expectedCost, 0.00001))
		})
	})

	Describe("IsUnlocked", func() {
		It("should return false when the building is locked", func() {
			building.count = 0
			Expect(building.IsUnlocked()).To(BeFalse())
		})

		It("should return true when the building is unlocked", func() {
			building.count = 1
			Expect(building.IsUnlocked()).To(BeTrue())
		})
	})

	Describe("String", func() {
		It("should return the correct string when locked", func() {
			building.count = 0
			expected := "Test Building (Locked, Cost: $10.00, Count: 0, Generate Rate: $0.50/s)"
			Expect(building.String(nil)).To(Equal(expected))
		})

		It("should return the correct string when unlocked", func() {
			building.count = 1
			expected := "Test Building (Next Cost: $11.50, Count: 1, Generate Rate: $0.50/s)"
			Expect(building.String(nil)).To(Equal(expected))
		})

		It("should return the correct string when unlocked with multiple purchases", func() {
			building.count = 3
			expectedCost := 10.0
			for i := 0; i < building.count; i++ {
				expectedCost *= 1.15
			}

			expected := fmt.Sprintf("Test Building (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", expectedCost, building.count, building.baseGenerateRate*float64(building.count))
			Expect(building.String(nil)).To(Equal(expected))
		})
	})

	Describe("GenerateIncome", func() {
		It("should return 0 when the building is locked", func() {
			building.count = 0
			Expect(building.GenerateIncome(10.0, nil)).To(Equal(0.0))
		})

		It("should calculate the correct income when the building is unlocked", func() {
			building.count = 2
			expectedIncome := 0.5 * 2 * 10.0
			Expect(building.GenerateIncome(10.0, nil)).To(BeNumerically("~", expectedIncome, 0.001))
		})
	})

	Describe("totalGenerateRate", func() {
		It("should calculate the correct total generate rate without upgrades", func() {
			building.count = 2
			Expect(building.totalGenerateRate(nil)).To(Equal(0.5 * 2))
		})

		It("should calculate the correct total generate rate with upgrades", func() {
			building.count = 2
			upgrades := []Upgrade{
				{isTargetManualWork: false, targetBuilding: 1, isPurchased: true, effect: func(rate float64) float64 {
					return rate * 1.5
				}},
			}
			Expect(building.totalGenerateRate(upgrades)).To(BeNumerically("~", 0.5*1.5*2, 0.00001))
		})
	})
})
