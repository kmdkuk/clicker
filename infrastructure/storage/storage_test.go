package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kmdkuk/clicker/domain/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock implementation of StorageDriver
type MockStorageDriver struct {
	SaveDataCalled bool
	LoadDataCalled bool
	Data           []byte
	SaveError      error
	LoadError      error
	Filename       string
}

func (m *MockStorageDriver) SaveData(data []byte) error {
	m.SaveDataCalled = true
	m.Data = data
	return m.SaveError
}

func (m *MockStorageDriver) LoadData() ([]byte, error) {
	m.LoadDataCalled = true
	return m.Data, m.LoadError
}

func (m *MockStorageDriver) GetKeyName() string {
	return m.Filename
}

// Mock implementation of GameState
type MockGameState struct {
	Money      float64
	Buildings  []model.Building
	Upgrades   []model.Upgrade
	ManualWork model.ManualWork
}

func (m *MockGameState) GetMoney() float64 {
	return m.Money
}

func (m *MockGameState) UpdateMoney(amount float64) {
	m.Money += amount
}

func (m *MockGameState) GetBuildings() []model.Building {
	return m.Buildings
}

func (m *MockGameState) GetUpgrades() []model.Upgrade {
	return m.Upgrades
}

func (m *MockGameState) GetManualWork() *model.ManualWork {
	return &m.ManualWork
}

func (m *MockGameState) SetBuildingCount(index int, count int) error {
	if index < 0 || index >= len(m.Buildings) {
		return errors.New("invalid building index")
	}
	m.Buildings[index].Count = count
	return nil
}

func (m *MockGameState) SetUpgrades(upgrades []model.Upgrade) {
	m.Upgrades = upgrades
}

func (m *MockGameState) SetUpgradesIsPurchased(index int, isPurchased bool) error {
	if index < 0 || index >= len(m.Upgrades) {
		return errors.New("invalid upgrade index")
	}
	m.Upgrades[index].IsPurchased = isPurchased
	return nil
}

func (m *MockGameState) SetManualWorkCount(count int) error {
	if count < 0 {
		return errors.New("invalid manual work count")
	}
	m.ManualWork.Count = count
	return nil
}
func (m *MockGameState) UpdateBuildings(now time.Time) {
	// Mock implementation, no action needed
}
func (m *MockGameState) GetTotalGenerateRate() float64 {
	return 0.0
}

func (m *MockGameState) GetBuildingCount(index int) (int, error) {
	if index < 0 || index >= len(m.Buildings) {
		return 0, errors.New("invalid building index")
	}
	return m.Buildings[index].Count, nil
}

