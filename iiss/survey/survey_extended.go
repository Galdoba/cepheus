package survey

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/travellermap"
)

type SpaceHex struct {
	Sector string
	Hex    string
	Name   string
	MW_UWP string
	X      string
	Y      string
}

func Localize(imported travellermap.WorldData) (*SpaceHex, error) {
	local := SpaceHex{}
	local.Name = imported.Name
	local.Hex = imported.Hex
	local.Sector = imported.Sector
	local.MW_UWP = imported.UWP
	local.X = fmt.Sprintf("%v", imported.WorldX)
	local.Y = fmt.Sprintf("%v", imported.WorldY)

	return &local, nil
}

func (sh *SpaceHex) Key() string {
	return fmt.Sprintf("%s/%s (%s: %s) [%v;%v]", sh.Name, sh.Sector, sh.Hex, sh.MW_UWP, sh.X, sh.Y)
}
