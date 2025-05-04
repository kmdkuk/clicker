package storage

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/infrastructure/state"
	"github.com/kmdkuk/clicker/infrastructure/storage/driver"
)

type Storage interface {
	SaveGameState(state state.GameState) error
	LoadGameState() (state.GameState, error)
}

type DefaultStorage struct {
	storageDriver driver.StorageDriver
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
	if err := s.createBackup(); err != nil {
		// Log the error but continue with the save
		fmt.Printf("Warning: Failed to create backup: %v\n", err)
	}

	return s.storageDriver.SaveData(data)
}

// LoadGameState loads and decodes the game state, recovering partial data if possible
func (s *DefaultStorage) LoadGameState() (state.GameState, error) {
	data, err := s.storageDriver.LoadData()
	if err != nil {
		return &state.DefaultGameState{}, fmt.Errorf("failed to load data: %w", err)
	}

	// Create a backup of the raw save data before processing
	if err := s.backupRawData(data); err != nil {
		fmt.Printf("Warning: Failed to backup raw data: %v\n", err)
	}

	// Try standard unmarshaling first
	var save Save
	if err := json.Unmarshal(data, &save); err != nil {
		fmt.Printf("failed to unmarshal: %v\n", err)
		// If standard unmarshaling fails, try partial recovery
		recoveredState, recoverErr := s.recoverPartialState(data)
		if recoverErr != nil {
			return &state.DefaultGameState{}, fmt.Errorf("cannot recover data: %w", recoverErr)
		}
		return recoveredState, nil
	}

	// Validate the save data
	validationErr := save.Validation()
	if validationErr == nil {
		// Normal path - convert valid save to game state
		return save.ConvertToGameState()
	}
	// If validation fails, try to recover what we can
	fixedSave, fixErr := s.fixInvalidSave(save, validationErr)
	if fixErr != nil {
		return &state.DefaultGameState{}, fmt.Errorf("failed to fix invalid save: %w", fixErr)
	}

	// Convert the fixed save to game state
	gameState, err := fixedSave.ConvertToGameState()
	if err != nil {
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

// backupRawData creates a backup of raw save data
func (s *DefaultStorage) backupRawData(data []byte) error {
	baseFilename := s.storageDriver.GetKeyName()
	if baseFilename == "" {
		baseFilename = config.DefaultSaveKey
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFilename := fmt.Sprintf("%s.%s.raw.bak",
		baseFilename,
		timestamp,
	)

	backupDriver := driver.NewStorageDriver(backupFilename)
	return backupDriver.SaveData(data)
}

// recoverPartialState attempts to recover any valid parts from corrupted JSON
func (s *DefaultStorage) recoverPartialState(data []byte) (state.GameState, error) {
	// Create a default state to merge recovered data into
	gameState := &state.DefaultGameState{}
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &m); err != nil {
		return gameState, fmt.Errorf("failed to unmarshal data for partial recovery: %w", err)
	}
	// Try to extract money
	var partialSave struct {
		Money      *float64 `json:"money"`
		Buildings  []int    `json:"buildings"`
		Upgradings []bool   `json:"upgradings"`
		ManualWork int      `json:"manualWork"`
		json.RawMessage
	}
	if err := unmarshalPartial(&partialSave.Money, m, "money"); err == nil && partialSave.Money != nil && *partialSave.Money > 0 {
		gameState.Money = *partialSave.Money
		fmt.Println("Partially recovered money from corrupted save: ", partialSave.Money)
	}

	if err := unmarshalPartial(&partialSave.Buildings, m, "buildings"); err == nil && partialSave.Buildings != nil {
		// Process valid buildings
		for i, building := range partialSave.Buildings {
			if i < len(gameState.Buildings) {
				// Only copy valid count values
				if building >= 0 {
					gameState.Buildings[i].Count = building
					fmt.Println("Partially recovered buildings count from corrupted save [", i, "]: ", building)
				}
			}
		}
	}

	// Try to extract upgrades
	if err := unmarshalPartial(&partialSave.Upgradings, m, "upgradings"); err == nil && partialSave.Upgradings != nil {
		// Process valid upgrades
		for i, isPurchased := range partialSave.Upgradings {
			if i < len(gameState.Upgrades) {
				gameState.Upgrades[i].IsPurchased = isPurchased
				fmt.Println("Partially recovered upgradings from corrupted save [", i, "]: ", isPurchased)
			}
		}
	}

	// Try to extract manual work
	if err := unmarshalPartial(&partialSave.ManualWork, m, "manualWork"); err == nil {
		if partialSave.ManualWork >= 0 {
			gameState.ManualWork.Count = partialSave.ManualWork
			fmt.Println("Partially recovered manual work from corrupted save: ", partialSave.ManualWork)
		}
	}

	// Log recovery attempt
	fmt.Println("Partially recovered game state from corrupted save")

	return gameState, nil
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