var _ = Describe("DefaultStorage", func() {
	var (
		mockDriver  *MockStorageDriver
		testStorage Storage
		testState   *MockGameState
	)

	BeforeEach(func() {
		mockDriver = &MockStorageDriver{
			Filename: "test_save.json",
		}
		testStorage = NewDefaultStorage(mockDriver)
		testState = &MockGameState{
			Money: 100.0,
			Buildings: []model.Building{
				{ID: 0, Name: "Building 1", Count: 5, BaseCost: 10},
				{ID: 1, Name: "Building 2", Count: 3, BaseCost: 50},
			},
			Upgrades: []model.Upgrade{
				{Name: "Upgrade 1", IsPurchased: true, Cost: 20},
				{Name: "Upgrade 2", IsPurchased: false, Cost: 100},
			},
			ManualWork: model.ManualWork{
				Count: 10,
			},
		}
	})

	AfterEach(func() {
		testDir := filepath.Dir(mockDriver.Filename)
		baseFilename := filepath.Base(mockDriver.Filename)

		if _, err := os.Stat(testDir); err == nil {
			files, err := filepath.Glob(filepath.Join(testDir, baseFilename+".*.bak"))
			if err == nil {
				for _, file := range files {
					err := os.Remove(file)
					if err != nil {
						fmt.Println("Error removing backup file:", err)
					}
				}
			}
		}
	})

	Describe("SaveGameState", func() {
		It("should convert game state to save format and save it", func() {
			err := testStorage.SaveGameState(testState)

			Expect(err).NotTo(HaveOccurred())
			Expect(mockDriver.SaveDataCalled).To(BeTrue())
			Expect(mockDriver.Data).NotTo(BeNil())

			// Verify the saved data contains the expected values
			var save Save
			err = json.Unmarshal(mockDriver.Data, &save)
			Expect(err).NotTo(HaveOccurred())

			Expect(save.Money).To(Equal(100.0))
			Expect(save.Buildings).To(HaveLen(2))
			Expect(save.Buildings[0]).To(Equal(5))
			Expect(save.Upgradings).To(HaveLen(2))
			Expect(save.Upgradings[0]).To(BeTrue())
			Expect(save.ManualWork).To(Equal(10))
		})

		It("should handle errors from SaveData", func() {
			mockDriver.SaveError = errors.New("save error")

			err := testStorage.SaveGameState(testState)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("save error"))
		})
	})

	Describe("LoadGameState", func() {
		Context("with valid save data", func() {
			BeforeEach(func() {
				// Create valid save data
				validSave := Save{
					Money:     250.0,
					Buildings: []int{7, 2},
					Upgradings: []bool{
						true,
					},
					ManualWork: 15,
				}

				data, err := json.Marshal(validSave)
				Expect(err).NotTo(HaveOccurred())
				mockDriver.Data = data
			})

			It("should load and convert save data to game state", func() {
				gameState, err := testStorage.LoadGameState()

				Expect(err).NotTo(HaveOccurred())
				Expect(mockDriver.LoadDataCalled).To(BeTrue())

				Expect(gameState.GetMoney()).To(Equal(250.0))
				Expect(gameState.GetBuildings()[0].Count).To(Equal(7))
				Expect(gameState.GetBuildings()[1].Count).To(Equal(2))
				Expect(gameState.GetUpgrades()).To(HaveLen(1))
				Expect(gameState.GetUpgrades()[0].IsPurchased).To(BeTrue())
				Expect(gameState.GetManualWork().Count).To(Equal(15))
			})
		})

		Context("with corrupted JSON data", func() {
			BeforeEach(func() {
				// Create corrupted JSON data
				mockDriver.Data = []byte(`{"money": 100.0, "buildings": [{"id": 0, "count": 5}, {"id": 1, "count": 3}], "upgradings": [{"id": 0, "isPurchased": false}], "manualWork": -10, "corrupted": true}`)
			})

			It("should attempt to recover partial state", func() {
				gameState, err := testStorage.LoadGameState()

				// Even with errors, we should get a usable game state
				Expect(err).NotTo(HaveOccurred())
				Expect(mockDriver.LoadDataCalled).To(BeTrue())

				// Should recover the money value
				Expect(gameState.GetMoney()).To(Equal(100.0))
			})
		})

		Context("with invalid save data", func() {
			BeforeEach(func() {
				// Create save data with validation errors
				invalidSave := Save{
					Money: -50.0, // Negative money
					Buildings: []int{
						-2, // Negative count
					},
					ManualWork: -5, // Negative manual work{
				}

				data, err := json.Marshal(invalidSave)
				Expect(err).NotTo(HaveOccurred())
				mockDriver.Data = data
			})

			It("should fix invalid values and return usable game state", func() {
				gameState, err := testStorage.LoadGameState()

				Expect(err).NotTo(HaveOccurred())
				Expect(mockDriver.LoadDataCalled).To(BeTrue())

				// Money should be fixed to non-negative value
				Expect(gameState.GetMoney()).To(BeNumerically(">=", 0))

				// Building count should be fixed to non-negative
				buildings := gameState.GetBuildings()
				if len(buildings) > 0 {
					Expect(buildings[0].Count).To(BeNumerically(">=", 0))
				}

				// Manual work count should be fixed
				Expect(gameState.GetManualWork().Count).To(BeNumerically(">=", 0))
			})
		})

		Context("when LoadData fails", func() {
			BeforeEach(func() {
				mockDriver.LoadError = errors.New("load error")
			})

			It("should return an error but still provide default game state", func() {
				gameState, err := testStorage.LoadGameState()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("load error"))

				// Should still return a usable default state
				Expect(gameState).NotTo(BeNil())
			})
		})
	})

	Describe("Recovery and backup functions", func() {
		// These are mostly tested through the LoadGameState and SaveGameState tests,
		// but we can add specific tests for edge cases

		Context("recoverPartialState", func() {
			It("should extract valid fields from partially corrupted JSON", func() {
				// This requires exposing recoverPartialState or testing through LoadGameState
				data := []byte(`{"money": 123.45, "buildings": [{"id": 0, "count": 7}], "manualWork": 99, "corrupted": true}`)
				mockDriver.Data = data

				gameState, err := testStorage.LoadGameState()

				Expect(err).NotTo(HaveOccurred())
				Expect(gameState.GetMoney()).To(Equal(123.45))
				Expect(gameState.GetManualWork().Count).To(Equal(99))
			})
		})

		Context("fixInvalidSave", func() {
			It("should replace invalid values with defaults", func() {
				invalidSave := Save{
					Money: -100.0,
					Buildings: []int{
						-5,
					},
				}

				data, _ := json.Marshal(invalidSave)
				mockDriver.Data = data

				gameState, err := testStorage.LoadGameState()

				Expect(err).NotTo(HaveOccurred())
				Expect(gameState.GetMoney()).To(BeNumerically(">=", 0))

				buildings := gameState.GetBuildings()
				if len(buildings) > 0 {
					Expect(buildings[0].Count).To(BeNumerically(">=", 0))
				}
			})
		})
	})
})
