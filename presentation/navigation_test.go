package presentation

import (
	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/presentation/input"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// GameStateReaderMock is a test mock for model.GameStateReader
type GameStateReaderMock struct {
	buildings    []model.Building
	upgrades     []model.Upgrade
	manualWork   model.ManualWork
	money        float64
	totalGenRate float64
}

func (g *GameStateReaderMock) GetMoney() float64 {
	return g.money
}

func (g *GameStateReaderMock) GetTotalGenerateRate() float64 {
	return g.totalGenRate
}

func (g *GameStateReaderMock) GetManualWork() *model.ManualWork {
	return &g.manualWork
}

func (g *GameStateReaderMock) GetBuildings() []model.Building {
	return g.buildings
}

func (g *GameStateReaderMock) GetUpgrades() []model.Upgrade {
	return g.upgrades
}

var _ = Describe("Navigation", func() {
	var (
		nav           *Navigation
		gameStateMock *GameStateReaderMock
	)

	BeforeEach(func() {
		// Setup test game state
		gameStateMock = &GameStateReaderMock{
			buildings: []model.Building{
				{Name: "Building 1"},
				{Name: "Building 2"},
				{Name: "Building 3"},
			},
			upgrades: []model.Upgrade{
				{Name: "Upgrade 1"},
				{Name: "Upgrade 2"},
			},
			manualWork: model.ManualWork{},
		}

		// Initialize navigation
		nav = NewNavigation([]int{len(gameStateMock.GetBuildings()), len(gameStateMock.GetUpgrades())})
	})

	Describe("Initial state", func() {
		It("should have initial cursor position and page number set to 0", func() {
			Expect(nav.GetCursor()).To(Equal(0))
			Expect(nav.GetPage()).To(Equal(0))
		})
	})

	Describe("Cursor movement", func() {
		Context("with up/down keys", func() {
			It("should move to next item with down key", func() {
				nav.HandleNavigation(input.KeyTypeDown)
				Expect(nav.GetCursor()).To(Equal(1))

				nav.HandleNavigation(input.KeyTypeDown)
				Expect(nav.GetCursor()).To(Equal(2))
			})

			It("should move to previous item with up key", func() {
				// First move down
				nav.HandleNavigation(input.KeyTypeDown)
				nav.HandleNavigation(input.KeyTypeDown)
				Expect(nav.GetCursor()).To(Equal(2))

				// Then move back up
				nav.HandleNavigation(input.KeyTypeUp)
				Expect(nav.GetCursor()).To(Equal(1))
			})

			It("should wrap to last item when pressing up key from first position", func() {
				// Start at first position
				Expect(nav.GetCursor()).To(Equal(0))

				// Press up to wrap to the last item
				nav.HandleNavigation(input.KeyTypeUp)
				// Buildings + ManualWork = 4 items total
				Expect(nav.GetCursor()).To(Equal(3))
			})

			It("should wrap to first item when pressing down key from last position", func() {
				// Move to last position
				nav.HandleNavigation(input.KeyTypeUp) // Wrap to the end
				Expect(nav.GetCursor()).To(Equal(3))

				// Press down to wrap back to the start
				nav.HandleNavigation(input.KeyTypeDown)
				Expect(nav.GetCursor()).To(Equal(0))
			})
		})
	})

	Describe("Page navigation", func() {
		It("should move to next page with right key", func() {
			Expect(nav.GetPage()).To(Equal(0))

			nav.HandleNavigation(input.KeyTypeRight)
			Expect(nav.GetPage()).To(Equal(1))
		})

		It("should move to previous page with left key", func() {
			// First move right
			nav.HandleNavigation(input.KeyTypeRight)
			Expect(nav.GetPage()).To(Equal(1))

			// Then move back left
			nav.HandleNavigation(input.KeyTypeLeft)
			Expect(nav.GetPage()).To(Equal(0))
		})

		It("should wrap to last page when pressing left key from first page", func() {
			Expect(nav.GetPage()).To(Equal(0))

			nav.HandleNavigation(input.KeyTypeLeft)
			Expect(nav.GetPage()).To(Equal(1)) // 2 pages total
		})

		It("should wrap to first page when pressing right key from last page", func() {
			nav.HandleNavigation(input.KeyTypeRight) // Move to page 2
			Expect(nav.GetPage()).To(Equal(1))

			nav.HandleNavigation(input.KeyTypeRight) // Wrap to page 1
			Expect(nav.GetPage()).To(Equal(0))
		})
	})

	Describe("Cursor position validation", func() {
		It("should adjust cursor position to valid range when changing pages", func() {
			// Move cursor to the last item
			nav.HandleNavigation(input.KeyTypeUp) // Use wrapping to go to the end
			Expect(nav.GetCursor()).To(Equal(3))

			// Change page (upgrades + manualWork = 3 items total)
			nav.HandleNavigation(input.KeyTypeRight)

			// Cursor should be adjusted to valid range
			Expect(nav.GetCursor()).To(BeNumerically("<=", 2))
		})

		It("should preserve cursor position when moving from page with fewer items to page with more items", func() {
			// Move to Upgrades page (fewer items)
			nav.HandleNavigation(input.KeyTypeRight)

			// Move cursor to position 2 (valid in this page)
			nav.HandleNavigation(input.KeyTypeDown)
			nav.HandleNavigation(input.KeyTypeDown)
			Expect(nav.GetCursor()).To(Equal(2))

			// Go back to Buildings page
			nav.HandleNavigation(input.KeyTypeLeft)

			// Cursor position should be preserved
			Expect(nav.GetCursor()).To(Equal(2))
		})
	})

	Describe("Other input keys", func() {
		It("should ignore non-navigation keys", func() {
			// Record initial state
			initialCursor := nav.GetCursor()
			initialPage := nav.GetPage()

			// Send a decision key (not handled by navigation)
			nav.HandleNavigation(input.KeyTypeDecision)

			// Position should remain unchanged
			Expect(nav.GetCursor()).To(Equal(initialCursor))
			Expect(nav.GetPage()).To(Equal(initialPage))
		})
	})
})
