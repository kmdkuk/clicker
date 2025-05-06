package config

type Config struct {
	EnableDebug  bool   // Enable or disable debug mode
	SaveKey      string // Key for saving game state
	ScreenWidth  int    // Width of the game screen
	ScreenHeight int    // Height of the game screen
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		EnableDebug:  false, // Debug mode is disabled by default
		SaveKey:      DefaultSaveKey,
		ScreenWidth:  800,
		ScreenHeight: 600,
	}
}

const (
	DefaultSaveKey string = "game_state.json"
	CostMultiplier        = 1.15
)
