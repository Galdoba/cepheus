package t5ss

import (
	"slices"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/travelzone"
)

// Name               string `json:"Name,omitempty"`
// Hex                string `json:"Hex,omitempty"`
// UWP                string `json:"UWP,omitempty"`
// PBG                string `json:"PBG,omitempty"`
// Zone               string `json:"Zone,omitempty"`
// Bases              string `json:"Bases,omitempty"`
// Allegiance         string `json:"Allegiance,omitempty"`
// Stellar            string `json:"Stellar,omitempty"`
// SS                 string `json:"SS,omitempty"`
// Ix                 string `json:"Ix,omitempty"`
// Ex                 string `json:"Ex,omitempty"`
// Cx                 string `json:"Cx,omitempty"`
// Nobility           string `json:"Nobility,omitempty"`
// Worlds             int    `json:"Worlds,omitempty"`
// ResourceUnits      int    `json:"ResourceUnits"`
// Subsector          int    `json:"Subsector"`
// Quadrant           int    `json:"Quadrant"`
// WorldX             int    `json:"WorldX"`
// WorldY             int    `json:"WorldY"`
// Remarks            string `json:"Remarks"`
// LegacyBaseCode     string `json:"LegacyBaseCode,omitempty"`
// Sector             string `json:"Sector,omitempty"`
// SubsectorName      string `json:"SubsectorName,omitempty"`
// SectorAbbreviation string `json:"SectorAbbreviation,omitempty"`
// AllegianceName     string `json:"AllegianceName,omitempty"`

func (wd WorldData) Coordinates() coordinates.Global {
	return coordinates.NewGlobal(wd.WorldX, wd.WorldY)
}

func (wd WorldData) TradeCodes() []classifications.Classification {
	list := strings.Split(wd.Remarks, " ")
	cls := []classifications.Classification{}
	for _, code := range list {
		cl := classifications.Classification(code)
		if slices.Contains(classifications.All(), cl) {
			cls = append(cls, cl)
		}

	}
	return cls
}

func (wd WorldData) TravelZone() travelzone.Zone {
	if travelzone.Zone(wd.Zone) == travelzone.Amber {
		return travelzone.Amber
	}
	if travelzone.Zone(wd.Zone) == travelzone.Red {
		return travelzone.Red
	}
	return ""
}
