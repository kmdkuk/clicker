package presentation

import (
	"time"

	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/domain/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockGameState struct {
	money                        float64
	manualWork                   model.ManualWork
	buildings                    []model.Building
	upgrades                     []model.Upgrade
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
func (m *MockGameState) GetManualWork() *model.ManualWork {
	return &m.manualWork
}
func (m *MockGameState) SetManualWork(manualWork model.ManualWork) {
	m.manualWork = manualWork
}
func (m *MockGameState) GetBuildings() []model.Building {
	return m.buildings
}
func (m *MockGameState) GetUpgrades() []model.Upgrade {
	return m.upgrades
}
func (m *MockGameState) SetUpgrades(upgrades []model.Upgrade) {
	m.upgrades = upgrades
}
func (m *MockGameState) ManualWorkAction() {
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
		decider           Decider
		manualWorkUseCase *MockManualWorkUseCase
		buildingUseCase   *MockBuildingUseCase
		upgradeUseCase    *MockUpgradeUseCase
	)

	BeforeEach(func() {
		manualWorkUseCase = &MockManualWorkUseCase{
			ManualWorkActionCalled: false,
		}
		buildingUseCase = &MockBuildingUseCase{
			PurchaseBuildingActionCalled: false,
			buildings: []dto.Building{
				{Name: "Building 1"},
				{Name: "Building 2"},
			},
		}
		upgradeUseCase = &MockUpgradeUseCase{
			PurchaseUpgradeActionCalled: false,
		}
		decider = NewDecider(
			manualWorkUseCase,
			buildingUseCase,
			upgradeUseCase,
		)
	})

	Context("Decide", func() {
		It("should call ManualWork when cursor is 0", func() {
			success, message := decider.Decide(0, 0)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal(""))
			Expect(manualWorkUseCase.ManualWorkActionCalled).To(BeTrue())
			Expect(buildingUseCase.PurchaseBuildingActionCalled).To(BeFalse())
			Expect(upgradeUseCase.PurchaseUpgradeActionCalled).To(BeFalse())
		})

		It("should call PurchaseBuildingAction when page is 0 and cursor is not 0", func() {
			success, message := decider.Decide(0, 1)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal(""))
			Expect(manualWorkUseCase.ManualWorkActionCalled).To(BeFalse())
			Expect(buildingUseCase.PurchaseBuildingActionCalled).To(BeTrue())
			Expect(upgradeUseCase.PurchaseUpgradeActionCalled).To(BeFalse())
		})

		It("should call PurchaseUpgradeAction when page is 1 and cursor is not 0", func() {
			success, message := decider.Decide(1, 1)
			Expect(success).To(BeTrue())
			Expect(message).To(Equal(""))
			Expect(manualWorkUseCase.ManualWorkActionCalled).To(BeFalse())
			Expect(buildingUseCase.PurchaseBuildingActionCalled).To(BeFalse())
			Expect(upgradeUseCase.PurchaseUpgradeActionCalled).To(BeTrue())
		})

		It("should return false for invalid page selection", func() {
			success, message := decider.Decide(2, 1)
			Expect(success).To(BeFalse())
			Expect(message).To(Equal("Invalid page selection"))
			Expect(manualWorkUseCase.ManualWorkActionCalled).To(BeFalse())
			Expect(buildingUseCase.PurchaseBuildingActionCalled).To(BeFalse())
			Expect(upgradeUseCase.PurchaseUpgradeActionCalled).To(BeFalse())
		})
	})
})
