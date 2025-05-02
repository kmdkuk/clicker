package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManualWork", func() {
	var (
		manualWork *ManualWork
		upgrades   []Upgrade
	)

	BeforeEach(func() {
		manualWork = &ManualWork{
			Name:  "Manual Click",
			Value: 1.0,
			Count: 0,
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
			value := manualWork.UpdateValue(upgrades)
			Expect(value).To(Equal(2.0)) // 1.0 * 2.0
		})

		It("should not apply unpurchased upgrades", func() {
			// Change first upgrade to unpurchased
			upgrades[0].IsPurchased = false
			value := manualWork.UpdateValue(upgrades)
			Expect(value).To(Equal(1.0)) // No upgrades applied
		})

		It("should apply multiple purchased upgrades cumulatively", func() {
			// Make both upgrades purchased
			upgrades[0].IsPurchased = true
			upgrades[1].IsPurchased = true
			value := manualWork.UpdateValue(upgrades)
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
			value := manualWork.UpdateValue(upgrades)
			Expect(value).To(Equal(2.0)) // Only the manual work upgrade should apply
		})
	})
})
