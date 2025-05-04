package usecase

import (
	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/infrastructure/state"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpgradeUseCase", func() {
	var (
		gameState *state.DefaultGameState
		useCase   *UpgradeUseCase
	)

	BeforeEach(func() {
		gameState = &state.DefaultGameState{
			Money: 1000,
			Buildings: []model.Building{
				{Name: "Building1", BaseCost: 100, Count: 0},
				{Name: "Building2", BaseCost: 200, Count: 0},
			},
			Upgrades: []model.Upgrade{
				{Name: "Upgrade1", IsPurchased: false, IsReleased: func(g model.GameStateReader) bool { return false }},
			},
		}
		useCase = NewUpgradeUseCase(gameState)
	})

	Describe("PurchaseUpgradeAction", func() {
		It("should successfully purchase an upgrade", func() {
			gameState.Upgrades[0].IsReleased = func(g model.GameStateReader) bool {
				return true
			}
			gameState.UpdateMoney(10.0) // Add enough money to purchase
			success, message := useCase.PurchaseUpgradeAction(0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal("Upgrade purchased successfully!"))
			Expect(gameState.Upgrades[0].IsPurchased).To(BeTrue())
		})

		It("should fail to purchase an upgrade if not enough money", func() {
			gameState.Upgrades[0].IsReleased = func(g model.GameStateReader) bool {
				return true
			}
			success, message := useCase.PurchaseUpgradeAction(0)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Not enough money for upgrade!"))
		})

		It("should fail to purchase an invalid upgrade", func() {
			success, message := useCase.PurchaseUpgradeAction(-1)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Invalid upgrade selection!"))
		})
	})
})
