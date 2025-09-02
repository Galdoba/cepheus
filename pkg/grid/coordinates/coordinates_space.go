package coordinates

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/grid/coordinates/convert"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/cube"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/global"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/local"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/sector"
)

type SpaceCoordinates struct {
	cube   cube.Cube
	global global.SpaceGlobal
	local  local.SpaceSectorLocal
}

type XY struct {
	X int
	Y int
}

type SectorBridge struct {
	Sectors map[XY]Sector
}

type Sector struct {
	X    int
	Y    int
	Name string
	Abb  string
}

// NewSpaceCoordinates - create space hex coordinates based on values provided:
//
// 2 values: GlobalWorldCoordinates (X, Y)
// 3 values: Cube Coordinates (Q,R,S)
// 4 values: Local Sector Coordinates (SectorX, SectorY, LocalX, LocalY)
func NewSpaceCoordinates(values ...int) SpaceCoordinates {
	sc := SpaceCoordinates{}
	switch len(values) {
	case 2:
		sc.global = global.NewSpaceGlobal(values[0], values[1])
		sc.cube = convert.GlobalToCube(sc.global)
		sc.local = convert.GlobalToLocal(sc.global)
	case 3:
		sc.cube = cube.NewCube(values[0], values[1], values[2])
		sc.global = convert.CubeToGlobal(sc.cube)
		sc.local = convert.CubeToLocal(sc.cube)
	case 4:
		sc.local = local.NewLocal(values[0], values[1], values[2], values[3])
		sc.cube = convert.LocalToCube(sc.local)
		sc.global = convert.LocalToGlobal(sc.local)
	default:
		panic(fmt.Sprintf("unsupported values quantity (%v)", len(values)))
	}

	return sc
}

func (sc SpaceCoordinates) StringSectorNameHex() string {
	return fmt.Sprintf("%v %v", sector.Name(sc.local.SectorX, sc.local.SectorY), sector.Hex(sc.local.X, sc.local.Y))
}

func (sc SpaceCoordinates) hexValues() (int, int, int) {
	return sc.cube.Q, sc.cube.R, sc.cube.S
}

func (sc SpaceCoordinates) SectorHex() string {
	w := fmt.Sprintf("%v", sc.local.X)
	if len(w) < 2 {
		w = "0" + w
	}
	h := fmt.Sprintf("%v", sc.local.Y)
	if len(h) < 2 {
		h = "0" + h
	}
	return w + h
}

func (sc SpaceCoordinates) CubeValues() (int, int, int) {
	return sc.cube.Q, sc.cube.R, sc.cube.S
}

func (sc SpaceCoordinates) GlobalValues() (int, int) {
	return sc.global.X, sc.global.Y
}

func (sc SpaceCoordinates) LocalValues() (int, int, int, int) {
	return sc.local.SectorX, sc.local.SectorY, sc.local.X, sc.local.Y
}

func (sc SpaceCoordinates) Validate() error {
	if err := convert.RoundTrip(sc.hexValues()); err != nil {
		return err
	}
	return nil
}
