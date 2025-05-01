package config

type Config struct {
	EnableDebug bool // デバッグモードの有効化
}

func NewConfig() *Config {
	return &Config{
		EnableDebug: false, // デフォルトは無効
	}
}
