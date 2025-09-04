package config

// Config is a travellermap configuration
type Config struct {
	App    AppData `toml:"travellermap"`
	Logger Logger  `toml:"logger"`
	Files  Files   `toml:"files"`
}

type AppData struct {
	Version string `toml:"version"`
}

type Logger struct {
	Enabled       bool   `toml:"enabled"`
	Filepath      string `toml:"path"`
	Level         string `toml:"logging_level"`
	ConsoleOutput bool   `toml:"console_logging"`
	Color         bool   `toml:"color_output"`
}

type Files struct {
	DataDirectory string `toml:"canonical_data_directory"`
	WorkSpaces    string `toml:"workspaces"`
}
