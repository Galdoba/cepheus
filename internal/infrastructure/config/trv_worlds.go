package config

import "github.com/Galdoba/cepheus/internal/domain/support/entities/paths"

type TrvWorldsCfg struct {
	Import TrvWorldsImport `toml:"import" comment:"import command configurations"`
}

type TrvWorldsImport struct {
	ImportDataPath        string `toml:"import_data_path"`
	DownloadRetryAttempts int    `toml:"download_retry_attempts"`
	CoordinatesRingSize   int    `toml:"coordinates_ring_size"`
}

func DefaultTrvWorldsConfig() TrvWorldsCfg {
	return TrvWorldsCfg{
		Import: TrvWorldsImport{
			ImportDataPath:        paths.ImportStoragePath(),
			DownloadRetryAttempts: 3,
			CoordinatesRingSize:   2,
		},
	}
}
