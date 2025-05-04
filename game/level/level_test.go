package level

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const buildings_count = 10

var _ = Describe("Level", func() {
	Describe("NewBuildings", func() {
		It("correct buildings_*", func() {
			Expect(len(building_names)).To(Equal(buildings_count))
			Expect(len(building_base_costs)).To(Equal(buildings_count))
			Expect(len(building_base_generate_rates)).To(Equal(buildings_count))
			for i := 0; i < buildings_count-1; i++ {
				Expect(building_base_costs[i]).To(BeNumerically("<", building_base_costs[i+1]), "Buildings should have increasing base costs")
				Expect(building_base_generate_rates[i]).To(BeNumerically("<", building_base_generate_rates[i+1]), "Buildings should have increasing generate rates")
			}
		})

		It("should not have occure panic", func() {
			Expect(func() {
				NewBuildings()
			}).NotTo(Panic())
		})
	})
})
