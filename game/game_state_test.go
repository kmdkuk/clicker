package game

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GameState", func() {
	var gameState *GameState

	BeforeEach(func() {
		gameState = NewGameState() // Update to use gameState
	})

	Describe("UpdateMoney", func() {
		It("should correctly add money", func() {
			gameState.UpdateMoney(10.0)
			Expect(gameState.money).To(Equal(10.0))
		})

		It("should correctly subtract money", func() {
			gameState.UpdateMoney(10.0)
			gameState.UpdateMoney(-5.0)
			Expect(gameState.money).To(Equal(5.0))
		})
	})

	Describe("updateBuildings", func() {
		It("should generate income from unlocked buildings", func() {
			now := time.Now()
			gameState.buildings[0].count = 1                 // Unlock the first building
			gameState.lastUpdate = now.Add(-1 * time.Second) // Simulate 1 second elapsed

			gameState.updateBuildings(now)
			Expect(gameState.money).To(Equal(gameState.buildings[0].baseGenerateRate))
		})

		It("should not generate income from locked buildings", func() {
			now := time.Now()
			gameState.lastUpdate = now.Add(-1 * time.Second) // Simulate 1 second elapsed

			gameState.updateBuildings(now)
			Expect(gameState.money).To(Equal(0.0))
		})
	})

	Describe("GetTotalGenerateRate", func() {
		It("should calculate the total generate rate from all unlocked buildings", func() {
			gameState.buildings[0].count = 1
			gameState.buildings[1].count = 2

			expectedRate := gameState.buildings[0].baseGenerateRate*1 + gameState.buildings[1].baseGenerateRate*2
			Expect(gameState.GetTotalGenerateRate()).To(BeNumerically("~", expectedRate, 0.00001))
		})

		It("should return 0 if no buildings are unlocked", func() {
			Expect(gameState.GetTotalGenerateRate()).To(Equal(0.0))
		})
	})

})
