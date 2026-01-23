package astrogation

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/support/entities/paths"
	"github.com/Galdoba/cepheus/internal/domain/worlds/aggregates/world"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
)

type SpaceHex struct {
	crd              coordinates.Cube
	jumpspaceMod_IN  int
	jumpspaceMod_OUT int
	hasGasGigant     int
	attraction       int
}

type Astrogation struct {
	id       string
	basic    *jsonstorage.Storage[t5ss.WorldData]
	surveyed *jsonstorage.Storage[world.WorldDTO]
}

type astroCfg struct {
	id          string
	external_DB string
	derived_DB  string
}

func defaultCfg() astroCfg {
	return astroCfg{
		id:          "generic",
		external_DB: paths.DefaultExternalDB_File(),
		derived_DB:  paths.DefaultDerivedDB_File(),
	}
}

func New() (*Astrogation, error) {
	cfg := defaultCfg()
	as := Astrogation{
		id: "generic",
	}
	db, err := jsonstorage.OpenStorage[t5ss.WorldData](cfg.external_DB)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	as.basic = db
	db2, err := jsonstorage.OpenStorage[world.WorldDTO](cfg.derived_DB)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	as.surveyed = db2

	return &as, nil
}

func Coordinates(entry t5ss.WorldData) coordinates.Cube {
	return coordinates.NewGlobal(entry.WorldX, entry.WorldY).ToCube()
}

func (as *Astrogation) DistanceBasic(source t5ss.WorldData, target t5ss.WorldData) int {
	from := Coordinates(source)
	to := Coordinates(target)
	return coordinates.Distance(from, to)
}

func (as *Astrogation) JumpMap(source t5ss.WorldData, radius int) []t5ss.WorldData {
	worlds := []t5ss.WorldData{}
	center := Coordinates(source)
	coords := coordinates.Spiral(center, radius)
	for _, crd := range coords {
		gc := crd.ToGlobal()
		key := fmt.Sprintf("%v", gc.DatabaseKey())
		if wd, err := as.basic.Read(key); err == nil {
			worlds = append(worlds, wd)
		} else {
			fmt.Println(err)
		}
	}
	return worlds

}

func (as *Astrogation) World(key string) (t5ss.WorldData, error) {
	return as.basic.Read(key)
}

func (as *Astrogation) ClosestNavyBase(crd coordinates.Global) int {
	for i := range 300 {
		coords := coordinates.Ring(crd.ToCube(), i)
		for _, c := range coords {
			if wd, err := as.basic.Read(c.ToGlobal().DatabaseKey()); err == nil {
				if strings.Contains(wd.Bases, "N") || strings.Contains(wd.Bases, "D") {
					return i
				}
			}
		}
	}
	return -1
}

// func (as *AstrogationBasic) TradeRoutes(source )
