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

	BeforeEach(func() {
		// Initialize a tab with test data before each test
		tab = NewTab(source, []string{"Buildings", "Upgrades"}, 0, 10, 20)
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
})
