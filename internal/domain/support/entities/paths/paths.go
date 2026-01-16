package paths

import (
	"path/filepath"

	"github.com/Galdoba/appcontext/xdg"
	"github.com/Galdoba/cepheus/internal/declare"
)

const (
	ImportDataFile = "second_survey_data.json"
	WorldsDataFile = "worlds_data.json"
)

func ImportStoragePath() string {
	path := xdg.New(declare.APP_TRV_WORLDS)
	return filepath.Join(path.PersistentDataDir(), ImportDataFile)
}

func WorldsStoragePath() string {
	path := xdg.New(declare.APP_TRV_WORLDS)
	return filepath.Join(path.PersistentDataDir(), WorldsDataFile)
}
