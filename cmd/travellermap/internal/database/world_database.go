package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/appcontext/jsonstore"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/Galdoba/cepheus/iiss/survey"
	"github.com/Galdoba/cepheus/pkg/travellermap"
)

type dbManager struct {
	canonical string
	session   string
	db        *jsonstore.JsonDB[survey.SpaceHex]
}

var dbm = &dbManager{}

func Create(actx *infra.Container, canonicalData, newMap string) (*jsonstore.JsonDB[survey.SpaceHex], error) {
	dbm = &dbManager{
		canonical: canonicalData,
		session:   newMap,
	}
	err := fmt.Errorf("db not opened")
	imported, err := Canonical(actx, canonicalData)
	if err != nil {
		return nil, fmt.Errorf("canonical: %v", err)
	}
	lenCan := len(imported.Worlds)
	fmt.Printf("canonical data: %v entries\n", lenCan)
	sessionMap := filepath.Join(actx.Config.Files.WorkSpaces, newMap+".json")
	storage, err := jsonstore.New[survey.SpaceHex](sessionMap)
	if err != nil {
		return nil, fmt.Errorf("jstor creation: %v", err)
	}
	dbm.db = storage
	wnum := 1
	for _, imp := range imported.Worlds {
		extended, err := survey.Localize(imp)
		if err != nil {
			return nil, fmt.Errorf("data localization failed: %v", err)
		}
		dbm.db.Insert(extended.Key(), extended)
		fmt.Printf("data localization: %v/%v \r", wnum, lenCan)
		wnum++
	}
	fmt.Println("")
	if err := dbm.db.Save(); err != nil {
		return nil, fmt.Errorf("map file saving failed: %v", err)
	}
	return dbm.db, nil
}

func Canonical(actx *infra.Container, key string) (travellermap.Database, error) {
	file := filepath.Join(actx.Config.Files.DataDirectory, key+".json")
	imported := travellermap.Database{}
	imported.Worlds = make(map[string]travellermap.WorldData)
	data, err := os.ReadFile(file)
	if err != nil {
		return imported, fmt.Errorf("read file: %v", err)
	}
	err = json.Unmarshal(data, &imported)
	if err != nil {
		return imported, fmt.Errorf("unmarashal canonical data: %v", err)
	}
	return imported, nil

}
