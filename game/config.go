package game

type Config struct {
	EnableDebug bool // Enable or disable debug mode
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		EnableDebug: false, // Debug mode is disabled by default
	}
}
