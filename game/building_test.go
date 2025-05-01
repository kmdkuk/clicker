package game

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Building", func() {
	var building Building

	BeforeEach(func() {
		building = Building{
			Name:         "Test Building",
			BaseCost:     10.0,
			GenerateRate: 0.5,
			Count:        0,
		}
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
			building.Count = 5
			expectedCost := 10.0
			for i := 0; i < 5; i++ {
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
			building.Count = 1
			Expect(building.IsUnlocked()).To(BeTrue())
		})
	})

	Describe("String", func() {
		It("should return the correct string when locked", func() {
			expected := "Test Building (Locked, Cost: $10.00, Count: 0, Generate Rate: $0.50/s)"
			Expect(building.String()).To(Equal(expected))
		})

		It("should return the correct string when unlocked", func() {
			building.Count = 1
			expected := "Test Building (Next Cost: $11.50, Count: 1, Generate Rate: $0.50/s)"
			Expect(building.String()).To(Equal(expected))
		})
	})

	Describe("GenerateIncome", func() {
		It("should return 0 when the building is locked", func() {
			Expect(building.GenerateIncome(10.0)).To(Equal(0.0))
		})

		It("should calculate the correct income when the building is unlocked", func() {
			building.Count = 2
			Expect(building.GenerateIncome(10.0)).To(Equal(0.5 * 2 * 10.0))
		})
	})
})
