package astrogation

import (
	"fmt"
	"strings"

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

type AstrogationBasic struct {
	id string
	db *jsonstorage.Storage[t5ss.WorldData]
}

func New(db_path string) (*AstrogationBasic, error) {
	as := AstrogationBasic{
		id: "generic",
	}
	db, err := jsonstorage.OpenStorage[t5ss.WorldData](db_path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	as.db = db

	return &as, nil
}

func Coordinates(entry t5ss.WorldData) coordinates.Cube {
	return coordinates.NewGlobal(entry.WorldX, entry.WorldY).ToCube()
}

func (as *AstrogationBasic) Distance(source t5ss.WorldData, target t5ss.WorldData) int {
	from := Coordinates(source)
	to := Coordinates(target)
	return coordinates.Distance(from, to)
}

func (as *AstrogationBasic) JumpMap(source t5ss.WorldData, radius int) []t5ss.WorldData {
	worlds := []t5ss.WorldData{}
	center := Coordinates(source)
	coords := coordinates.Spiral(center, radius)
	for _, crd := range coords {
		gc := crd.ToGlobal()
		key := fmt.Sprintf("%v", gc.DatabaseKey())
		if wd, err := as.db.Read(key); err == nil {
			worlds = append(worlds, wd)
		} else {
			fmt.Println(err)
		}
	}
	return worlds

}

func (as *AstrogationBasic) World(key string) (t5ss.WorldData, error) {
	return as.db.Read(key)
}

func (as *AstrogationBasic) ClosestNavyBase(crd coordinates.Global) int {
	for i := range 300 {
		coords := coordinates.Ring(crd.ToCube(), i)
		for _, c := range coords {
			if wd, err := as.db.Read(c.ToGlobal().DatabaseKey()); err == nil {
				if strings.Contains(wd.Bases, "N") || strings.Contains(wd.Bases, "D") {
					return i
				}
			}
		}
	}
	return -1
}

// func (as *AstrogationBasic) TradeRoutes(source )
