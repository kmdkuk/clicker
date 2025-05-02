package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade", func() {
	var (
		upgrade   Upgrade
		gameState *GameStateMock
	)

	BeforeEach(func() {
		gameState = NewGameStateMock()

		upgrade = Upgrade{
			Name:               "Test Upgrade",
			Cost:               50.0,
			TargetBuilding:     1,
			IsTargetManualWork: false,
			IsPurchased:        false,
			Effect: func(rate float64) float64 {
				return rate * 2.0
			},
			IsReleased: func(g GameStateReader) bool {
				// Released based on money amount
				return g.GetMoney() >= 25.0
			},
		}
	})

	Describe("String", func() {
		It("should show as purchased when upgrade is purchased", func() {
			upgrade.IsPurchased = true
			result := upgrade.String(gameState)
			Expect(result).To(ContainSubstring("Purchased"))
		})

		It("should show as available when released but not purchased", func() {
			gameState.Money = 30.0 // Above release threshold
			result := upgrade.String(gameState)
			Expect(result).To(ContainSubstring("Selling Cost"))
		})

		It("should show as locked when not released", func() {
			gameState.Money = 20.0 // Below release threshold
			result := upgrade.String(gameState)
			Expect(result).To(ContainSubstring("Locked"))
		})

		It("should include the upgrade cost in the string", func() {
			result := upgrade.String(gameState)
			Expect(result).To(ContainSubstring("50.00"))
		})
	})
})
