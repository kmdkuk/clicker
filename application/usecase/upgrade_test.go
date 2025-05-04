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

func (m *MockGameState) SetUpgradesIsPurchasedWithID(ID string, isPurchased bool) error {
	for i := range m.Upgrades {
		if m.Upgrades[i].Name == ID {
			m.Upgrades[i].IsPurchased = isPurchased
			return nil
		}
	}
	return errors.New("upgrade not found")
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
					ID:          "0",
					Name:        "Basic Upgrade",
					Cost:        50.0,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true },
				},
				{
					ID:          "1",
					Name:        "Premium Upgrade",
					Cost:        150.0,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true },
				},
				{
					ID:          "2",
					Name:        "Limited Upgrade",
					Cost:        200.0,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true },
				},
				{
					ID:          "3",
					Name:        "Purchased Upgrade",
					Cost:        300.0,
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

	Describe("GetUpgradesIsReleasedCostSorted", func() {
		Context("when upgrades exist in different release states", func() {
			BeforeEach(func() {
				mockGameState.Upgrades = []model.Upgrade{
					{
						Name:        "Mid Cost Released",
						Cost:        50.0,
						IsPurchased: false,
						IsReleased:  func(_ model.GameStateReader) bool { return true },
					},
					{
						Name:        "High Cost Released",
						Cost:        100.0,
						IsPurchased: false,
						IsReleased:  func(_ model.GameStateReader) bool { return true },
					},
					{
						Name:        "Low Cost Released",
						Cost:        25.0,
						IsPurchased: false,
						IsReleased:  func(_ model.GameStateReader) bool { return true },
					},
					{
						Name:        "Lowest Cost Unreleased",
						Cost:        10.0,
						IsPurchased: false,
						IsReleased:  func(_ model.GameStateReader) bool { return false },
					},
					{
						Name:        "Highest Cost Released",
						Cost:        200.0,
						IsPurchased: true,
						IsReleased:  func(_ model.GameStateReader) bool { return true },
					},
				}
				upgradeUseCase = usecase.NewUpgradeUseCase(mockGameState)
			})

			It("should return only released upgrades", func() {
				upgrades := upgradeUseCase.GetUpgradesIsReleasedCostSorted()

				// Check we only get released upgrades
				Expect(upgrades).To(HaveLen(4))
				for _, upgrade := range upgrades {
					Expect(upgrade.IsReleased).To(BeTrue())
				}

				// Verify unreleased upgrades are filtered out
				found := false
				for _, upgrade := range upgrades {
					if upgrade.Name == "Lowest Cost Unreleased" {
						found = true
					}
				}
				Expect(found).To(BeFalse(), "Unreleased upgrades should be filtered out")
			})

			It("should sort upgrades by cost in ascending order", func() {
				upgrades := upgradeUseCase.GetUpgradesIsReleasedCostSorted()

				// Check sorting
				Expect(upgrades).To(HaveLen(4))
				for i := 0; i < len(upgrades)-1; i++ {
					Expect(upgrades[i].Cost).To(BeNumerically("<=", upgrades[i+1].Cost),
						"Upgrades should be sorted by cost in ascending order")
				}

				// Verify exact order
				Expect(upgrades[0].Name).To(Equal("Low Cost Released"), "Lowest cost should be first")
				Expect(upgrades[3].Name).To(Equal("Highest Cost Released"), "Highest cost should be last")
			})
		})

		Context("when all upgrades are unreleased", func() {
			BeforeEach(func() {
				mockGameState.Upgrades = []model.Upgrade{
					{
						Name:       "Unreleased 1",
						Cost:       50.0,
						IsReleased: func(_ model.GameStateReader) bool { return false },
					},
					{
						Name:       "Unreleased 2",
						Cost:       20.0,
						IsReleased: func(_ model.GameStateReader) bool { return false },
					},
				}
				upgradeUseCase = usecase.NewUpgradeUseCase(mockGameState)
			})

			It("should return an empty slice", func() {
				upgrades := upgradeUseCase.GetUpgradesIsReleasedCostSorted()
				Expect(upgrades).To(BeEmpty())
			})
		})
	})

	Context("edge cases for purchasing upgrades", func() {
		It("should succeed when money is exactly equal to upgrade cost", func() {
			// Set money to exactly the upgrade cost
			mockGameState.Money = mockGameState.Upgrades[0].Cost

			success, message := upgradeUseCase.PurchaseUpgradeAction(0)

			Expect(success).To(BeTrue())
			Expect(message).To(Equal("Upgrade purchased successfully!"))
			Expect(mockGameState.Money).To(Equal(float64(0)))
		})

		It("should maintain unchanged state when purchase fails", func() {
			originalMoney := mockGameState.Money
			originalUpgradeState := make([]bool, len(mockGameState.Upgrades))
			for i, upgrade := range mockGameState.Upgrades {
				originalUpgradeState[i] = upgrade.IsPurchased
			}

			// Attempt to purchase unavailable upgrade
			upgradeUseCase.PurchaseUpgradeAction(2)

			// Verify state remains unchanged
			Expect(mockGameState.Money).To(Equal(originalMoney))
			for i, upgrade := range mockGameState.Upgrades {
				Expect(upgrade.IsPurchased).To(Equal(originalUpgradeState[i]))
			}
		})
	})

	Context("when upgrade release status changes", func() {
		var dynamicReleaseStatus bool

		BeforeEach(func() {
			dynamicReleaseStatus = false
			mockGameState.Upgrades = []model.Upgrade{
				{
					Name:        "Dynamic Upgrade",
					Cost:        30.0,
					IsPurchased: false,
					IsReleased: func(_ model.GameStateReader) bool {
						return dynamicReleaseStatus
					},
				},
			}
			upgradeUseCase = usecase.NewUpgradeUseCase(mockGameState)
		})

		It("should reflect changes in release status", func() {
			// Initially not released
			upgrades := upgradeUseCase.GetUpgrades()
			Expect(upgrades[0].IsReleased).To(BeFalse())

			// Change release status
			dynamicReleaseStatus = true

			// Should now be released
			upgrades = upgradeUseCase.GetUpgrades()
			Expect(upgrades[0].IsReleased).To(BeTrue())

			// Should also appear in sorted released upgrades
			sortedUpgrades := upgradeUseCase.GetUpgradesIsReleasedCostSorted()
			Expect(sortedUpgrades).To(HaveLen(1))
			Expect(sortedUpgrades[0].Name).To(Equal("Dynamic Upgrade"))
		})
	})

	Describe("PurchaseUpgradeAction with filtered upgrades", func() {
		// フィルタリングしたアップグレードリストとオリジナルのリストの不一致を検証するテスト
		var (
			mockGameState    *MockGameState
			upgradeUseCase   *usecase.UpgradeUseCase
			originalUpgrades []model.Upgrade
		)

		BeforeEach(func() {
			// オリジナルのアップグレードリスト（フィルタリング前）
			originalUpgrades = []model.Upgrade{
				{ // インデックス0
					ID:          "0",
					Name:        "Hidden Upgrade 1",
					Cost:        10,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return false }, // 非表示
				},
				{ // インデックス1
					ID:          "1",
					Name:        "First Visible Upgrade",
					Cost:        20,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true }, // 表示
				},
				{ // インデックス2
					ID:          "2",
					Name:        "Hidden Upgrade 2",
					Cost:        30,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return false }, // 非表示
				},
				{ // インデックス3
					ID:          "3",
					Name:        "Second Visible Upgrade",
					Cost:        40,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true }, // 表示
				},
				{ // インデックス4
					ID:          "4",
					Name:        "Third Visible Upgrade",
					Cost:        50,
					IsPurchased: false,
					IsReleased:  func(_ model.GameStateReader) bool { return true }, // 表示
				},
			}

			mockGameState = &MockGameState{
				Money:    100.0, // 十分な資金
				Upgrades: originalUpgrades,
			}
			upgradeUseCase = usecase.NewUpgradeUseCase(mockGameState)
		})

		Context("when using cursor from filtered list", func() {
			It("should purchase to purchase intended upgrade", func() {
				filteredUpgrades := upgradeUseCase.GetUpgradesIsReleasedCostSorted()
				Expect(filteredUpgrades).To(HaveLen(3))
				Expect(filteredUpgrades[0].Name).To(Equal("First Visible Upgrade"))
				Expect(filteredUpgrades[1].Name).To(Equal("Second Visible Upgrade"))
				Expect(filteredUpgrades[2].Name).To(Equal("Third Visible Upgrade"))

				actualIndex := 3
				Expect(mockGameState.Upgrades[actualIndex].IsPurchased).To(BeFalse())

				uiCursor := 1 // 0-indexed
				success, _ := upgradeUseCase.PurchaseUpgradeAction(uiCursor)

				Expect(success).To(BeTrue())
				intendedUpgrade := mockGameState.Upgrades[actualIndex]
				Expect(intendedUpgrade.Name).To(Equal("Second Visible Upgrade"))
				Expect(intendedUpgrade.IsPurchased).To(BeTrue())
			})
		})
	})
})
