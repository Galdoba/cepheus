package t5ss

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
)

type WorldData struct {
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
	ResourceUnits      int    `json:"ResourceUnits"`
	Subsector          int    `json:"Subsector"`
	Quadrant           int    `json:"Quadrant"`
	WorldX             int    `json:"WorldX"`
	WorldY             int    `json:"WorldY"`
	Remarks            string `json:"Remarks"`
	LegacyBaseCode     string `json:"LegacyBaseCode,omitempty"`
	Sector             string `json:"Sector,omitempty"`
	SubsectorName      string `json:"SubsectorName,omitempty"`
	SectorAbbreviation string `json:"SectorAbbreviation,omitempty"`
	AllegianceName     string `json:"AllegianceName,omitempty"`
}

type WorldBatch struct {
	List []WorldData `json:"Worlds,omitempty"`
}

func (w WorldData) Import_DB_Key() string {
	return fmt.Sprintf("{%v, %v}", w.WorldX, w.WorldY)
}

func (w WorldData) Details_DB_Key() string {
	return fmt.Sprintf("%v [%v/%v %v] %v", w.Name, w.SubsectorName, w.Sector, w.Hex, w.Import_DB_Key())
}

func (w WorldData) Coordinates() coordinates.Global {
	return coordinates.NewGlobal(w.WorldX, w.WorldY)
}

func (w WorldData) NormalizeName() string {
	if w.Name != "" {
		return w.Name
	}
	nameParts := []string{}
	if w.SubsectorName != "" {
		nameParts = append(nameParts, w.SubsectorName)
	}
	if w.Hex != "" {
		part := w.Hex
		if len(w.UWP) > 3 {
			sah := ""
			for i, u := range strings.Split(w.UWP, "") {
				switch i {
				case 1, 2, 3:
					sah += u
				}
			}
			part += "-" + sah
		}
		nameParts = append(nameParts, part)
	}
	return strings.Join(nameParts, " ")
}

func (w WorldData) ConfirmedBases() []string {
	bases := []string{}
	if strings.Contains(w.UWP, "?") {
		return bases
	}
	for _, b := range strings.Split(w.Bases, "") {
		bases = append(bases, b)
	}
	if len(bases) == 0 {
		bases = append(bases, "no bases confirmed")
	}
	return bases
}
