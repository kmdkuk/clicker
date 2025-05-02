package state

import (
	"encoding/json"
)

type Storage interface {
	SaveGameState(state GameState) error
	LoadGameState() (GameState, error)
}

type DefaultStorage struct {
	storageDriver StorageDriver
}

func NewDefaultStorage(driver StorageDriver) Storage {
	return &DefaultStorage{
		storageDriver: driver,
	}
}

func (s *DefaultStorage) SaveGameState(state GameState) error {
	save := ConverToSave(state)
	data, err := json.Marshal(save)
	if err != nil {
		return err
	}
	return s.storageDriver.SaveData(data)
}

func (s *DefaultStorage) LoadGameState() (GameState, error) {
	data, err := s.storageDriver.LoadData()
	if err != nil {
		return &DefaultGameState{}, err
	}
	var save Save
	if err := json.Unmarshal(data, &save); err != nil {
		return &DefaultGameState{}, err
	}

	if err := save.Validation(); err != nil {
		return &DefaultGameState{}, err
	}

	return save.ConvertToGameState()
}
