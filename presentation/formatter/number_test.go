package formatter

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Number Formatter", func() {
	Context("FormatLargeNumber", func() {
		It("should handle small numbers", func() {
			Expect(FormatLargeNumber(0)).To(Equal("0.00"))
			Expect(FormatLargeNumber(5)).To(Equal("5.00"))
			Expect(FormatLargeNumber(42)).To(Equal("42.0"))
			Expect(FormatLargeNumber(999)).To(Equal("999"))
			Expect(FormatLargeNumber(3.14)).To(Equal("3.14"))
		})

		It("should format thousands (K)", func() {
			Expect(FormatLargeNumber(1000)).To(Equal("1.00K"))
			Expect(FormatLargeNumber(1500)).To(Equal("1.50K"))
			Expect(FormatLargeNumber(2750)).To(Equal("2.75K"))
			Expect(FormatLargeNumber(9999)).To(Equal("9.99K"))
			Expect(FormatLargeNumber(10000)).To(Equal("10.0K"))
			Expect(FormatLargeNumber(10500)).To(Equal("10.5K"))
			Expect(FormatLargeNumber(100000)).To(Equal("100K"))
			Expect(FormatLargeNumber(999999)).To(Equal("999K"))
		})

		It("should format millions (M)", func() {
			Expect(FormatLargeNumber(1000000)).To(Equal("1.00M"))
			Expect(FormatLargeNumber(1500000)).To(Equal("1.50M"))
			Expect(FormatLargeNumber(27500000)).To(Equal("27.5M"))
			Expect(FormatLargeNumber(999999999)).To(Equal("999M"))
		})

		It("should format billions (B)", func() {
			Expect(FormatLargeNumber(1000000000)).To(Equal("1.00B"))
			Expect(FormatLargeNumber(1500000000)).To(Equal("1.50B"))
			Expect(FormatLargeNumber(2750000000)).To(Equal("2.75B"))
		})

		It("should format trillions (T)", func() {
			Expect(FormatLargeNumber(1e12)).To(Equal("1.00T"))
			Expect(FormatLargeNumber(1.5e12)).To(Equal("1.50T"))
			Expect(FormatLargeNumber(2.75e12)).To(Equal("2.75T"))
		})

		It("should handle negative numbers", func() {
			Expect(FormatLargeNumber(-5)).To(Equal("-5.00"))
			Expect(FormatLargeNumber(-1500)).To(Equal("-1.50K"))
			Expect(FormatLargeNumber(-1e6)).To(Equal("-1.00M"))
		})

		It("should handle extremely large numbers", func() {
			Expect(FormatLargeNumber(1e30)).To(Equal("1.00e+30"))
			Expect(FormatLargeNumber(1e31)).To(Equal("1.00e+31"))
			Expect(FormatLargeNumber(1e32)).To(Equal("1.00e+32"))
		})
	})

	Context("FormatCurrency", func() {
		It("should add currency symbol to formatted numbers", func() {
			Expect(FormatCurrency(0, "$")).To(Equal("$ 0.00"))
			Expect(FormatCurrency(1500, "$")).To(Equal("$ 1.50K"))
			Expect(FormatCurrency(1e6, "$")).To(Equal("$ 1.00M"))
			Expect(FormatCurrency(1e6, "¥")).To(Equal("¥ 1.00M"))
		})
	})
})
