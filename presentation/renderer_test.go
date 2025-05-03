package presentation

import (
	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/presentation/input"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock for GameStateReader
type mockGameStateReader struct {
	money        float64
	totalGenRate float64
	manualWork   *model.ManualWork
	buildings    []model.Building
	upgrades     []model.Upgrade
}

// ManualWorkAction implements input.GameStateWriter.
func (m *mockGameStateReader) ManualWorkAction() {
}

// PurchaseBuildingAction implements input.GameStateWriter.
func (m *mockGameStateReader) PurchaseBuildingAction(cursor int) (bool, string) {
	return false, ""
}

// PurchaseUpgradeAction implements input.GameStateWriter.
func (m *mockGameStateReader) PurchaseUpgradeAction(cursor int) (bool, string) {
	return false, ""
}

func (m *mockGameStateReader) GetMoney() float64 {
	return m.money
}

func (m *mockGameStateReader) GetTotalGenerateRate() float64 {
	return m.totalGenRate
}

func (m *mockGameStateReader) GetManualWork() *model.ManualWork {
	return m.manualWork
}

func (m *mockGameStateReader) GetBuildings() []model.Building {
	return m.buildings
}

func (m *mockGameStateReader) GetUpgrades() []model.Upgrade {
	return m.upgrades
}

var _ = Describe("Renderer", func() {
	var (
		renderer   *DefaultRenderer
		gameState  *mockGameStateReader
		testConfig *config.Config
		mockScreen *ebiten.Image
	)

	BeforeEach(func() {
		// Setup test GameState
		gameState = &mockGameStateReader{
			money:        100.0,
			totalGenRate: 5.5,
			manualWork:   &model.ManualWork{Name: "Manual Work: $1.0"},
			buildings: []model.Building{
				model.Building{Name: "Building 1: $10"},
				model.Building{Name: "Building 2: $50"},
			},
			upgrades: []model.Upgrade{
				model.Upgrade{Name: "Upgrade 1: $20"},
				model.Upgrade{Name: "Upgrade 2: $100"},
			},
		}

		// Setup test Config
		testConfig = &config.Config{
			EnableDebug: true,
		}

		// Create Renderer
		renderer = NewRenderer(testConfig, gameState, input.NewDecider(gameState)).(*DefaultRenderer)

		// Create mock screen
		mockScreen = ebiten.NewImage(640, 480)
	})

	Describe("Initialization", func() {
		It("should be properly initialized", func() {
			Expect(renderer).NotTo(BeNil())
		})
	})

	Describe("Popup functionality", func() {
		Context("Popup display and state checks", func() {
			It("should have popup inactive in initial state", func() {
				Expect(renderer.IsPopupActive()).To(BeFalse())
			})

			It("should activate popup when ShowPopup is called", func() {
				renderer.ShowPopup("Test message")
				Expect(renderer.IsPopupActive()).To(BeTrue())
			})

			It("HandlePopup should return appropriate values based on popup state", func() {
				renderer.ShowPopup("Test message")
				// Should return true because popup is active
				renderer.HandleInput(input.KeyTypeNone)
				Expect(renderer.IsPopupActive()).To(BeTrue())

				// Close the popup
				renderer.HandleInput(input.KeyTypeDecision)
				Expect(renderer.IsPopupActive()).To(BeFalse())

				// Should return false because popup is inactive now
				renderer.HandleInput(input.KeyTypeNone)
				Expect(renderer.IsPopupActive()).To(BeFalse())
			})
		})
	})

	Describe("Popup input handling", func() {
		BeforeEach(func() {
			// Display popup before each test
			renderer.ShowPopup("Test message")
			Expect(renderer.IsPopupActive()).To(BeTrue())
		})

		It("should close popup when HandlePopupInput is called with decision key", func() {
			// Verify popup closes with decision key
			renderer.HandleInput(input.KeyTypeDecision)
			Expect(renderer.IsPopupActive()).To(BeFalse())
		})

		It("should not close popup when HandlePopupInput is called with non-decision keys", func() {
			// Verify popup doesn't close with non-decision keys
			renderer.HandleInput(input.KeyTypeUp)
			Expect(renderer.IsPopupActive()).To(BeTrue())

			renderer.HandleInput(input.KeyTypeDown)
			Expect(renderer.IsPopupActive()).To(BeTrue())
		})

	})

	Describe("Input handling with popup", func() {
		It("should handle popup closure through main input method", func() {
			// Display popup
			renderer.ShowPopup("Test message")
			Expect(renderer.IsPopupActive()).To(BeTrue())

			// HandleInput with decision key should close the popup
			renderer.HandleInput(input.KeyTypeDecision)
			Expect(renderer.IsPopupActive()).To(BeFalse())
		})
		It("should skip normal navigation when popup is active", func() {
			// Display popup
			renderer.ShowPopup("Test message")

			// Save initial state
			initialPage := renderer.navigation.GetPage()
			initialCursor := renderer.navigation.GetCursor()

			// Send navigation keys
			renderer.HandleInput(input.KeyTypeRight) // Try to change page
			renderer.HandleInput(input.KeyTypeDown)  // Try to move cursor

			// Verify navigation doesn't work when popup is active
			Expect(renderer.navigation.GetPage()).To(Equal(initialPage))
			Expect(renderer.navigation.GetCursor()).To(Equal(initialCursor))
		})

		It("should resume normal navigation after popup is closed", func() {
			// Display popup and then close it
			renderer.ShowPopup("Test message")
			renderer.HandleInput(input.KeyTypeDecision)

			// Save initial state
			initialCursor := renderer.navigation.GetCursor()

			// Verify cursor movement works
			renderer.HandleInput(input.KeyTypeDown)
			Expect(renderer.navigation.GetCursor()).NotTo(Equal(initialCursor))
		})
	})

	Describe("Drawing functionality", func() {
		It("should execute the drawing process without panicking", func() {
			// Verify that the drawing method executes without any runtime errors
			Expect(func() {
				renderer.Draw(mockScreen)
			}).NotTo(Panic())
		})

		It("should execute drawing process normally even when popup is active", func() {
			renderer.ShowPopup("Test popup")
			Expect(func() {
				renderer.Draw(mockScreen)
			}).NotTo(Panic())
		})
	})

	Describe("Navigation and cursor management", func() {
		It("should initialize cursor and page to default values", func() {
			// Verify initial states
			Expect(renderer.navigation.GetPage()).To(Equal(0))
			Expect(renderer.navigation.GetCursor()).To(Equal(0))
		})

		Context("when handling navigation inputs", func() {
			It("should move cursor up and down", func() {
				// Initial cursor position should be 0
				Expect(renderer.navigation.GetCursor()).To(Equal(0))

				// Move cursor down
				renderer.HandleInput(input.KeyTypeDown)
				Expect(renderer.navigation.GetCursor()).To(Equal(1))

				// Move cursor up
				renderer.HandleInput(input.KeyTypeUp)
				Expect(renderer.navigation.GetCursor()).To(Equal(0))
			})

			It("should wrap cursor when reaching boundaries", func() {
				// Move cursor up from top position (should wrap to bottom)
				renderer.HandleInput(input.KeyTypeUp)
				totalItems := len(gameState.GetBuildings()) + 1 // Manual work + buildings
				Expect(renderer.navigation.GetCursor()).To(Equal(totalItems - 1))
			})

			It("should change pages using left/right keys", func() {
				// Initial page should be 0
				Expect(renderer.navigation.GetPage()).To(Equal(0))

				// Move to next page
				renderer.HandleInput(input.KeyTypeRight)
				Expect(renderer.navigation.GetPage()).To(Equal(1))

				// Move back to previous page
				renderer.HandleInput(input.KeyTypeLeft)
				Expect(renderer.navigation.GetPage()).To(Equal(0))

				// Navigate left from first page should wrap to last page
				renderer.HandleInput(input.KeyTypeLeft)
				Expect(renderer.navigation.GetPage()).To(Equal(1)) // Assuming 2 pages total
			})

			It("should validate cursor position when switching pages", func() {
				// Move to page 1
				renderer.HandleInput(input.KeyTypeRight)
				Expect(renderer.navigation.GetPage()).To(Equal(1))

				// Set cursor to position that might be invalid on other pages
				for i := 0; i < 5; i++ {
					renderer.HandleInput(input.KeyTypeDown)
				}

				// Move back to page 0
				renderer.HandleInput(input.KeyTypeLeft)

				// Cursor should be validated within bounds of page 0
				Expect(renderer.navigation.GetCursor()).To(BeNumerically("<=", len(gameState.GetBuildings())))
			})
		})
	})

	Describe("Debug message functionality", func() {
		It("should set and retrieve debug messages", func() {
			testMessage := "Test debug message"
			renderer.DebugMessage(testMessage)
			Expect(renderer.GetDebugMessage()).To(Equal(testMessage))
		})

		Context("when debug mode is enabled", func() {
			BeforeEach(func() {
				testConfig.EnableDebug = true
			})

			It("should display debug messages", func() {
				// Since we can't directly check drawing, we just verify no panic
				renderer.DebugMessage("Debug test")
				Expect(func() {
					renderer.Draw(mockScreen)
				}).NotTo(Panic())
			})
		})

		Context("when debug mode is disabled", func() {
			BeforeEach(func() {
				testConfig.EnableDebug = false
			})

			It("should not display debug messages", func() {
				renderer.DebugMessage("Debug test")
				// Again, we just verify no panic
				Expect(func() {
					renderer.Draw(mockScreen)
				}).NotTo(Panic())
			})
		})
	})

	Describe("Decision handling", func() {
		It("should trigger decision action based on cursor and page", func() {
			// Move to a building position
			renderer.HandleInput(input.KeyTypeDown) // Move to building 1

			// Should not panic when making a decision
			Expect(func() {
				renderer.HandleInput(input.KeyTypeDecision)
			}).NotTo(Panic())
		})

		It("should show popup with message from decider", func() {
			// Mock the decider to return a specific message
			// Note: This would require exposing decider as a field or creating a test-specific method
			// For this test, we're assuming the decider works as expected

			// Instead we can verify the popup doesn't show with empty message
			renderer.HandleInput(input.KeyTypeDecision)
			Expect(renderer.IsPopupActive()).To(BeFalse()) // Our mock returns empty string
		})
	})
})
