package ui

import (
	"github.com/kmdkuk/clicker/input"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Popup", func() {
	var popup Popup

	BeforeEach(func() {
		popup = Popup{}
	})

	Describe("Show", func() {
		It("should activate the popup with the correct message", func() {
			message := "Test Message"
			popup.Show(message)

			Expect(popup.Active).To(BeTrue())
			Expect(popup.Message).To(Equal(message))
		})
	})

	Describe("Close", func() {
		It("should deactivate the popup", func() {
			popup.Active = true
			popup.Message = "Test Message"

			popup.Close()

			Expect(popup.Active).To(BeFalse())
			Expect(popup.Message).To(Equal("Test Message")) // Message should remain unchanged
		})
	})

	Describe("HandleInput", func() {
		It("should close the popup when KeyTypeDecision is pressed", func() {
			popup.Active = true
			popup.Message = "Test Message"

			popup.HandleInput(input.KeyTypeDecision)

			Expect(popup.Active).To(BeFalse())
		})

		It("should not close the popup for other key types", func() {
			popup.Active = true
			popup.Message = "Test Message"

			popup.HandleInput(input.KeyTypeUp)

			Expect(popup.Active).To(BeTrue())
			Expect(popup.Message).To(Equal("Test Message"))
		})
	})
})
