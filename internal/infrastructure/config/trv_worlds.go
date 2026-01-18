package config

import "github.com/Galdoba/cepheus/internal/domain/support/entities/paths"

type TrvWorldsCfg struct {
	Import TrvWorldsImport `toml:"import" comment:"import command configurations"`
	World  Worlds          `toml:"world" comment:"worlds data configuration"`
}

type TrvWorldsImport struct {
	External_DB_File      string `toml:"external_db_file"`
	DownloadRetryAttempts int    `toml:"download_retry_attempts"`
	CoordinatesRingSize   int    `toml:"coordinates_ring_size"`
}

type Worlds struct {
	Derived_DB_File string `toml:"dirived_db_file"`
}

func DefaultTrvWorldsConfig() TrvWorldsCfg {
	return TrvWorldsCfg{
		Import: TrvWorldsImport{
			External_DB_File:      paths.DefaultExternalDB_File(),
			DownloadRetryAttempts: 3,
			CoordinatesRingSize:   2,
		},
		World: Worlds{
			Derived_DB_File: paths.DefaultDerivedDB_File(),
		},
	}
}
