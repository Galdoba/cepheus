package survey

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/travellermap"
)

type SpaceHex struct {
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

func Localize(imported travellermap.WorldData) (*SpaceHex, error) {
	local := SpaceHex{}
	local.Name = imported.Name
	local.Hex = imported.Hex
	local.UWP = imported.UWP
	local.PBG = imported.PBG
	local.Zone = imported.Zone
	local.Bases = imported.Bases
	local.Allegiance = imported.Allegiance
	local.Stellar = imported.Stellar
	local.SS = imported.SS
	local.Ix = imported.Ix
	local.Ex = imported.Ex
	local.Cx = imported.Cx
	local.Nobility = imported.Nobility
	local.Worlds = imported.Worlds
	local.ResourceUnits = imported.ResourceUnits
	local.Subsector = imported.Subsector
	local.Quadrant = imported.Quadrant
	local.WorldX = imported.WorldX
	local.WorldY = imported.WorldY
	local.Remarks = imported.Remarks
	local.LegacyBaseCode = imported.LegacyBaseCode
	local.Sector = imported.Sector
	local.SubsectorName = imported.SubsectorName
	local.SectorAbbreviation = imported.SectorAbbreviation
	local.AllegianceName = imported.AllegianceName
	return &local, nil
}

func (sh *SpaceHex) Key() string {
	return fmt.Sprintf("%s/%s (%s: %s) {%v;%v}", sh.Name, sh.Sector, sh.Hex, sh.UWP, sh.WorldX, sh.WorldY)
}
