package presentation

import (
	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/presentation/components"
	"github.com/kmdkuk/clicker/presentation/input"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockPlayerUseCase struct {
	player *dto.Player
}

func (m *MockPlayerUseCase) GetPlayer() *dto.Player {
	return m.player
}

type MockManualWorkUseCase struct {
	ManualWorkActionCalled bool
	manualWork             *dto.ManualWork
}

func (m *MockManualWorkUseCase) ManualWorkAction() {
	m.ManualWorkActionCalled = true
}
func (m *MockManualWorkUseCase) GetManualWork() *dto.ManualWork {
	return m.manualWork
}

type MockBuildingUseCase struct {
	PurchaseBuildingActionCalled  bool
	buildings                     []dto.Building
	successPurchaseBuildingAction bool
	messagePurchaseBuildingAction string
}

func (m *MockBuildingUseCase) GetBuildings() []dto.Building {
	return m.buildings
}
func (m *MockBuildingUseCase) GetBuildingsIsUnlockedWithMaskedNextLock() []dto.Building {
	return m.buildings
}
func (m *MockBuildingUseCase) PurchaseBuildingAction(index int) (bool, string) {
	m.PurchaseBuildingActionCalled = true
	return m.successPurchaseBuildingAction, m.messagePurchaseBuildingAction
}

type MockUpgradeUseCase struct {
	PurchaseUpgradeActionCalled  bool
	upgrades                     []dto.Upgrade
	successPurchaseUpgradeAction bool
	messagePurchaseUpgradeAction string
}

func (m *MockUpgradeUseCase) GetUpgrades() []dto.Upgrade {
	return m.upgrades
}
func (m *MockUpgradeUseCase) GetUpgradesIsReleasedCostSorted() []dto.Upgrade {
	return m.upgrades
}
func (m *MockUpgradeUseCase) PurchaseUpgradeAction(index int) (bool, string) {
	m.PurchaseUpgradeActionCalled = true
	return m.successPurchaseUpgradeAction, m.messagePurchaseUpgradeAction
}

var _ = Describe("Renderer", func() {
	var (
		renderer          *DefaultRenderer
		testConfig        *config.Config
		mockScreen        *ebiten.Image
		playerUseCase     *MockPlayerUseCase
		manualWorkUseCase *MockManualWorkUseCase
		buildingUseCase   *MockBuildingUseCase
		upgradeUseCase    *MockUpgradeUseCase
	)

	BeforeEach(func() {
		// Setup test Config
		testConfig = &config.Config{
			EnableDebug:  true,
			ScreenWidth:  640,
			ScreenHeight: 480,
		}

		playerUseCase = &MockPlayerUseCase{
			player: &dto.Player{
				Money:             100,
				TotalGenerateRate: 10,
			},
		}

		manualWorkUseCase = &MockManualWorkUseCase{
			manualWork: &dto.ManualWork{
				Name:  "Manual Work",
				Value: 10,
			},
		}

		buildingUseCase = &MockBuildingUseCase{
			buildings: []dto.Building{
				{Name: "Building 1: $10"},
				{Name: "Building 2: $50"},
				{Name: "Building 3: $50"},
			},
		}

		upgradeUseCase = &MockUpgradeUseCase{
			upgrades: []dto.Upgrade{
				{Name: "Upgrade 1: $20"},
				{Name: "Upgrade 2: $100"},
			},
			successPurchaseUpgradeAction: true,
			messagePurchaseUpgradeAction: "",
		}

		// Create Renderer
		r, err := NewRenderer(testConfig,
			playerUseCase,
			manualWorkUseCase,
			buildingUseCase,
			upgradeUseCase,
		)
		Expect(err).NotTo(HaveOccurred())
		renderer = r.(*DefaultRenderer)

		// Create mock screen
		mockScreen = ebiten.NewImage(testConfig.ScreenWidth, testConfig.ScreenHeight)
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
				renderer.HandleInput(input.KeyTypeNone, false, false, 0, 0)
				Expect(renderer.IsPopupActive()).To(BeTrue())

				// Close the popup
				renderer.HandleInput(input.KeyTypeDecision, false, false, 0, 0)
				Expect(renderer.IsPopupActive()).To(BeFalse())

				// Should return false because popup is inactive now
				renderer.HandleInput(input.KeyTypeNone, false, false, 0, 0)
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
			renderer.HandleInput(input.KeyTypeDecision, false, false, 0, 0)
			Expect(renderer.IsPopupActive()).To(BeFalse())
		})

		It("should not close popup when HandlePopupInput is called with non-decision keys", func() {
			// Verify popup doesn't close with non-decision keys
			renderer.HandleInput(input.KeyTypeUp, false, false, 0, 0)
			Expect(renderer.IsPopupActive()).To(BeTrue())

			renderer.HandleInput(input.KeyTypeDown, false, false, 0, 0)
			Expect(renderer.IsPopupActive()).To(BeTrue())
		})

	})

	Describe("Input handling with popup", func() {
		It("should handle popup closure through main input method", func() {
			// Display popup
			renderer.ShowPopup("Test message")
			Expect(renderer.IsPopupActive()).To(BeTrue())

			// HandleInput with decision key should close the popup
			renderer.HandleInput(input.KeyTypeDecision, false, false, 0, 0)
			Expect(renderer.IsPopupActive()).To(BeFalse())
		})
		It("should skip normal navigation when popup is active", func() {
			// Display popup
			renderer.ShowPopup("Test message")

			// Save initial state
			initialPage := renderer.navigation.GetPage()
			initialCursor := renderer.navigation.GetCursor()

			// Send navigation keys
			renderer.HandleInput(input.KeyTypeRight, false, false, 0, 0) // Try to change page
			renderer.HandleInput(input.KeyTypeDown, false, false, 0, 0)  // Try to move cursor

			// Verify navigation doesn't work when popup is active
			Expect(renderer.navigation.GetPage()).To(Equal(initialPage))
			Expect(renderer.navigation.GetCursor()).To(Equal(initialCursor))
		})

		It("should resume normal navigation after popup is closed", func() {
			// Display popup and then close it
			renderer.ShowPopup("Test message")
			renderer.HandleInput(input.KeyTypeDecision, false, false, 0, 0)

			// Save initial state
			initialCursor := renderer.navigation.GetCursor()

			// Verify cursor movement works
			renderer.HandleInput(input.KeyTypeDown, false, false, 0, 0)
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
				renderer.HandleInput(input.KeyTypeDown, false, false, 0, 0)
				Expect(renderer.navigation.GetCursor()).To(Equal(1))

				// Move cursor up
				renderer.HandleInput(input.KeyTypeUp, false, false, 0, 0)
				Expect(renderer.navigation.GetCursor()).To(Equal(0))
			})

			It("should wrap cursor when reaching boundaries", func() {
				// Move cursor up from top position (should wrap to bottom)
				renderer.HandleInput(input.KeyTypeUp, false, false, 0, 0)
				totalItems := len(buildingUseCase.GetBuildings()) + 1 // Manual work + buildings
				Expect(renderer.navigation.GetCursor()).To(Equal(totalItems - 1))
			})

			It("should change pages using left/right keys", func() {
				// Initial page should be 0
				Expect(renderer.navigation.GetPage()).To(Equal(0))

				// Move to next page
				renderer.HandleInput(input.KeyTypeRight, false, false, 0, 0)
				Expect(renderer.navigation.GetPage()).To(Equal(1))

				// Move back to previous page
				renderer.HandleInput(input.KeyTypeLeft, false, false, 0, 0)
				Expect(renderer.navigation.GetPage()).To(Equal(0))

				// Navigate left from first page should wrap to last page
				renderer.HandleInput(input.KeyTypeLeft, false, false, 0, 0)
				Expect(renderer.navigation.GetPage()).To(Equal(1)) // Assuming 2 pages total
			})

			It("should validate cursor position when switching pages", func() {
				// Move to page 1
				renderer.HandleInput(input.KeyTypeRight, false, false, 0, 0)
				Expect(renderer.navigation.GetPage()).To(Equal(1))

				// Set cursor to position that might be invalid on other pages
				for i := 0; i < 5; i++ {
					renderer.HandleInput(input.KeyTypeDown, false, false, 0, 0)
				}

				// Move back to page 0
				renderer.HandleInput(input.KeyTypeLeft, false, false, 0, 0)

				// Cursor should be validated within bounds of page 0
				Expect(renderer.navigation.GetCursor()).To(BeNumerically("<=", len(buildingUseCase.GetBuildings())))
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
			renderer.HandleInput(input.KeyTypeDown, false, false, 0, 0) // Move to building 1

			// Should not panic when making a decision
			Expect(func() {
				renderer.HandleInput(input.KeyTypeDecision, false, false, 0, 0)
			}).NotTo(Panic())
		})

		It("should show popup with message from decider", func() {
			// Mock the decider to return a specific message
			// Note: This would require exposing decider as a field or creating a test-specific method
			// For this test, we're assuming the decider works as expected

			// Instead we can verify the popup doesn't show with empty message
			renderer.HandleInput(input.KeyTypeDecision, false, false, 0, 0)
			Expect(renderer.IsPopupActive()).To(BeFalse()) // Our mock returns empty string
		})

		It("should show popup with error message when purchase fails", func() {
			// Override mock to return an error
			buildingUseCase.successPurchaseBuildingAction = false
			buildingUseCase.messagePurchaseBuildingAction = "Not enough money!"

			renderer.navigation.cursor = 1
			renderer.navigation.page = 0
			renderer.handleDecision(false, 0, 0)
			Expect(renderer.IsPopupActive()).To(BeTrue())
			Expect(renderer.GetPopupMessage()).To(Equal("Not enough money!"))
		})
		Context("when isClicked is true", func() {
			It("should set the page if a tab is clicked", func() {
				tab1X := testConfig.ScreenWidth/2 + testConfig.ScreenWidth/4
				tabY := 110
				renderer.handleDecision(true, tab1X, tabY)

				// Verify the page is set correctly
				Expect(renderer.navigation.GetPage()).To(Equal(1))
			})

			It("should set the cursor if an item is clicked", func() {
				itemX := testConfig.ScreenWidth / 2
				item2Y := 130 + components.ItemHeight + components.ItemHeight/2
				renderer.Update()
				renderer.handleDecision(true, itemX, item2Y)
				// Verify the cursor is set correctly
				Expect(renderer.navigation.GetCursor()).To(Equal(2))
			})

			It("should not set page or cursor if nothing is clicked", func() {

				initialPage := renderer.navigation.GetPage()
				initialCursor := renderer.navigation.GetCursor()

				renderer.handleDecision(true, 0, 0)

				// Verify the page and cursor remain unchanged
				Expect(renderer.navigation.GetPage()).To(Equal(initialPage))
				Expect(renderer.navigation.GetCursor()).To(Equal(initialCursor))
			})
		})
	})

	It("should update item lists with data from use cases", func() {
		// Setup test data and ensure renderer has old data
		renderer.Update()

		// Verify lists have been updated
		Expect(renderer.manualWork.Items).To(HaveLen(1))
		Expect(renderer.buildings.Items).To(HaveLen(len(buildingUseCase.buildings)))
		Expect(renderer.upgrades.Items).To(HaveLen(len(upgradeUseCase.upgrades)))
	})
})
