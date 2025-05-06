package components

import (
	"github.com/kmdkuk/clicker/presentation/input"

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
		BeforeEach(func() {
			popup.Show("Test Message")
		})

		It("should close popup when decision key is pressed", func() {
			popup.HandleInput(input.KeyTypeDecision, false)
			Expect(popup.Active).To(BeFalse())
		})

		It("should not close popup when non-decision keys are pressed", func() {
			popup.HandleInput(input.KeyTypeUp, false)
			Expect(popup.Active).To(BeTrue())

			popup.HandleInput(input.KeyTypeDown, false)
			Expect(popup.Active).To(BeTrue())

			popup.HandleInput(input.KeyTypeLeft, false)
			Expect(popup.Active).To(BeTrue())

			popup.HandleInput(input.KeyTypeRight, false)
			Expect(popup.Active).To(BeTrue())
		})

		It("should not process input when popup is inactive", func() {
			popup.Close()
			// Don't process key input when popup is inactive
			popup.HandleInput(input.KeyTypeDecision, false)
			Expect(popup.Active).To(BeFalse())
		})
	})

	Describe("IsActive", func() {
		It("should return correct active state", func() {
			popup.Active = false
			Expect(popup.IsActive()).To(BeFalse())

			popup.Active = true
			Expect(popup.IsActive()).To(BeTrue())
		})
	})

	Describe("GetMessage", func() {
		It("should return the current message", func() {
			testMessage := "Test GetMessage"
			popup.Message = testMessage
			Expect(popup.GetMessage()).To(Equal(testMessage))
		})
	})
})
