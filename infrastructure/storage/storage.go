package storage

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/kmdkuk/clicker/infrastructure/state"
	"github.com/kmdkuk/clicker/infrastructure/storage/driver"
)

type Storage interface {
	SaveGameState(state state.GameState) error
	LoadGameState() (state.GameState, error)
}

type DefaultStorage struct {
	storageDriver        driver.StorageDriver
	haveOccuredLoadError bool
}

func NewDefaultStorage(driver driver.StorageDriver) Storage {
	return &DefaultStorage{
		storageDriver: driver,
	}
}

// SaveGameState encodes the game state to JSON and saves it
func (s *DefaultStorage) SaveGameState(state state.GameState) error {
	// Convert to save format
	save := ConverToSave(state)

	// Marshal to JSON
	data, err := json.Marshal(save)
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Create a backup of the current save before writing new data
	if s.haveOccuredLoadError {
		if err := s.createBackup(); err != nil {
			// Log the error but continue with the save
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
		s.haveOccuredLoadError = false
	}

	return s.storageDriver.SaveData(data)
}

// LoadGameState loads and decodes the game state, recovering partial data if possible
func (s *DefaultStorage) LoadGameState() (state.GameState, error) {
	s.haveOccuredLoadError = false
	data, err := s.storageDriver.LoadData()
	if err != nil {
		s.haveOccuredLoadError = true
		return &state.DefaultGameState{}, fmt.Errorf("failed to load data: %w", err)
	}

	var oldSaveConverted Save
	var oldSave oldSave
	// Try to unmarshal old save format first
	if err := json.Unmarshal(data, &oldSave); err != nil {
		oldSave, err = s.recoverPartialOldSave(data)
		if err != nil {
			fmt.Printf("Error recovering partial save as old save: %v\n", err)
		}
	}
	fmt.Printf("Marshaled old save: %+v\n", oldSave)
	// Convert old save to new format
	oldSaveConverted = Save{
		Money:      oldSave.Money,
		Buildings:  oldSave.Buildings,
		Upgradings: make([]upgrade, len(oldSave.Upgradings)),
		ManualWork: oldSave.ManualWork,
	}
	for i, u := range oldSave.Upgradings {
		oldSaveConverted.Upgradings[i] = upgrade{
			ID:          u.ID,
			IsPurchased: u.IsPurchased,
		}
	}

	var save Save
	// Try standard unmarshaling first
	if err := json.Unmarshal(data, &save); err != nil {
		s.haveOccuredLoadError = true
		fmt.Printf("failed to unmarshal: %v\n", err)
		// If standard unmarshaling fails, try partial recovery
		recoveredSave, recoverErr := s.recoverPartialSave(data)
		if recoverErr != nil {
			return &state.DefaultGameState{}, fmt.Errorf("cannot recover data: %w", recoverErr)
		}
		recoveredSave.merge(oldSaveConverted)
		gameState, err := recoveredSave.ConvertToGameState()
		if err != nil {
			s.haveOccuredLoadError = true
			return &state.DefaultGameState{}, fmt.Errorf("failed to convert recovered save: %w", err)
		}
		// Auto-save the fixed state
		if err := s.SaveGameState(gameState); err != nil {
			fmt.Printf("Warning: Failed to save fixed state: %v\n", err)
		}
		return gameState, nil
	}

	// Validate the save data
	validationErr := save.Validation()
	if validationErr == nil {
		// Normal path - convert valid save to game state
		save.merge(oldSaveConverted)
		return save.ConvertToGameState()
	}
	s.haveOccuredLoadError = true
	fmt.Printf("Validation error: %v\n", validationErr)
	// If validation fails, try to recover what we can
	fixedSave, fixErr := s.fixInvalidSave(save, validationErr)
	if fixErr != nil {
		s.haveOccuredLoadError = true
		return &state.DefaultGameState{}, fmt.Errorf("failed to fix invalid save: %w", fixErr)
	}

	// Convert the fixed save to game state
	fixedSave.merge(oldSaveConverted)
	gameState, err := fixedSave.ConvertToGameState()
	if err != nil {
		s.haveOccuredLoadError = true
		return &state.DefaultGameState{}, fmt.Errorf("failed to convert fixed save: %w", err)
	}

	// Auto-save the fixed state
	if err := s.SaveGameState(gameState); err != nil {
		fmt.Printf("Warning: Failed to save fixed state: %v\n", err)
	}

	return gameState, nil
}

// createBackup creates a backup of the current save file
func (s *DefaultStorage) createBackup() error {
	data, err := s.storageDriver.LoadData()
	if err != nil {
		return err
	}

	// Get base filename from driver
	baseFilename := s.storageDriver.GetKeyName()
	if baseFilename == "" {
		baseFilename = "save.json"
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupFilename := filepath.Join(
		filepath.Dir(baseFilename),
		fmt.Sprintf("%s.%s.bak",
			filepath.Base(baseFilename),
			timestamp,
		),
	)

	// Create a backup driver
	backupDriver := driver.NewStorageDriver(backupFilename)
	return backupDriver.SaveData(data)
}

// recoverPartialSave attempts to recover any valid parts from corrupted JSON
func (s *DefaultStorage) recoverPartialSave(data []byte) (Save, error) {
	// Create a default state to merge recovered data into
	save := Save{}
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &m); err != nil {
		return save, fmt.Errorf("failed to unmarshal data for partial recovery: %w", err)
	}
	// Try to extract money
	var partialSave struct {
		Money      *float64  `json:"money"`
		Buildings  []int     `json:"buildings"`
		Upgradings []upgrade `json:"upgradings"`
		ManualWork int       `json:"manualWork"`
		json.RawMessage
	}
	if err := unmarshalPartial(&partialSave.Money, m, "money"); err == nil && partialSave.Money != nil && *partialSave.Money > 0 {
		save.Money = *partialSave.Money
		fmt.Println("Partially recovered money from corrupted save: ", partialSave.Money)
	}

	if err := unmarshalPartial(&partialSave.Buildings, m, "buildings"); err == nil && partialSave.Buildings != nil {
		// Process valid buildings
		for i, building := range partialSave.Buildings {
			if i < len(save.Buildings) {
				// Only copy valid count values
				if building >= 0 {
					save.Buildings[i] = building
					fmt.Println("Partially recovered buildings count from corrupted save [", i, "]: ", building)
				}
			}
		}
	}

	// Try to extract upgrades
	if err := unmarshalPartial(&partialSave.Upgradings, m, "upgradings"); err == nil && partialSave.Upgradings != nil {
		// Process valid upgrades
		for i := range partialSave.Upgradings {
			for j := range save.Upgradings {
				if partialSave.Upgradings[i].ID == save.Upgradings[j].ID {
					save.Upgradings[j].IsPurchased = partialSave.Upgradings[i].IsPurchased
					fmt.Println("Partially recovered upgradings from corrupted save [", save.Upgradings[j].ID, "]: ", save.Upgradings[j].IsPurchased)
					break
				}
			}
		}
	}

	// Try to extract manual work
	if err := unmarshalPartial(&partialSave.ManualWork, m, "manualWork"); err == nil {
		if partialSave.ManualWork >= 0 {
			save.ManualWork = partialSave.ManualWork
			fmt.Println("Partially recovered manual work from corrupted save: ", partialSave.ManualWork)
		}
	}

	// Log recovery attempt
	fmt.Println("Partially recovered game state from corrupted save")

	return save, nil
}

func (s *DefaultStorage) recoverPartialOldSave(data []byte) (oldSave, error) {
	// Create a default state to merge recovered data into
	save := oldSave{}
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &m); err != nil {
		return save, fmt.Errorf("failed to unmarshal data for partial recovery: %w", err)
	}
	// Try to extract money
	var partialSave struct {
		Money      *float64  `json:"Money"`
		Buildings  []int     `json:"Buildings"`
		Upgradings []upgrade `json:"Upgradings"`
		ManualWork int       `json:"ManualWork"`
		json.RawMessage
	}
	if err := unmarshalPartial(&partialSave.Money, m, "Money"); err == nil && partialSave.Money != nil && *partialSave.Money > 0 {
		save.Money = *partialSave.Money
		fmt.Println("Partially recovered money from corrupted save: ", partialSave.Money)
	}

	if err := unmarshalPartial(&partialSave.Buildings, m, "Buildings"); err == nil && partialSave.Buildings != nil {
		// Process valid buildings
		save.Buildings = make([]int, len(partialSave.Buildings))
		for i, building := range partialSave.Buildings {
			if i < len(save.Buildings) {
				// Only copy valid count values
				if building >= 0 {
					save.Buildings[i] = building
					fmt.Println("Partially recovered buildings count from corrupted save [", i, "]: ", building)
				}
			}
		}
	}

	// Try to extract upgrades
	if err := unmarshalPartial(&partialSave.Upgradings, m, "Upgradings"); err == nil && partialSave.Upgradings != nil {
		// Process valid upgrades
		for i := range partialSave.Upgradings {
			for j := range save.Upgradings {
				if partialSave.Upgradings[i].ID == save.Upgradings[j].ID {
					save.Upgradings[j].IsPurchased = partialSave.Upgradings[i].IsPurchased
					fmt.Println("Partially recovered upgradings from corrupted save [", save.Upgradings[j].ID, "]: ", save.Upgradings[j].IsPurchased)
					break
				}
			}
		}
	}

	// Try to extract manual work
	if err := unmarshalPartial(&partialSave.ManualWork, m, "ManualWork"); err == nil {
		if partialSave.ManualWork >= 0 {
			save.ManualWork = partialSave.ManualWork
			fmt.Println("Partially recovered manual work from corrupted save: ", partialSave.ManualWork)
		}
	}

	// Log recovery attempt
	fmt.Println("Partially recovered game state from corrupted save")

	return save, nil
}

