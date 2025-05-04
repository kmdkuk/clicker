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
		useCase   *ManualWorkUseCase
	)

	BeforeEach(func() {
		gameState = &state.DefaultGameState{
			Money: 1000,
			ManualWork: model.ManualWork{
				Name:      "Manual Work",
				Count:     0,
				BaseValue: 1.0,
			},
			Upgrades: []model.Upgrade{
				{
					Name:        "Upgrade1",
					IsPurchased: true,
					Effect: func(v float64) float64 {
						return v * 1.1
					},
					IsTargetManualWork: true,
				},
			},
		}
		useCase = NewManualWorkUseCase(gameState)
	})

	Describe("GetManualWork", func() {
		It("should return correctly manual work", func() {
			manualWork := useCase.GetManualWork()
			Expect(manualWork.Name).To(Equal("Manual Work"))
			Expect(manualWork.Value).To(BeNumerically("~", 1*1.1, 0.0001))
		})
	})

	Describe("ManualWorkAction", func() {
		It("should update money and manual work count", func() {
			useCase.ManualWorkAction()
			Expect(gameState.Money).To(BeNumerically("~", 1000+1*1.1, 0.0001))
			Expect(gameState.ManualWork.Count).To(Equal(1))
		})
	})
})
