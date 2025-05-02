package state

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Save", func() {
	var save Save

	BeforeEach(func() {
		save = Save{
			Money:      100.0,
			Buildings:  []int{1, 2, 3},
			Upgradings: []bool{true},
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
			save.Buildings = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			Expect(save.Validation()).To(HaveOccurred())
		})

		It("should return false if Upgradings length is invalid", func() {
			save.Upgradings = []bool{true, false, true}
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
			save.Upgradings = []bool{true, false, true, true} // 長さが不正
			_, err := save.ConvertToGameState()
			Expect(err).To(HaveOccurred())
		})
	})
})
