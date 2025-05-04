package storage

import (
	"github.com/kmdkuk/clicker/game/level"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Save", func() {
	var save Save

	BeforeEach(func() {
		save = Save{
			Money:     100.0,
			Buildings: []int{1, 2, 3},
			Upgradings: []upgrade{
				{
					ID:          "0_0",
					IsPurchased: true,
				},
			},
			ManualWork: 10,
		}
	})

	Describe("IsValid", func() {
		It("should return true for a valid Save", func() {
			Expect(save.Validation()).ToNot(HaveOccurred())
		})

		It("should return false if Money is negative", func() {
			save.Money = -1
			Expect(save.Validation()).To(HaveOccurred())
		})

		It("should return false if Buildings length is invalid", func() {
			var buildings []int
			for range level.NewBuildings() {
				buildings = append(buildings, 0)
			}
			buildings = append(buildings, 0)
			save.Buildings = buildings
			Expect(save.Validation()).To(HaveOccurred())
		})

		It("should return false if Upgradings length is invalid", func() {
			us := level.NewUpgrades()
			var upgrades []upgrade
			for _, u := range us {
				upgrades = append(upgrades, upgrade{
					ID:          u.ID,
					IsPurchased: false,
				})
			}
			upgrades = append(upgrades, upgrade{
				ID:          "invalid_upgrade",
				IsPurchased: true,
			})
			save.Upgradings = upgrades
			Expect(save.Validation()).To(HaveOccurred())
		})

		It("should return false if ManualWork is negative", func() {
			save.ManualWork = -1
			Expect(save.Validation()).To(HaveOccurred())
		})
	})

	Describe("ConvertToGameState", func() {
		It("should convert Save to GameState successfully", func() {
			gameState, err := save.ConvertToGameState()
			Expect(err).ToNot(HaveOccurred())
			Expect(gameState.GetMoney()).To(Equal(save.Money))
			Expect(gameState.GetManualWork().Count).To(Equal(save.ManualWork))
			// 他のフィールドも必要に応じて検証
		})

		It("should return an error if setting ManualWork fails", func() {
			save.ManualWork = -1 // 無効な値を設定
			_, err := save.ConvertToGameState()
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if setting Buildings fails", func() {
			save.Buildings = []int{-1} // 無効な値を設定
			_, err := save.ConvertToGameState()
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if setting Upgradings fails", func() {
			save.Upgradings = []upgrade{
				{
					ID:          "0_0",
					IsPurchased: true,
				},
				{
					ID:          "0_1",
					IsPurchased: false,
				},
				{
					ID:          "1_0",
					IsPurchased: true,
				},
				{
					ID:          "invalid_upgrade",
					IsPurchased: true,
				},
			}
			_, err := save.ConvertToGameState()
			Expect(err).To(HaveOccurred())
		})
	})
})
