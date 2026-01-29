package filepaths

import (
	"path/filepath"

	"github.com/Galdoba/appcontext/xdg"
	"github.com/Galdoba/cepheus/internal/declare"
)

const (
	ImportDataFile = "second_survey_data.json"
	WorldsDataFile = "worlds_data.json"
)

func DefaultExternalDB_File() string {
	return xdg.Location(xdg.ForData(),
		xdg.WithProgramName(declare.APP_TRV_WORLDS),
		xdg.WithFileName(ImportDataFile))
}

func DefaultDerivedDB_File() string {
	return xdg.Location(xdg.ForData(),
		xdg.WithProgramName(declare.APP_TRV_WORLDS),
		xdg.WithFileName(WorldsDataFile))
}

func RandomTablesDirectory() string {
	base := xdg.Location(xdg.ForData(), xdg.WithProgramName(declare.APP_TRV_WORLDS))
	return filepath.Join(base, "random_tables")
}
