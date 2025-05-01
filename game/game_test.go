package game

import (
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

// MockDecisionProcessor はテスト用の決定プロセッサ
type MockDecider struct {
	success bool
	message string
	called  bool
	page    int
	cursor  int
}

func (mdp *MockDecider) Decide(page, cursor int) (bool, string) {
	mdp.called = true
	mdp.page = page
	mdp.cursor = cursor
	return mdp.success, mdp.message
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

	Describe("handleDecision", func() {
		It("should delegate decision processing to DecisionProcessor", func() {
			mockProcessor := &MockDecider{
				success: true,
				message: "Test message",
			}

			game.decider = mockProcessor
			game.page = 1
			game.cursor = 2

			game.handleDecision()

			Expect(mockProcessor.called).To(BeTrue())
			Expect(mockProcessor.page).To(Equal(1))
			Expect(mockProcessor.cursor).To(Equal(2))
			Expect(game.popup.Active).To(BeTrue())
			Expect(game.popup.Message).To(Equal("Test message"))
		})
	})

	Describe("handleDecision with page=0", func() {
		It("should add money for manual work", func() {
			game.cursor = 0 // Select manual work
			game.handleDecision()
			Expect(game.gameState.GetMoney()).To(Equal(0.1))
		})

		It("should purchase a building if enough money is available", func() {
			game.cursor = 1                  // Select the first building
			game.gameState.UpdateMoney(10.0) // Add enough money to purchase
			game.handleDecision()

			Expect(game.gameState.GetMoney()).To(BeNumerically("<", 10.0)) // Money should decrease
			Expect(game.gameState.GetBuildings()[0].count).To(Equal(1))    // Building count should increase
		})

		It("should show a popup if not enough money is available", func() {
			game.cursor = 1 // Select the first building
			game.handleDecision()

			Expect(game.popup.Active).To(BeTrue()) // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money to unlock!"))
		})

		It("should show a popup if not enough money is available when unlocked", func() {
			game.cursor = 1                            // Select the first building
			game.gameState.GetBuildings()[0].count = 1 // Unlock the building

			game.handleDecision()

			Expect(game.popup.Active).To(BeTrue()) // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money to purchase!"))
		})

		It("should correctly apply upgrades when performing manual work", func() {
			// Setup: Enable an upgrade that doubles manual work value
			game.gameState.SetUpgrades([]Upgrade{
				{
					name:               "Double Manual Work",
					isTargetManualWork: true,
					isPurchased:        true,
					effect: func(value float64) float64 {
						return value * 2
					},
				},
			})
			game.cursor = 0 // Select manual work

			// Perform manual work
			game.handleDecision()

			// Expect the money to increase by the upgraded manual work value
			expectedMoney := 0.1 * 2
			Expect(game.gameState.GetMoney()).To(Equal(expectedMoney))
		})
	})

	Describe("handleDecision with page=1 and cursor > 0", func() {
		BeforeEach(func() {
			game.page = 1 // Set to the second page
			game.gameState.SetUpgrades([]Upgrade{
				{
					name:               "Test Upgrade 1",
					isPurchased:        false,
					isTargetManualWork: false,
					targetBuilding:     1,
					cost:               10.0,
					effect: func(value float64) float64 {
						return value * 2
					},
					isReleased: func(g GameState) bool {
						return true
					},
				},
				{
					name:               "Test Upgrade 2",
					isPurchased:        false,
					isTargetManualWork: false,
					targetBuilding:     1,
					cost:               20.0,
					effect: func(value float64) float64 {
						return value + 5
					},
					isReleased: func(g GameState) bool {
						return true
					},
				},
			})
		})

		It("should purchase an upgrade if enough money is available", func() {
			game.cursor = 1                  // Select the first upgrade
			game.gameState.UpdateMoney(10.0) // Add enough money to purchase the upgrade

			game.handleDecision()

			Expect(game.gameState.GetMoney()).To(BeNumerically("<", 10.0))          // Money should decrease
			Expect(game.gameState.GetUpgrades()[0].isPurchased).To(BeTrue())        // Upgrade should be marked as purchased
			Expect(game.popup.Active).To(BeTrue())                                  // Popup should be active
			Expect(game.popup.Message).To(Equal("Upgrade purchased successfully!")) // Correct popup message
		})

		It("should show a popup if not enough money is available for the upgrade", func() {
			game.cursor = 1                 // Select the first upgrade
			game.gameState.UpdateMoney(5.0) // Not enough money to purchase the upgrade

			game.handleDecision()

			Expect(game.gameState.GetUpgrades()[0].isPurchased).To(BeFalse())     // Upgrade should not be purchased
			Expect(game.popup.Active).To(BeTrue())                                // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money for upgrade!")) // Correct popup message
		})

		It("should not allow purchasing an already purchased upgrade", func() {
			game.cursor = 1                                    // Select the first upgrade
			game.gameState.UpdateMoney(10.0)                   // Add enough money to purchase the upgrade
			game.gameState.GetUpgrades()[0].isPurchased = true // Mark the upgrade as already purchased

			game.handleDecision()

			Expect(game.gameState.GetMoney()).To(Equal(10.0))                  // Money should not decrease
			Expect(game.popup.Active).To(BeTrue())                             // Popup should be active
			Expect(game.popup.Message).To(Equal("Upgrade already purchased!")) // Correct popup message
		})

		It("should show a popup if the upgrade is not yet available", func() {
			game.page = 1   // Set to the second page
			game.cursor = 1 // Select the first upgrade
			game.gameState.SetUpgrades([]Upgrade{
				{
					name:               "Test Upgrade 1",
					isTargetManualWork: false,
					isPurchased:        false,
					targetBuilding:     1,
					cost:               10.0,
					effect: func(value float64) float64 {
						return value * 2
					},
					isReleased: func(g GameState) bool {
						return false // Upgrade is not yet available
					},
				},
			})

			game.gameState.UpdateMoney(10.0) // Add enough money to purchase the upgrade
			game.handleDecision()

			// Assert that the upgrade was not purchased
			Expect(game.gameState.GetUpgrades()[0].isPurchased).To(BeFalse())

			// Assert that the popup is active with the correct message
			Expect(game.popup.Active).To(BeTrue())
			Expect(game.popup.Message).To(Equal("Upgrade not available yet!"))
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
			Expect(game.cursor).To(Equal(len(game.gameState.GetBuildings()))) // Cursor should wrap to the last item
		})

		It("should wrap the cursor to the top when moving down from the bottom", func() {
			game.cursor = len(game.gameState.GetBuildings()) // Start at the last item
			mockInputHandler.pressedKey = KeyTypeDown

			game.handleInput()
			Expect(game.cursor).To(Equal(0)) // Cursor should wrap to the first item
		})

		It("should trigger handleDecision when KeyTypeDecision is pressed", func() {
			mockInputHandler.pressedKey = KeyTypeDecision
			game.cursor = 0 // Select manual work

			game.handleInput()
			Expect(game.gameState.GetMoney()).To(BeNumerically("~", 0.1, 0.0001)) // Money should increase
		})

		It("should move to the next page when KeyTypeRight is pressed", func() {
			game.page = 0 // Start on the first page
			mockInputHandler.pressedKey = KeyTypeRight

			game.handleInput()
			Expect(game.page).To(Equal(1)) // Page should move to the second page
		})

		It("should wrap to the first page when KeyTypeRight is pressed on the last page", func() {
			game.page = 1 // Start on the last page
			mockInputHandler.pressedKey = KeyTypeRight

			game.handleInput()
			Expect(game.page).To(Equal(0)) // Page should wrap to the first page
		})

		It("should move to the previous page when KeyTypeLeft is pressed", func() {
			game.page = 1 // Start on the second page
			mockInputHandler.pressedKey = KeyTypeLeft

			game.handleInput()
			Expect(game.page).To(Equal(0)) // Page should move to the first page
		})

		It("should wrap to the last page when KeyTypeLeft is pressed on the first page", func() {
			game.page = 0 // Start on the first page
			mockInputHandler.pressedKey = KeyTypeLeft

			game.handleInput()
			Expect(game.page).To(Equal(1)) // Page should wrap to the last page
		})
	})
})
