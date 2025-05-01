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
	})

	Describe("updateBuildings", func() {
		It("should generate income from unlocked buildings", func() {
			now := time.Now()
			game.buildings[0].count = 1                 // Unlock the first building
			game.lastUpdate = now.Add(-1 * time.Second) // Simulate 1 second elapsed

			game.updateBuildings(now)
			Expect(game.money).To(Equal(game.buildings[0].baseGenerateRate))
		})

		It("should not generate income from locked buildings", func() {
			now := time.Now()
			game.lastUpdate = now.Add(-1 * time.Second) // Simulate 1 second elapsed

			game.updateBuildings(now)
			Expect(game.money).To(Equal(0.0))
		})
	})

	Describe("handleDecision with page=0", func() {
		It("should add money for manual work", func() {
			game.cursor = 0 // Select manual work
			game.handleDecision()
			Expect(game.money).To(Equal(0.1))
		})

		It("should purchase a building if enough money is available", func() {
			game.cursor = 1        // Select the first building
			game.UpdateMoney(10.0) // Add enough money to purchase
			game.handleDecision()

			Expect(game.money).To(BeNumerically("<", 10.0)) // Money should decrease
			Expect(game.buildings[0].count).To(Equal(1))    // Building count should increase
		})

		It("should show a popup if not enough money is available", func() {
			game.cursor = 1 // Select the first building
			game.handleDecision()

			Expect(game.popup.Active).To(BeTrue()) // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money to unlock!"))
		})

		It("should show a popup if not enough money is available when unlocked", func() {
			game.cursor = 1             // Select the first building
			game.buildings[0].count = 1 // Unlock the building

			game.handleDecision()

			Expect(game.popup.Active).To(BeTrue()) // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money to purchase!"))
		})

		It("should correctly apply upgrades when performing manual work", func() {
			// Setup: Enable an upgrade that doubles manual work value
			game.upgrades = []Upgrade{
				{
					name:               "Double Manual Work",
					isTargetManualWork: true,
					isPurchased:        true,
					effect: func(value float64) float64 {
						return value * 2
					},
				},
			}
			game.cursor = 0 // Select manual work

			// Perform manual work
			game.handleDecision()

			// Expect the money to increase by the upgraded manual work value
			expectedMoney := 0.1 * 2
			Expect(game.money).To(Equal(expectedMoney))
		})
	})

	Describe("handleDecision with page=1 and cursor > 0", func() {
		BeforeEach(func() {
			game.page = 1 // Set to the second page
			game.upgrades = []Upgrade{
				{
					name:               "Test Upgrade 1",
					isPurchased:        false,
					isTargetManualWork: false,
					targetBuilding:     1,
					cost:               10.0,
					effect: func(value float64) float64 {
						return value * 2
					},
					isReleased: func(*Game) bool {
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
					isReleased: func(*Game) bool {
						return true
					},
				},
			}
		})

		It("should purchase an upgrade if enough money is available", func() {
			game.cursor = 1        // Select the first upgrade
			game.UpdateMoney(10.0) // Add enough money to purchase the upgrade

			game.handleDecision()

			Expect(game.money).To(BeNumerically("<", 10.0))                         // Money should decrease
			Expect(game.upgrades[0].isPurchased).To(BeTrue())                       // Upgrade should be marked as purchased
			Expect(game.popup.Active).To(BeTrue())                                  // Popup should be active
			Expect(game.popup.Message).To(Equal("Upgrade purchased successfully!")) // Correct popup message
		})

		It("should show a popup if not enough money is available for the upgrade", func() {
			game.cursor = 1       // Select the first upgrade
			game.UpdateMoney(5.0) // Not enough money to purchase the upgrade

			game.handleDecision()

			Expect(game.upgrades[0].isPurchased).To(BeFalse())                    // Upgrade should not be purchased
			Expect(game.popup.Active).To(BeTrue())                                // Popup should be active
			Expect(game.popup.Message).To(Equal("Not enough money for upgrade!")) // Correct popup message
		})

		It("should not allow purchasing an already purchased upgrade", func() {
			game.cursor = 1                     // Select the first upgrade
			game.UpdateMoney(10.0)              // Add enough money to purchase the upgrade
			game.upgrades[0].isPurchased = true // Mark the upgrade as already purchased

			game.handleDecision()

			Expect(game.money).To(Equal(10.0))                                 // Money should not decrease
			Expect(game.popup.Active).To(BeTrue())                             // Popup should be active
			Expect(game.popup.Message).To(Equal("Upgrade already purchased!")) // Correct popup message
		})

		It("should show a popup if the upgrade is not yet available", func() {
			game.page = 1   // Set to the second page
			game.cursor = 1 // Select the first upgrade
			game.upgrades = []Upgrade{
				{
					name:               "Test Upgrade 1",
					isTargetManualWork: false,
					isPurchased:        false,
					targetBuilding:     1,
					cost:               10.0,
					effect: func(value float64) float64 {
						return value * 2
					},
					isReleased: func(*Game) bool {
						return false // Upgrade is not yet available
					},
				},
			}

			game.UpdateMoney(10.0) // Add enough money to purchase the upgrade
			game.handleDecision()

			// Assert that the upgrade was not purchased
			Expect(game.upgrades[0].isPurchased).To(BeFalse())

			// Assert that the popup is active with the correct message
			Expect(game.popup.Active).To(BeTrue())
			Expect(game.popup.Message).To(Equal("Upgrade not available yet!"))
		})
	})

	Describe("GetTotalGenerateRate", func() {
		It("should calculate the total generate rate from all unlocked buildings", func() {
			game.buildings[0].count = 1
			game.buildings[1].count = 2

			expectedRate := game.buildings[0].baseGenerateRate*1 + game.buildings[1].baseGenerateRate*2
			Expect(game.GetTotalGenerateRate()).To(BeNumerically("~", expectedRate, 0.00001))
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
			Expect(game.money).To(Equal(game.manualWork.Value(game.upgrades))) // Money should increase
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
