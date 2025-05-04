//go:build js && wasm
// +build js,wasm

package driver

import (
	"errors"
	"syscall/js"

	"github.com/kmdkuk/clicker/config"
)

func NewStorageDriver(key string) StorageDriver {
	if key == "" {
		return &StorageWasm{
			key: config.DefaultSaveKey,
		}
	}
	return &StorageWasm{
		key: key,
	}
}

type StorageWasm struct {
	key string
}

func (s *StorageWasm) SaveData(data []byte) error {
	localStorage := js.Global().Get("localStorage")
	if localStorage.IsUndefined() {
		return errors.New("localStorage is not available")
	}
	localStorage.Call("setItem", s.key, string(data))
	return nil
}

func (s *StorageWasm) LoadData() ([]byte, error) {
	localStorage := js.Global().Get("localStorage")
	if localStorage.IsUndefined() {
		return nil, errors.New("localStorage is not available")
	}
	data := localStorage.Call("getItem", s.key)
	if data.IsNull() || data.IsUndefined() {
		return nil, nil // Return nil if no data is found.
	}
	return []byte(data.String()), nil
}

func (s *StorageWasm) GetKeyName() string {
	return s.key
}
