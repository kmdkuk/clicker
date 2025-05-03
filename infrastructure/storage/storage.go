package storage

import (
	"encoding/json"

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

func (s *DefaultStorage) SaveGameState(state state.GameState) error {
	save := ConverToSave(state)
	data, err := json.Marshal(save)
	if err != nil {
		return err
	}
	return s.storageDriver.SaveData(data)
}

func (s *DefaultStorage) LoadGameState() (state.GameState, error) {
	data, err := s.storageDriver.LoadData()
	if err != nil {
		return &state.DefaultGameState{}, err
	}
	var save Save
	if err := json.Unmarshal(data, &save); err != nil {
		return &state.DefaultGameState{}, err
	}

	if err := save.Validation(); err != nil {
		return &state.DefaultGameState{}, err
	}

	return save.ConvertToGameState()
}
