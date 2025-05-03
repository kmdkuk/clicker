package components

import (
	"github.com/kmdkuk/clicker/domain/model"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockListItem for testing
type MockListItem struct {
	StringValue string
}

func (m *MockListItem) String(gameState model.GameStateReader) string {
	return m.StringValue
}

var _ = Describe("List", func() {
	var (
		list       *List
		mockScreen *ebiten.Image
		gameState  *GameStateReaderMock
		items      []ListItem
	)

	BeforeEach(func() {
		mockScreen = ebiten.NewImage(640, 480)
		gameState = &GameStateReaderMock{}

		items = []ListItem{
			&MockListItem{StringValue: "Item 1"},
			&MockListItem{StringValue: "Item 2"},
			&MockListItem{StringValue: "Item 3"},
		}

		list = NewList(gameState, items, true, 10, 20)
	})

	Describe("NewList", func() {
		It("should initialize a list with the provided parameters", func() {
			Expect(list.Items).To(HaveLen(3))
			Expect(list.Visible).To(BeTrue())
		})
	})

	Describe("Draw", func() {
		It("should not panic when drawing", func() {
			Expect(func() {
				list.Draw(mockScreen, 0)
			}).NotTo(Panic())
		})

		It("should not draw when not visible", func() {
			list.Visible = false
			Expect(func() {
				list.Draw(mockScreen, 0)
			}).NotTo(Panic())
		})

		It("should handle cursor out of range", func() {
			Expect(func() {
				list.Draw(mockScreen, 10) // Out of range
			}).NotTo(Panic())
		})
	})

	Describe("ConvertBuildingToListItems", func() {
		It("should convert a slice of buildings to list items", func() {
			buildings := []model.Building{
				{Name: "Building 1"},
				{Name: "Building 2"},
			}

			listItems := ConvertBuildingToListItems(buildings)
			Expect(listItems).To(HaveLen(2))
		})
	})

	Describe("ConvertUpgradeToListItems", func() {
		It("should convert a slice of upgrades to list items", func() {
			upgrades := []model.Upgrade{
				{Name: "Upgrade 1"},
				{Name: "Upgrade 2"},
			}

			listItems := ConvertUpgradeToListItems(upgrades)
			Expect(listItems).To(HaveLen(2))
		})
	})
})
