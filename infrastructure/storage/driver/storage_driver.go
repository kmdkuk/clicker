package driver

type StorageDriver interface {
	SaveData(data []byte) error
	LoadData() ([]byte, error)
}
