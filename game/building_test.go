package game

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Building", func() {
	var building Building

	BeforeEach(func() {
		building = Building{
			name:             "Test Building",
			baseCost:         10.0,
			baseGenerateRate: 0.5,
			count:            0,
		}
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
			building.count = 5
			expectedCost := 10.0
			for i := 0; i < building.count; i++ {
				expectedCost *= 1.15
			}
			Expect(building.Cost()).To(Equal(expectedCost))
		})
	})

	Describe("IsUnlocked", func() {
		It("should return false when the building is locked", func() {
			Expect(building.IsUnlocked()).To(BeFalse())
		})

		It("should return true when the building is unlocked", func() {
			building.count = 1
			Expect(building.IsUnlocked()).To(BeTrue())
		})
	})

	Describe("String", func() {
		It("should return the correct string when locked", func() {
			expected := "Test Building (Locked, Cost: $10.00, Count: 0, Generate Rate: $0.50/s)"
			Expect(building.String()).To(Equal(expected))
		})

		It("should return the correct string when unlocked", func() {
			building.count = 1
			expected := "Test Building (Next Cost: $11.50, Count: 1, Generate Rate: $0.50/s)"
			Expect(building.String()).To(Equal(expected))
		})

		It("should return the correct string when unlocked with multiple purchases", func() {
			building.count = 3
			expectedCost := 10.0
			for i := 0; i < building.count; i++ {
				expectedCost *= 1.15
			}

			expected := fmt.Sprintf("Test Building (Next Cost: $%.2f, Count: %d, Generate Rate: $%.2f/s)", expectedCost, building.count, building.baseGenerateRate*float64(building.count))
			Expect(building.String()).To(Equal(expected))
		})
	})

	Describe("GenerateIncome", func() {
		It("should return 0 when the building is locked", func() {
			Expect(building.GenerateIncome(10.0)).To(Equal(0.0))
		})

		It("should calculate the correct income when the building is unlocked", func() {
			building.count = 2
			Expect(building.GenerateIncome(10.0)).To(Equal(0.5 * 2 * 10.0))
		})
	})
})
