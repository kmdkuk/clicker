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
				{Name: "Building1", BaseCost: 100, Count: 0},
				{Name: "Building2", BaseCost: 200, Count: 0},
			},
			Upgrades: []model.Upgrade{},
		}
		useCase = NewBuildingUseCase(gameState)
	})

	Describe("PurchaseBuildingAction", func() {
		It("should successfully purchase a building", func() {
			gameState.UpdateMoney(10.0) // Add enough money to purchase
			success, message := useCase.PurchaseBuildingAction(0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal("Building purchased successfully!"))
			Expect(gameState.Buildings[0].Count).To(Equal(1))
		})

		It("should fail to purchase a building if not enough money", func() {
			success, message := useCase.PurchaseBuildingAction(0)
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
