package config

import "github.com/Galdoba/cepheus/internal/domain/support/entities/paths"

type TrvWorldsCfg struct {
	Import TrvWorldsImport `toml:"import" comment:"import command configurations"`
	World  Worlds          `toml:"world" comment:"worlds data configuration"`
}

type TrvWorldsImport struct {
	ImportDataPath        string `toml:"import_data_path"`
	DownloadRetryAttempts int    `toml:"download_retry_attempts"`
	CoordinatesRingSize   int    `toml:"coordinates_ring_size"`
}

type Worlds struct {
	WorldsDataPath string `toml:"worlds_data_path"`
}

func DefaultTrvWorldsConfig() TrvWorldsCfg {
	return TrvWorldsCfg{
		Import: TrvWorldsImport{
			ImportDataPath:        paths.ImportStoragePath(),
			DownloadRetryAttempts: 3,
			CoordinatesRingSize:   2,
		},
		World: Worlds{
			WorldsDataPath: paths.WorldsStoragePath(),
		},
	}
}
