package usecase

import (
	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/infrastructure/state"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildingUseCase", func() {
	var (
		gameState *state.DefaultGameState
		useCase   *BuildingUseCase
	)

	BeforeEach(func() {
		gameState = &state.DefaultGameState{
			Money: 1000,
			Buildings: []model.Building{
				{Name: "Building1", BaseCost: 100, Count: 2, BaseGenerateRate: 1.0},
				{Name: "Building2", BaseCost: 200, Count: 1, BaseGenerateRate: 1.0},
				{Name: "Building3", BaseCost: 300, Count: 0, BaseGenerateRate: 1.0},
			},
			Upgrades: []model.Upgrade{},
		}
		useCase = NewBuildingUseCase(gameState)
	})

	Describe("GetBuildings", func() {
		It("should return the correct building information", func() {
			buildings := useCase.GetBuildings()
			Expect(len(buildings)).To(Equal(3))
			building := buildings[0]
			Expect(building.Name).To(Equal("Building1"))
			Expect(building.IsUnlocked).To(BeTrue())
			Expect(building.Count).To(Equal(2))
			Expect(building.Cost).To(BeNumerically("~", 100.0*1.15*1.15, 0.0001))
			Expect(building.TotalGenerateRate).To(Equal(1.0 * 2))
		})
	})

	Describe("PurchaseBuildingAction", func() {
		It("should successfully purchase a building", func() {
			gameState.UpdateMoney(10.0) // Add enough money to purchase
			success, message := useCase.PurchaseBuildingAction(0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal("Building purchased successfully!"))
			Expect(gameState.Buildings[0].Count).To(Equal(3))
		})

		It("should fail to purchase a unlocked building if not enough money", func() {
			gameState.Buildings[0].BaseCost = 2000 // Set cost higher than available money
			success, message := useCase.PurchaseBuildingAction(0)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Not enough money to purchase!"))
		})

		It("should fail to purchase a locked building if not enough money", func() {
			gameState.Buildings[2].BaseCost = 2000 // Set cost higher than available money
			success, message := useCase.PurchaseBuildingAction(2)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Not enough money to unlock!"))
		})

		It("should fail to purchase an invalid building", func() {
			success, message := useCase.PurchaseBuildingAction(-1)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Invalid building selection!"))
		})
	})
})
