package travellermap

type WorldList struct {
	Worlds []WorldData
}

type WorldData struct {
	Name               string
	Hex                string
	UWP                string
	PBG                string
	Zone               string
	Bases              string
	Allegiance         string
	Stellar            string
	SS                 string
	Ix                 string
	Ex                 string
	Cx                 string
	Nobility           string
	Worlds             int
	ResourceUnits      int
	Subsector          int
	Quadrant           int
	WorldX             int
	WorldY             int
	Remarks            string
	LegacyBaseCode     string
	Sector             string
	SubsectorName      string
	SectorAbbreviation string
	AllegianceName     string
}
