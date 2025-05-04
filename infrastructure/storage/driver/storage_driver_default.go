//go:build !js && !wasm
// +build !js,!wasm

package driver

import (
	"os"
)

func NewStorageDriver(key string) StorageDriver {
	if key == "" {
		return &DefaultStorageDriver{
			path: DefaultSaveKey,
		}
	}
	return &DefaultStorageDriver{
		path: key,
	}
}

type DefaultStorageDriver struct {
	path string
}

func (s *DefaultStorageDriver) SaveData(data []byte) error {
	return os.WriteFile(s.path, data, 0644)
}
func (s *DefaultStorageDriver) LoadData() ([]byte, error) {
	return os.ReadFile(s.path)
}

func (s *DefaultStorageDriver) GetKeyName() string {
	return s.path
}
