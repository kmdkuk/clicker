package usecase

import (
	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/infrastructure/state"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlayerUseCase", func() {
	var (
		gameState *state.DefaultGameState
		useCase   *PlayerUseCase
	)

	BeforeEach(func() {
		gameState = &state.DefaultGameState{
			Money: 1000,
			Buildings: []model.Building{
				{Name: "Building1", BaseCost: 100, Count: 2, BaseGenerateRate: 1.0},
				{Name: "Building2", BaseCost: 200, Count: 1, BaseGenerateRate: 1.0},
				{Name: "Building3", BaseCost: 300, Count: 0, BaseGenerateRate: 1.0},
			},
			Upgrades: []model.Upgrade{
				{
					Name:        "Upgrade1",
					IsPurchased: true,
					Effect: func(v float64) float64 {
						return v * 1.1
					},
					IsTargetManualWork: false,
					TargetBuilding:     0,
				},
			},
		}
		useCase = NewPlayerUsecase(gameState)
	})

	Describe("GetPlayer", func() {
		It("should return correctly Player", func() {
			player := useCase.GetPlayer()
			Expect(player.Money).To(Equal(gameState.Money))
			Expect(player.TotalGenerateRate).To(BeNumerically("~", gameState.GetTotalGenerateRate(), 0.0001))
		})
	})
})
