package driver

const DefaultSaveKey = "game_state.json"

type StorageDriver interface {
	SaveData(data []byte) error
	LoadData() ([]byte, error)
	GetKeyName() string
}
