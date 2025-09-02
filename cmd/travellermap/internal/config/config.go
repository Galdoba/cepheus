package config

// Config is a travellermap configuration
type Config struct {
	App    AppData `toml:"travellermap"`
	Logger Logger  `toml:"logger"`
}

type AppData struct {
	Version string `toml:"version"`
}

type Logger struct {
	Enabled bool   `toml:"enabled"`
	Level   string `toml:"logging_level"`
	Color   bool   `toml:"color_output"`
}

func New() *Config {
	cfg := Config{
		App: AppData{
			Version: "0.0.x",
		},
		Logger: Logger{
			Enabled: false,
			Level:   "info",
			Color:   false,
		},
	}
	return &cfg
}
