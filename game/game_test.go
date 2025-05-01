package game

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockInputHandler is a mock implementation of InputHandler
type MockInputHandler struct {
	pressedKey KeyType
}

func (m *MockInputHandler) Update() {
	// No-op for mock
}

func (m *MockInputHandler) GetPressedKey() KeyType {
	return m.pressedKey
}

var _ = Describe("Game", func() {
	var game *Game
	var mockInputHandler *MockInputHandler

	BeforeEach(func() {
		config := &Config{EnableDebug: false}
		mockInputHandler = &MockInputHandler{}
		game = NewGame(config)
		game.inputHandler = mockInputHandler // Replace inputHandler with mock
	})

	Describe("UpdateMoney", func() {
		It("should correctly add money", func() {
			game.UpdateMoney(10.0)
			Expect(game.money).To(Equal(10.0))
		})

		It("should correctly subtract money", func() {
			game.UpdateMoney(10.0)
			game.UpdateMoney(-5.0)
			Expect(game.money).To(Equal(5.0))
		})

		It("should round money to avoid floating-point errors", func() {
			game.UpdateMoney(0.1)
			game.UpdateMoney(0.2)
			Expect(game.money).To(Equal(0.3))
		})
	})

	Describe("updateBuildings", func() {
		It("should generate income from unlocked buildings", func() {
			game.buildings[0].Count = 1                        // Unlock the first building
			game.lastUpdate = time.Now().Add(-1 * time.Second) // Simulate 1 second elapsed

			game.updateBuildings()
			Expect(game.money).To(Equal(game.buildings[0].GenerateRate))
		})

		It("should not generate income from locked buildings", func() {
			game.lastUpdate = time.Now().Add(-1 * time.Second) // Simulate 1 second elapsed

			game.updateBuildings()
			Expect(game.money).To(Equal(0.0))
		})
	})

	Describe("handleDecision", func() {
		It("should add money for manual work", func() {
			game.cursor = 0 // Select manual work
			game.handleDecision()
			Expect(game.money).To(Equal(game.manualWork.Value))
		})

		It("should purchase a building if enough money is available", func() {
			game.cursor = 1        // Select the first building
			game.UpdateMoney(10.0) // Add enough money to purchase
			game.handleDecision()

			Expect(game.money).To(BeNumerically("<", 10.0)) // Money should decrease
			Expect(game.buildings[0].Count).To(Equal(1))    // Building count should increase
		})

		It("should show a popup if not enough money is available", func() {
			game.cursor = 1 // Select the first building
			game.handleDecision()

			Expect(game.popup.Active).To(BeTrue()) // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money to unlock!"))
		})
	})

	Describe("GetTotalGenerateRate", func() {
		It("should calculate the total generate rate from all unlocked buildings", func() {
			game.buildings[0].Count = 1
			game.buildings[1].Count = 2

			expectedRate := game.buildings[0].GenerateRate*1 + game.buildings[1].GenerateRate*2
			Expect(game.GetTotalGenerateRate()).To(Equal(expectedRate))
		})

		It("should return 0 if no buildings are unlocked", func() {
			Expect(game.GetTotalGenerateRate()).To(Equal(0.0))
		})
	})

	Describe("handleInput", func() {
		It("should move the cursor up when KeyTypeUp is pressed", func() {
			game.cursor = 1 // Start at the second item
			mockInputHandler.pressedKey = KeyTypeUp

			game.handleInput()
			Expect(game.cursor).To(Equal(0)) // Cursor should move to the first item
		})

		It("should move the cursor down when KeyTypeDown is pressed", func() {
			game.cursor = 0 // Start at the first item
			mockInputHandler.pressedKey = KeyTypeDown

			game.handleInput()
			Expect(game.cursor).To(Equal(1)) // Cursor should move to the second item
		})

		It("should wrap the cursor to the bottom when moving up from the top", func() {
			game.cursor = 0 // Start at the first item
			mockInputHandler.pressedKey = KeyTypeUp

			game.handleInput()
			Expect(game.cursor).To(Equal(len(game.buildings))) // Cursor should wrap to the last item
		})

		It("should wrap the cursor to the top when moving down from the bottom", func() {
			game.cursor = len(game.buildings) // Start at the last item
			mockInputHandler.pressedKey = KeyTypeDown

			game.handleInput()
			Expect(game.cursor).To(Equal(0)) // Cursor should wrap to the first item
		})

		It("should trigger handleDecision when KeyTypeDecision is pressed", func() {
			mockInputHandler.pressedKey = KeyTypeDecision
			game.cursor = 0 // Select manual work

			game.handleInput()
			Expect(game.money).To(Equal(game.manualWork.Value)) // Money should increase
		})
	})
})
