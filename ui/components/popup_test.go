package components

import (
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
})
