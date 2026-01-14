package travellermap

type WorldList struct {
	Worlds []WorldData `json:"worlds"`
}

type SpaceMap struct {
	Worlds map[string]WorldData
}

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
