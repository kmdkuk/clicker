package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Display", func() {
	var (
		display    *Display
		gameState  *GameStateReaderMock
		mockScreen *ebiten.Image
	)

	BeforeEach(func() {
		gameState = &GameStateReaderMock{
			money:        123.45,
			totalGenRate: 6.78,
		}
		display = NewDisplay(gameState)
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
				display.DrawMoney(mockScreen)
			}).NotTo(Panic())
		})
	})
})