// fixInvalidSave attempts to fix validation errors in the save data
func (s *DefaultStorage) fixInvalidSave(save Save, validationErr error) (Save, error) {
	// Log the validation error
	fmt.Printf("Fixing invalid save: %v\n", validationErr)

	// Create default save to fill in missing pieces
	defaultSave := ConverToSave(&state.DefaultGameState{})

	// Fix money if negative
	if save.Money < 0 {
		save.Money = defaultSave.Money
	}

	// Fix buildings
	if save.Buildings == nil {
		save.Buildings = defaultSave.Buildings
	} else {
		for i, building := range save.Buildings {
			// Fix negative counts
			if building < 0 {
				save.Buildings[i] = 0
			}
		}
	}

	// Fix upgrades
	if save.Upgradings == nil {
		save.Upgradings = defaultSave.Upgradings
	} else {
		for i, isPurchased := range save.Upgradings {
			if i < len(defaultSave.Upgradings) {
				save.Upgradings[i] = isPurchased
			}
		}
	}

	// Fix manual work
	if save.ManualWork < 0 {
		save.ManualWork = defaultSave.ManualWork
	}

	// Validate the fixed save
	if err := save.Validation(); err != nil {
		// If we still have validation errors, log them but continue with what we have
		fmt.Printf("Warning: Still have validation errors after fixing: %v\n", err)
	}

	return save, nil
}

func unmarshalPartial(to interface{}, m map[string]json.RawMessage, s string) error {
	if err := json.Unmarshal(m[s], to); err != nil {
		return fmt.Errorf("failed to unmarshal data for partial recovery: %w", err)
	}
	delete(m, s)
	return nil
}

func (s *Save) merge(other Save) {
	if s.Money < other.Money {
		s.Money = other.Money
	}
	if s.ManualWork < other.ManualWork {
		s.ManualWork = other.ManualWork
	}
	s.Buildings = append(s.Buildings, make([]int, len(other.Buildings)-len(s.Buildings))...)
	for i, b := range s.Buildings {
		if i < len(other.Buildings) && other.Buildings[i] > b {
			s.Buildings[i] = other.Buildings[i]
		}
	}
	s.Upgradings = append(s.Upgradings, make([]upgrade, len(other.Upgradings)-len(s.Upgradings))...)
	for i := range s.Upgradings {
		for j := range other.Upgradings {
			if s.Upgradings[i].ID == other.Upgradings[j].ID {
				if other.Upgradings[j].IsPurchased {
					s.Upgradings[i].IsPurchased = other.Upgradings[j].IsPurchased
				}
			}
		}
	}
}
