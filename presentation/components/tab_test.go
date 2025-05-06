package components

import (
	"bytes"

	"github.com/kmdkuk/clicker/assets/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tab Component", func() {
	var tab *Tab
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.BebasNeueRegular_ttf))
	Expect(err).NotTo(HaveOccurred())
	var (
		x = 10
		y = 20
	)

	BeforeEach(func() {
		// Initialize a tab with test data before each test
		tab = NewTab(source, []string{"Buildings", "Upgrades"}, 0, x, y)
	})

	Context("initialization", func() {
		It("should correctly initialize with provided values", func() {
			Expect(tab.titles).To(Equal([]string{"Buildings", "Upgrades"}))
			Expect(tab.activePage).To(Equal(0))
			Expect(tab.x).To(Equal(10))
			Expect(tab.y).To(Equal(20))
		})

		It("should support initialization with non-zero default page", func() {
			customTab := NewTab(source, []string{"Tab1", "Tab2", "Tab3"}, 1, 5, 15)
			Expect(customTab.activePage).To(Equal(1))
		})
	})

	Context("page management", func() {
		It("should get the active page correctly", func() {
			Expect(tab.GetActivePage()).To(Equal(0))

			tab.activePage = 1
			Expect(tab.GetActivePage()).To(Equal(1))
		})

		It("should set the active page within valid range", func() {
			tab.SetActivePage(1)
			Expect(tab.activePage).To(Equal(1))

			// Back to first page
			tab.SetActivePage(0)
			Expect(tab.activePage).To(Equal(0))
		})

		It("should ignore invalid page indices", func() {
			// Try setting out of range index
			tab.SetActivePage(2)                // Beyond the range of our test tabs
			Expect(tab.activePage).To(Equal(0)) // Should remain unchanged

			tab.SetActivePage(-1)
			Expect(tab.activePage).To(Equal(0)) // Should remain unchanged
		})
	})

	Context("drawing functionality", func() {
		It("should not panic when drawing", func() {
			// Just verify that Draw doesn't panic
			// We can't easily test the actual drawing results
			screen := ebiten.NewImage(320, 240)
			Expect(func() { tab.Draw(screen) }).NotTo(Panic())
		})

		// Note: Testing the exact visual output would require more complex testing
		// involving image comparison or mocking the ebitenutil.DebugPrintAt function
	})

	Context("with multiple tabs", func() {
		It("should handle multiple tabs correctly", func() {
			multiTab := NewTab(source, []string{"Tab1", "Tab2", "Tab3", "Tab4"}, 0, 10, 20)

			Expect(multiTab.titles).To(HaveLen(4))

			// Check we can navigate through all tabs
			for i := 0; i < len(multiTab.titles); i++ {
				multiTab.SetActivePage(i)
				Expect(multiTab.GetActivePage()).To(Equal(i))
			}
		})
	})

	Describe("GetHoverPage", func() {
		var (
			tab         *Tab
			screenWidth = 640
			x           = 10
			y           = 20
		)

		BeforeEach(func() {
			// Initialize a Tab instance for testing
			tab = NewTab(nil, []string{"Page 1", "Page 2", "Page 3"}, 0, x, y)
		})

		It("should return the correct page index when hovering over a tab", func() {
			tabWidth := (screenWidth - 2*x) / 3    // Assuming 3 tabs
			mouseX := 0*tabWidth + tabWidth/2 + 10 // Within the x range of the first tab
			mouseY := y + ItemHeight/2             // Within the y range of the tabs

			page := tab.GetHoverPage(screenWidth, mouseX, mouseY)
			Expect(page).To(Equal(0)) // First tab

			mouseX = 1*tabWidth + tabWidth/2 + 10 // Within the x range of the second tab
			page = tab.GetHoverPage(screenWidth, mouseX, mouseY)
			Expect(page).To(Equal(1)) // Second tab

			mouseX = 2*tabWidth + tabWidth/2 + 10 // Within the x range of the third tab
			page = tab.GetHoverPage(screenWidth, mouseX, mouseY)
			Expect(page).To(Equal(2)) // Third tab
		})

		It("should return -1 when hovering outside the tab bounds", func() {
			mouseX := 5  // Outside the x range of the tabs
			mouseY := 25 // Within the y range of the tabs

			page := tab.GetHoverPage(screenWidth, mouseX, mouseY)
			Expect(page).To(Equal(-1))

			mouseX = 15 // Within the x range of the tabs
			mouseY = 5  // Outside the y range of the tabs

			page = tab.GetHoverPage(screenWidth, mouseX, mouseY)
			Expect(page).To(Equal(-1))
		})
	})
})
