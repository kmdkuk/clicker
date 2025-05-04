package config

type Config struct {
	EnableDebug bool   // Enable or disable debug mode
	SaveKey     string // Key for saving game state
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		EnableDebug: false, // Debug mode is disabled by default
		SaveKey:     DefaultSaveKey,
	}
}

const (
	DefaultSaveKey string = "game_state.json"
	CostMultiplier        = 1.15
)
