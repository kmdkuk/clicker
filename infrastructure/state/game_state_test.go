package state

import (
	"time"

	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/game/level"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultGameState", func() {
	var gameState DefaultGameState

	BeforeEach(func() {
		gameState = DefaultGameState{
			Money:      0,
			ManualWork: model.ManualWork{Name: "Manual Work: $0.1", BaseValue: 0.1, Count: 0},
			Buildings:  level.NewBuildings(),
			Upgrades:   level.NewUpgrades(),
			LastUpdate: time.Now(),
		} // Update to use gameState
	})

	Describe("UpdateMoney", func() {
		It("should correctly add money", func() {
			gameState.UpdateMoney(10.0)
			Expect(gameState.GetMoney()).To(Equal(10.0))
		})

		It("should correctly subtract money", func() {
			gameState.UpdateMoney(10.0)
			gameState.UpdateMoney(-5.0)
			Expect(gameState.GetMoney()).To(Equal(5.0))
		})
	})

	Describe("updateBuildings", func() {
		It("should generate income from unlocked buildings", func() {
			now := time.Now()
			gameState.Buildings[0].Count = 1                 // Unlock the first building
			gameState.LastUpdate = now.Add(-1 * time.Second) // Simulate 1 second elapsed

			gameState.UpdateBuildings(now)
			Expect(gameState.GetMoney()).To(Equal(gameState.Buildings[0].BaseGenerateRate))
		})

		It("should not generate income from locked buildings", func() {
			now := time.Now()
			gameState.LastUpdate = now.Add(-1 * time.Second) // Simulate 1 second elapsed

			gameState.UpdateBuildings(now)
			Expect(gameState.GetMoney()).To(Equal(0.0))
		})
	})

	Describe("GetTotalGenerateRate", func() {
		It("should calculate the total generate rate from all unlocked buildings", func() {
			gameState.Buildings[0].Count = 1
			gameState.Buildings[1].Count = 2

			expectedRate := gameState.Buildings[0].BaseGenerateRate*1 + gameState.Buildings[1].BaseGenerateRate*2
			Expect(gameState.GetTotalGenerateRate()).To(BeNumerically("~", expectedRate, 0.00001))
		})

		It("should return 0 if no buildings are unlocked", func() {
			Expect(gameState.GetTotalGenerateRate()).To(Equal(0.0))
		})
	})

	Describe("PurchaseBuildingAction", func() {
		It("should successfully purchase a building", func() {
			gameState.UpdateMoney(10.0) // Add enough money to purchase
			success, message := gameState.PurchaseBuildingAction(0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal("Building purchased successfully!"))
			Expect(gameState.Buildings[0].Count).To(Equal(1))
		})

		It("should fail to purchase a building if not enough money", func() {
			success, message := gameState.PurchaseBuildingAction(0)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Not enough money to unlock!"))
		})

		It("should fail to purchase an invalid building", func() {
			success, message := gameState.PurchaseBuildingAction(-1)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Invalid building selection!"))
		})
	})

	Describe("PurchaseUpgradeAction", func() {
		It("should successfully purchase an upgrade", func() {
			gameState.Upgrades[0].IsReleased = func(g model.GameStateReader) bool {
				return true
			}
			gameState.UpdateMoney(10.0) // Add enough money to purchase
			success, message := gameState.PurchaseUpgradeAction(0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal("Upgrade purchased successfully!"))
			Expect(gameState.Upgrades[0].IsPurchased).To(BeTrue())
		})

		It("should fail to purchase an upgrade if not enough money", func() {
			gameState.Upgrades[0].IsReleased = func(g model.GameStateReader) bool {
				return true
			}
			success, message := gameState.PurchaseUpgradeAction(0)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Not enough money for upgrade!"))
		})

		It("should fail to purchase an invalid upgrade", func() {
			success, message := gameState.PurchaseUpgradeAction(-1)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Invalid upgrade selection!"))
		})
	})

})
