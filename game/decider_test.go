package game

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockGameState struct {
	money                        float64
	manualWork                   ManualWork
	buildings                    []Building
	upgrades                     []Upgrade
	manualWorkCalled             bool
	updateBuildingsCalled        bool
	getTotalGenerateRateCalled   bool
	purchaseBuildingActionCalled bool
	purchaseUpgradeActionCalled  bool
}

func (m *MockGameState) UpdateMoney(amount float64) {
	m.money += amount
}

func (m *MockGameState) GetMoney() float64 {
	return m.money
}
func (m *MockGameState) GetManualWork() *ManualWork {
	return &m.manualWork
}
func (m *MockGameState) SetManualWork(manualWork ManualWork) {
	m.manualWork = manualWork
}
func (m *MockGameState) GetBuildings() []Building {
	return m.buildings
}
func (m *MockGameState) GetUpgrades() []Upgrade {
	return m.upgrades
}
func (m *MockGameState) SetUpgrades(upgrades []Upgrade) {
	m.upgrades = upgrades
}
func (m *MockGameState) ManualWork() {
	m.manualWorkCalled = true
	m.UpdateMoney(m.manualWork.Work(m.upgrades))
}
func (m *MockGameState) UpdateBuildings(now time.Time) {
	m.updateBuildingsCalled = true
}
func (m *MockGameState) GetTotalGenerateRate() float64 {
	m.getTotalGenerateRateCalled = true
	return 0.0
}

func (m *MockGameState) PurchaseBuildingAction(buildingIndex int) (bool, string) {
	m.purchaseBuildingActionCalled = true
	return true, ""
}

func (m *MockGameState) PurchaseUpgradeAction(buildingIndex int) (bool, string) {
	m.purchaseUpgradeActionCalled = true
	return true, ""
}

var _ = Describe("Decider", func() {
	var (
		gameState *MockGameState
		decider   Decider
	)

	BeforeEach(func() {
		gameState = &MockGameState{
			money: 0,
			buildings: []Building{
				{name: "Building 1", baseGenerateRate: 1.0},
				{name: "Building 2", baseGenerateRate: 2.0},
			},
			upgrades: []Upgrade{
				{name: "Upgrade 1", isReleased: func(g GameState) bool { return true }, isTargetManualWork: true, effect: func(value float64) float64 { return value * 2 }},
			},
			manualWorkCalled:             false,
			updateBuildingsCalled:        false,
			getTotalGenerateRateCalled:   false,
			purchaseBuildingActionCalled: false,
			purchaseUpgradeActionCalled:  false,
		}
		decider = NewDefaultDecider(gameState)
	})

	Context("Decide", func() {
		It("should call ManualWork when cursor is 0", func() {
			success, message := decider.Decide(0, 0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal(""))
			Expect(gameState.manualWorkCalled).To(BeTrue())
			Expect(gameState.updateBuildingsCalled).To(BeFalse())
			Expect(gameState.getTotalGenerateRateCalled).To(BeFalse())
			Expect(gameState.purchaseBuildingActionCalled).To(BeFalse())
			Expect(gameState.purchaseUpgradeActionCalled).To(BeFalse())
			Expect(gameState.GetManualWork().count).To(Equal(1))
		})

		It("should call PurchaseBuildingAction when page is 0 and cursor is not 0", func() {
			success, message := decider.Decide(0, 1)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal(""))
			Expect(gameState.manualWorkCalled).To(BeFalse())
			Expect(gameState.updateBuildingsCalled).To(BeFalse())
			Expect(gameState.getTotalGenerateRateCalled).To(BeFalse())
			Expect(gameState.purchaseBuildingActionCalled).To(BeTrue())
			Expect(gameState.purchaseUpgradeActionCalled).To(BeFalse())
		})

		It("should call PurchaseUpgradeAction when page is 1 and cursor is not 0", func() {
			success, message := decider.Decide(1, 1)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal(""))
			Expect(gameState.manualWorkCalled).To(BeFalse())
			Expect(gameState.updateBuildingsCalled).To(BeFalse())
			Expect(gameState.getTotalGenerateRateCalled).To(BeFalse())
			Expect(gameState.purchaseBuildingActionCalled).To(BeFalse())
			Expect(gameState.purchaseUpgradeActionCalled).To(BeTrue())
		})

		It("should return false for invalid page selection", func() {
			success, message := decider.Decide(2, 1)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Invalid page selection"))
			Expect(gameState.manualWorkCalled).To(BeFalse())
			Expect(gameState.updateBuildingsCalled).To(BeFalse())
			Expect(gameState.getTotalGenerateRateCalled).To(BeFalse())
			Expect(gameState.purchaseBuildingActionCalled).To(BeFalse())
			Expect(gameState.purchaseUpgradeActionCalled).To(BeFalse())
		})
	})
})
