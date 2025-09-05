package survey

// ImportFormat is a representation of site data for starsystem. It contains all Traveller5 Second Survey Data plus coordinates data and description.
type ImportFormat struct {
	Name               string `json:"Name,omitempty"`
	Hex                string `json:"Hex,omitempty"`
	UWP                string `json:"UWP,omitempty"`
	PBG                string `json:"PBG,omitempty"`
	Zone               string `json:"Zone,omitempty"`
	Bases              string `json:"Bases,omitempty"`
	Allegiance         string `json:"Allegiance,omitempty"`
	Stellar            string `json:"Stellar,omitempty"`
	SS                 string `json:"SS,omitempty"`
	Ix                 string `json:"Ix,omitempty"`
	Ex                 string `json:"Ex,omitempty"`
	Cx                 string `json:"Cx,omitempty"`
	Nobility           string `json:"Nobility,omitempty"`
	Worlds             int    `json:"Worlds,omitempty"`
	ResourceUnits      int    `json:"ResourceUnits,omitempty"`
	Subsector          int    `json:"Subsector,omitempty"`
	Quadrant           int    `json:"Quadrant,omitempty"`
	WorldX             int    `json:"WorldX,omitempty"`
	WorldY             int    `json:"WorldY,omitempty"`
	Remarks            string `json:"Remarks,omitempty"`
	LegacyBaseCode     string `json:"LegacyBaseCode,omitempty"`
	Sector             string `json:"Sector,omitempty"`
	SubsectorName      string `json:"SubsectorName,omitempty"`
	SectorAbbreviation string `json:"SectorAbbreviation,omitempty"`
	AllegianceName     string `json:"AllegianceName,omitempty"`
}
