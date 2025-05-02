package state

type StorageDriver interface {
	SaveData(data []byte) error
	LoadData() ([]byte, error)
}
