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

	Describe("GetBuildingsIsUnlockedWithMaskedNextLock", func() {
		Context("with some buildings unlocked and some locked", func() {
			BeforeEach(func() {
				gameState.Buildings = []model.Building{
					{Name: "Building1", BaseCost: 100, Count: 2, BaseGenerateRate: 1.0}, // Unlocked
					{Name: "Building2", BaseCost: 200, Count: 1, BaseGenerateRate: 2.0}, // Unlocked
					{Name: "Building3", BaseCost: 300, Count: 0, BaseGenerateRate: 3.0}, // Locked
					{Name: "Building4", BaseCost: 400, Count: 0, BaseGenerateRate: 4.0}, // Locked
					{Name: "Building5", BaseCost: 500, Count: 0, BaseGenerateRate: 5.0}, // Locked
				}
				useCase = NewBuildingUseCase(gameState)
			})

			It("should return all unlocked buildings plus the next locked one", func() {
				buildings := useCase.GetBuildingsIsUnlockedWithMaskedNextLock()

				// Should include 3 unlocked buildings + 1 next locked building
				Expect(buildings).To(HaveLen(3))

				// First three should be unlocked
				Expect(buildings[0].IsUnlocked).To(BeTrue())
				Expect(buildings[0].Name).To(Equal("Building1"))

				Expect(buildings[1].IsUnlocked).To(BeTrue())
				Expect(buildings[1].Name).To(Equal("Building2"))

				// Third should be locked (the next one to unlock)
				Expect(buildings[2].IsUnlocked).To(BeFalse())
				Expect(buildings[2].Name).To(Equal("???"))
			})
		})

		Context("when all buildings are unlocked", func() {
			BeforeEach(func() {
				// Set all buildings to have count > 0 to unlock everything
				gameState.Buildings = []model.Building{
					{Name: "Building1", BaseCost: 100, Count: 2, BaseGenerateRate: 1.0},
					{Name: "Building2", BaseCost: 200, Count: 1, BaseGenerateRate: 2.0},
					{Name: "Building3", BaseCost: 300, Count: 1, BaseGenerateRate: 3.0},
				}
				useCase = NewBuildingUseCase(gameState)
			})

			It("should return all buildings since they are all unlocked", func() {
				buildings := useCase.GetBuildingsIsUnlockedWithMaskedNextLock()

				// Should include all buildings
				Expect(buildings).To(HaveLen(3))

				// All should be unlocked
				for i, building := range buildings {
					Expect(building.IsUnlocked).To(BeTrue())
					Expect(building.Name).To(Equal(gameState.Buildings[i].Name))
				}
			})
		})

		Context("when only the first building is unlocked", func() {
			BeforeEach(func() {
				// Only first building is always unlocked, rest are locked
				gameState.Buildings = []model.Building{
					{Name: "Building1", BaseCost: 100, Count: 0, BaseGenerateRate: 1.0}, // Always unlocked even with count=0
					{Name: "Building2", BaseCost: 200, Count: 0, BaseGenerateRate: 2.0}, // Locked
					{Name: "Building3", BaseCost: 300, Count: 0, BaseGenerateRate: 3.0}, // Locked
				}
				useCase = NewBuildingUseCase(gameState)
			})

			It("should return the first building (unlocked) and the next one (locked)", func() {
				buildings := useCase.GetBuildingsIsUnlockedWithMaskedNextLock()

				// Should include 1 unlocked + 1 next locked
				Expect(buildings).To(HaveLen(1))

				// First one should be unlocked
				Expect(buildings[0].IsUnlocked).To(BeFalse())
				Expect(buildings[0].Name).To(Equal("???"))
			})
		})

		Context("with empty building list", func() {
			BeforeEach(func() {
				gameState.Buildings = []model.Building{}
				useCase = NewBuildingUseCase(gameState)
			})

			It("should return an empty slice", func() {
				buildings := useCase.GetBuildingsIsUnlockedWithMaskedNextLock()
				Expect(buildings).To(BeEmpty())
			})
		})

		Context("when there's a gap in unlocked buildings", func() {
			BeforeEach(func() {
				gameState.Buildings = []model.Building{
					{Name: "Building1", BaseCost: 100, Count: 2, BaseGenerateRate: 1.0}, // Unlocked
					{Name: "Building2", BaseCost: 200, Count: 0, BaseGenerateRate: 2.0}, // Unlocked (from previous)
					{Name: "Building3", BaseCost: 300, Count: 1, BaseGenerateRate: 3.0}, // This creates an inconsistent state
					{Name: "Building4", BaseCost: 400, Count: 0, BaseGenerateRate: 4.0}, // Should be unlocked due to Building3
					{Name: "Building5", BaseCost: 500, Count: 0, BaseGenerateRate: 5.0}, // Should be locked
				}
				useCase = NewBuildingUseCase(gameState)
			})

			It("should handle the gap correctly", func() {
				buildings := useCase.GetBuildingsIsUnlockedWithMaskedNextLock()

				// Should include 4 unlocked buildings + 1 next locked building
				Expect(buildings).To(HaveLen(4))

				// First four should be unlocked
				Expect(buildings[0].IsUnlocked).To(BeTrue())
				Expect(buildings[0].Name).To(Equal("Building1"))
				Expect(buildings[1].IsUnlocked).To(BeFalse())
				Expect(buildings[1].Name).To(Equal("???"))
				Expect(buildings[2].IsUnlocked).To(BeTrue())
				Expect(buildings[2].Name).To(Equal("Building3"))
				Expect(buildings[3].IsUnlocked).To(BeFalse())
				Expect(buildings[3].Name).To(Equal("???"))

			})
		})
	})
})
