package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kmdkuk/clicker/application/dto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Display", func() {
	var (
		display    *Display
		mockScreen *ebiten.Image
		playerDTO  *dto.Player
	)

	BeforeEach(func() {
		playerDTO = &dto.Player{
			Money:             123.45,
			TotalGenerateRate: 6.78,
		}
		display = NewDisplay(10, 10)
		mockScreen = ebiten.NewImage(640, 480)
	})

	Describe("NewDisplay", func() {
		It("should initialize with the provided game state", func() {
			Expect(display).NotTo(BeNil())
		})
	})

	Describe("DrawMoney", func() {
		It("should not panic when drawing money information", func() {
			Expect(func() {
				display.DrawMoney(mockScreen, playerDTO)
			}).NotTo(Panic())
		})
	})
})
