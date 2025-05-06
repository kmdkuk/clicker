package components

import (
	"bytes"

	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/assets/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockListItem for testing
type MockListItem struct {
	StringValue string
}

func (m *MockListItem) String() string {
	return m.StringValue
}

var _ = Describe("List", func() {
	var (
		list       *List
		mockScreen *ebiten.Image
		items      []ListItem
	)
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.BebasNeueRegular_ttf))
	Expect(err).NotTo(HaveOccurred())

	BeforeEach(func() {
		mockScreen = ebiten.NewImage(640, 480)

		items = []ListItem{
			&MockListItem{StringValue: "Item 1"},
			&MockListItem{StringValue: "Item 2"},
			&MockListItem{StringValue: "Item 3"},
		}

		list = NewList(source, true, 10, 20)
		list.Items = items
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
			buildings := []dto.Building{
				{Name: "Building 1"},
				{Name: "Building 2"},
			}

			listItems := ConvertBuildingToListItems(buildings)
			Expect(listItems).To(HaveLen(2))
		})
	})

	Describe("ConvertUpgradeToListItems", func() {
		It("should convert a slice of upgrades to list items", func() {
			upgrades := []dto.Upgrade{
				{Name: "Upgrade 1"},
				{Name: "Upgrade 2"},
			}

			listItems := ConvertUpgradeToListItems(upgrades)
			Expect(listItems).To(HaveLen(2))
		})
	})
	Describe("List with view port", func() {
		var (
			list       *List
			mockScreen *ebiten.Image
		)

		BeforeEach(func() {
			// Create a test list with viewport size 3
			list = NewListWithViewport(source, true, 10, 20, 3)

			// Add test items
			list.Items = []ListItem{
				&MockListItem{StringValue: "Item 0"},
				&MockListItem{StringValue: "Item 1"},
				&MockListItem{StringValue: "Item 2"},
				&MockListItem{StringValue: "Item 3"},
				&MockListItem{StringValue: "Item 4"},
				&MockListItem{StringValue: "Item 5"},
			}

			mockScreen = ebiten.NewImage(640, 480)
		})

		Describe("Initialization", func() {
			It("should initialize with correct default values", func() {
				list := NewList(source, true, 10, 20)
				Expect(list.Visible).To(BeTrue())
				Expect(list.x).To(Equal(10))
				Expect(list.y).To(Equal(20))
				Expect(list.scrollPos).To(Equal(0))
			})

			It("should initialize with custom viewport size", func() {
				list := NewListWithViewport(source, false, 5, 15, 5)
				Expect(list.Visible).To(BeFalse())
				Expect(list.x).To(Equal(5))
				Expect(list.y).To(Equal(15))
				Expect(list.scrollPos).To(Equal(0))
				Expect(list.viewportSize).To(Equal(5)) // Custom viewport size
			})
		})

		Describe("Scrolling behavior", func() {
			Context("when cursor is outside viewport", func() {
				It("should scroll down when cursor is below viewport", func() {
					// Initial state
					Expect(list.scrollPos).To(Equal(0))

					// Move cursor below viewport
					cursor := 4 // This is beyond the viewport (items 0,1,2)
					list.Draw(mockScreen, cursor)

					// Scrolling should adjust to keep cursor in view
					Expect(list.scrollPos).To(BeNumerically(">", 0))
					start, end := list.GetVisibleRange()
					Expect(cursor).To(BeNumerically(">=", start))
					Expect(cursor).To(BeNumerically("<", end))
				})

				It("should scroll up when cursor is above viewport", func() {
					// Set initial scroll position down
					list.scrollPos = 3

					// Move cursor above viewport
					cursor := 1 // This is before the viewport (items 3,4,5)
					list.Draw(mockScreen, cursor)

					// Scrolling should adjust to show cursor
					Expect(list.scrollPos).To(BeNumerically("<=", cursor))
					start, _ := list.GetVisibleRange()
					Expect(cursor).To(BeNumerically(">=", start))
				})
			})

			Context("when using Scroll method", func() {
				It("should adjust scroll position within valid range", func() {
					// Scroll down
					list.Scroll(2)
					Expect(list.scrollPos).To(Equal(2))

					// Scroll more (but limit should be applied)
					list.Scroll(10)
					maxScroll := len(list.Items) - list.viewportSize
					Expect(list.scrollPos).To(Equal(maxScroll))

					// Scroll back up
					list.Scroll(-1)
					Expect(list.scrollPos).To(Equal(maxScroll - 1))

					// Try to scroll too far up
					list.Scroll(-10)
					Expect(list.scrollPos).To(Equal(0))
				})
			})

			Context("with empty list", func() {
				It("should handle empty list gracefully", func() {
					list.Items = []ListItem{}

					// Shouldn't panic when drawing empty list
					Expect(func() {
						list.Draw(mockScreen, 0)
					}).NotTo(Panic())

					// Scroll operations should be safe
					list.Scroll(5)
					Expect(list.scrollPos).To(Equal(0))

					start, end := list.GetVisibleRange()
					Expect(start).To(Equal(0))
					Expect(end).To(Equal(0))
				})
			})

			Context("when list is shorter than viewport", func() {
				It("should not allow scrolling with short list", func() {
					list.Items = []ListItem{
						&MockListItem{StringValue: "Item 0"},
						&MockListItem{StringValue: "Item 1"},
					}

					// Try to scroll
					list.Scroll(1)
					Expect(list.scrollPos).To(Equal(0))

					start, end := list.GetVisibleRange()
					Expect(start).To(Equal(0))
					Expect(end).To(Equal(2)) // All items visible
				})
			})
		})

		Describe("GetVisibleRange", func() {
			It("should return correct visible range", func() {
				// Default position
				start, end := list.GetVisibleRange()
				Expect(start).To(Equal(0))
				Expect(end).To(Equal(3)) // First 3 items

				// Scroll to middle
				list.scrollPos = 2
				start, end = list.GetVisibleRange()
				Expect(start).To(Equal(2))
				Expect(end).To(Equal(5)) // Items 2,3,4

				// Scroll near end
				list.scrollPos = 3
				start, end = list.GetVisibleRange()
				Expect(start).To(Equal(3))
				Expect(end).To(Equal(6)) // Items 3,4,5
			})

			It("should handle partial viewport at the end", func() {
				list.scrollPos = 4 // Only 2 items left from position 4
				start, end := list.GetVisibleRange()
				Expect(start).To(Equal(4))
				Expect(end).To(Equal(6)) // Items 4,5 (only 2 items)
			})
		})

		Describe("Draw method behavior", func() {
			It("should not draw when list is not visible", func() {
				list.Visible = false

				// Draw shouldn't modify scroll position when not visible
				initialScroll := list.scrollPos
				list.Draw(mockScreen, 4)
				Expect(list.scrollPos).To(Equal(initialScroll))
			})

			It("should draw scrollbar when list exceeds viewport", func() {
				// This is hard to test without accessing private methods
				// But we can at least verify it doesn't panic
				Expect(func() {
					list.Draw(mockScreen, 2)
				}).NotTo(Panic())
			})

			It("should highlight the item at cursor position", func() {
				// Again, hard to test UI rendering directly
				// But we can ensure the function runs without issues
				cursor := 2
				Expect(func() {
					list.Draw(mockScreen, cursor)
				}).NotTo(Panic())
			})
		})

		Describe("Integration with DTO items", func() {
			It("should handle Building items correctly", func() {
				buildings := []dto.Building{
					{Name: "Building 1", Cost: 100},
					{Name: "Building 2", Cost: 200},
				}

				listItems := ConvertBuildingToListItems(buildings)
				Expect(listItems).To(HaveLen(2))

				list.Items = listItems
				Expect(func() {
					list.Draw(mockScreen, 0)
				}).NotTo(Panic())
			})

			It("should handle Upgrade items correctly", func() {
				upgrades := []dto.Upgrade{
					{Name: "Upgrade 1", Cost: 100},
					{Name: "Upgrade 2", Cost: 200},
				}

				listItems := ConvertUpgradeToListItems(upgrades)
				Expect(listItems).To(HaveLen(2))

				list.Items = listItems
				Expect(func() {
					list.Draw(mockScreen, 0)
				}).NotTo(Panic())
			})
		})
	})

	Describe("GetHoverCursor", func() {
		var (
			list        *List
			screenWidth = 640
			x           = 10
			y           = 20
		)

		BeforeEach(func() {
			// Create a test list with viewport size 3
			list = NewListWithViewport(source, true, x, y, 3)

			// Add test items
			list.Items = []ListItem{
				&MockListItem{StringValue: "Item 0"},
				&MockListItem{StringValue: "Item 1"},
				&MockListItem{StringValue: "Item 2"},
				&MockListItem{StringValue: "Item 3"},
				&MockListItem{StringValue: "Item 4"},
			}

			mockScreen = ebiten.NewImage(screenWidth, 480)
		})

		It("should return the correct index when hovering over an item", func() {
			mouseX := screenWidth / 2  // Within the x range of the list
			mouseY := y + ItemHeight/2 // Within the y range of the first item

			cursor := list.GetHoverCursor(screenWidth, mouseX, mouseY)
			Expect(cursor).To(Equal(0)) // First item

			mouseY = y + ItemHeight + ItemHeight/2 // Within the y range of the second item
			cursor = list.GetHoverCursor(screenWidth, mouseX, mouseY)
			Expect(cursor).To(Equal(1)) // Second item
		})

		It("should return -1 when hovering outside the list bounds", func() {
			mouseX := 5                // Outside the x range of the list
			mouseY := y + ItemHeight/2 // Within the y range of the first item

			cursor := list.GetHoverCursor(screenWidth, mouseX, mouseY)
			Expect(cursor).To(Equal(-1))

			mouseY = y + ItemHeight + ItemHeight/2

			cursor = list.GetHoverCursor(screenWidth, mouseX, mouseY)
			Expect(cursor).To(Equal(-1))

			mouseX = screenWidth / 2  // Within the x range of the list
			mouseY = y - ItemHeight/2 // Outside the y range of the list

			cursor = list.GetHoverCursor(screenWidth, mouseX, mouseY)
			Expect(cursor).To(Equal(-1))
		})

		It("should return -1 when the list is not visible", func() {
			list.Visible = false

			mouseX := screenWidth / 2
			mouseY := y + ItemHeight/2

			cursor := list.GetHoverCursor(screenWidth, mouseX, mouseY)
			Expect(cursor).To(Equal(-1))
		})
	})
})
