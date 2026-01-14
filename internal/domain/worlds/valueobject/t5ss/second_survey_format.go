package t5ss

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates"
)

type SpaceCoordinates string

type WorldData struct {
	Name               string `json:"name"`
	Hex                string `json:"hex"`
	UWP                string `json:"uwp"`
	PBG                string `json:"pbg"`
	Zone               string `json:"zone"`
	Bases              string `json:"bases"`
	Allegiance         string `json:"allegiance"`
	Stellar            string `json:"stellar"`
	SS                 string `json:"ss"`
	Ix                 string `json:"ix"`
	Ex                 string `json:"ex"`
	Cx                 string `json:"cx"`
	Nobility           string `json:"nobility"`
	Worlds             int    `json:"worlds"`
	ResourceUnits      int    `json:"resource_units"`
	Subsector          int    `json:"subsector"`
	Quadrant           int    `json:"quadrant"`
	WorldX             int    `json:"world_x"`
	WorldY             int    `json:"world_y"`
	Remarks            string `json:"remarks"`
	LegacyBaseCode     string `json:"legacy_base_code"`
	Sector             string `json:"sector"`
	SubsectorName      string `json:"subsector_name"`
	SectorAbbreviation string `json:"sector_abbreviation"`
	AllegianceName     string `json:"allegiance_name"`
}

func (w WorldData) Import_DB_Key() string {
	return fmt.Sprintf("{%v, %v}", w.WorldX, w.WorldY)
}

func (w WorldData) Details_DB_Key() string {
	return fmt.Sprintf("%v [%v/%v %v] %v", w.Name, w.SubsectorName, w.Sector, w.Hex, w.Import_DB_Key())
}

func (w WorldData) Coordinates() coordinates.SpaceCoordinates {
	sc := coordinates.NewSpaceCoordinates(w.WorldX, w.WorldY)
	return sc
}
