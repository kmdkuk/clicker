package usecase_test

import (
	"errors"
	"time"

	"github.com/kmdkuk/clicker/application/usecase"
	"github.com/kmdkuk/clicker/domain/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockGameState struct {
	Money               float64
	Upgrades            []model.Upgrade
	SetUpgradeCallCount int
	SetUpgradeError     error
	UpdateMoneyAmount   float64
}

func (m *MockGameState) GetMoney() float64 {
	return m.Money
}

func (m *MockGameState) SetMoney(amount float64) {
	m.Money = amount
}

func (m *MockGameState) SetUpgrades(upgrades []model.Upgrade) {
	m.Upgrades = upgrades
}

func (m *MockGameState) GetUpgrades() []model.Upgrade {
	return m.Upgrades
}

func (m *MockGameState) SetUpgradesIsPurchased(upgradeIndex int, isPurchased bool) error {
	if m.SetUpgradeError != nil {
		return m.SetUpgradeError
	}
	m.SetUpgradeCallCount++
	m.Upgrades[upgradeIndex].IsPurchased = isPurchased
	return nil
}

func (m *MockGameState) UpdateMoney(amount float64) {
	m.UpdateMoneyAmount = amount
	m.Money += amount
}

func (m *MockGameState) GetTotalGenerateRate() float64 {
	return 0.0
}
func (m *MockGameState) GetManualWork() *model.ManualWork {
	return &model.ManualWork{}
}
func (m *MockGameState) SetManualWorkCount(count int) error {
	return nil
}
func (m *MockGameState) GetBuildings() []model.Building {
	return nil
}
func (m *MockGameState) SetBuildingCount(index int, count int) error {
	return nil
}
func (m *MockGameState) UpdateBuildings(_ time.Time) {
}

var _ = Describe("UpgradeUseCase", func() {
	var (
		mockGameState  *MockGameState
		upgradeUseCase *usecase.UpgradeUseCase
	)

	BeforeEach(func() {
		mockGameState = &MockGameState{
			Money: 100.0,
			Upgrades: []model.Upgrade{
				{
					Name:        "Basic Upgrade",
					Cost:        50.0,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true },
				},
				{
					Name:        "Premium Upgrade",
					Cost:        150.0,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true },
				},
				{
					Name:        "Limited Upgrade",
					Cost:        30.0,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return false },
				},
				{
					Name:        "Purchased Upgrade",
					Cost:        20.0,
					IsPurchased: true,
					IsReleased:  func(_ model.GameStateReader) bool { return true },
				},
			},
		}
		upgradeUseCase = usecase.NewUpgradeUseCase(mockGameState)
	})

	Describe("GetUpgrades", func() {
		Context("when upgrades exist in the game state", func() {
			It("should return the correct number of upgrade DTOs", func() {
				upgrades := upgradeUseCase.GetUpgrades()
				Expect(upgrades).To(HaveLen(4))
			})

			It("should correctly map model upgrades to DTOs", func() {
				upgrades := upgradeUseCase.GetUpgrades()

				// Check first upgrade
				Expect(upgrades[0].Name).To(Equal("Basic Upgrade"))
				Expect(upgrades[0].Cost).To(Equal(50.0))
				Expect(upgrades[0].IsPurchased).To(BeFalse())
				Expect(upgrades[0].IsReleased).To(BeTrue())

				// Check already purchased upgrade
				Expect(upgrades[3].Name).To(Equal("Purchased Upgrade"))
				Expect(upgrades[3].IsPurchased).To(BeTrue())
			})
		})

		Context("when no upgrades exist in the game state", func() {
			BeforeEach(func() {
				mockGameState.Upgrades = []model.Upgrade{}
				upgradeUseCase = usecase.NewUpgradeUseCase(mockGameState)
			})

			It("should return an empty slice", func() {
				upgrades := upgradeUseCase.GetUpgrades()
				Expect(upgrades).To(BeEmpty())
			})
		})
	})

	Describe("PurchaseUpgradeAction", func() {
		Context("when the upgrade index is valid", func() {
			It("should successfully purchase an available upgrade", func() {
				success, message := upgradeUseCase.PurchaseUpgradeAction(0)

				Expect(success).To(BeTrue())
				Expect(message).To(Equal("Upgrade purchased successfully!"))
				Expect(mockGameState.SetUpgradeCallCount).To(Equal(1))
				Expect(mockGameState.UpdateMoneyAmount).To(Equal(-50.0))
				Expect(mockGameState.Money).To(Equal(50.0))
			})

			It("should fail when trying to purchase an already purchased upgrade", func() {
				success, message := upgradeUseCase.PurchaseUpgradeAction(3)

				Expect(success).To(BeFalse())
				Expect(message).To(Equal("Upgrade already purchased!"))
				Expect(mockGameState.SetUpgradeCallCount).To(Equal(0))
			})

			It("should fail when trying to purchase an unreleased upgrade", func() {
				success, message := upgradeUseCase.PurchaseUpgradeAction(2)

				Expect(success).To(BeFalse())
				Expect(message).To(Equal("Upgrade not available yet!"))
				Expect(mockGameState.SetUpgradeCallCount).To(Equal(0))
			})

			It("should fail when not enough money", func() {
				success, message := upgradeUseCase.PurchaseUpgradeAction(1) // Premium upgrade costs 150

				Expect(success).To(BeFalse())
				Expect(message).To(Equal("Not enough money for upgrade!"))
				Expect(mockGameState.SetUpgradeCallCount).To(Equal(0))
			})

			It("should handle errors from SetUpgradesIsPurchased", func() {
				mockGameState.SetUpgradeError = errors.New("database error")

				success, message := upgradeUseCase.PurchaseUpgradeAction(0)

				Expect(success).To(BeFalse())
				Expect(message).To(Equal("Failed to purchase upgrade!"))
				Expect(mockGameState.Money).To(Equal(100.0)) // Money should not be deducted
			})
		})

		Context("when the upgrade index is invalid", func() {
			It("should fail with negative index", func() {
				success, message := upgradeUseCase.PurchaseUpgradeAction(-1)

				Expect(success).To(BeFalse())
				Expect(message).To(Equal("Invalid upgrade selection!"))
			})

			It("should fail with index beyond array bounds", func() {
				success, message := upgradeUseCase.PurchaseUpgradeAction(10)

				Expect(success).To(BeFalse())
				Expect(message).To(Equal("Invalid upgrade selection!"))
			})
		})
	})
})
