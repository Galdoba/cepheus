package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/appcontext/jsonstore"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/Galdoba/cepheus/iiss/survey"
	"github.com/Galdoba/cepheus/internal/domain/generic/services/travellermap"
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
	if err := os.MkdirAll(filepath.Dir(dbm.db.Path()), 0755); err != nil {
		return nil, fmt.Errorf("map file directorycreation failed: %v", err)
	}
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

func Open(path string) (*jsonstore.JsonDB[survey.SpaceHex], error) {
	storage, err := jsonstore.Load[survey.SpaceHex](path)
	if err != nil {
		return nil, fmt.Errorf("jstor load: %v", err)
	}
	return storage, nil
}

type searchOption struct {
	caseSensitive bool
	maxValues     int
	random        bool
}

type SearchOption func(*searchOption)

func CaseSensitive(cs bool) SearchOption {
	return func(so *searchOption) {
		so.caseSensitive = cs
	}
}

func Search(db *jsonstore.JsonDB[survey.SpaceHex], key string, opts ...SearchOption) []*survey.SpaceHex {
	so := searchOption{
		caseSensitive: false,
	}
	for _, modify := range opts {
		modify(&so)
	}
	entries, err := db.GetAll()
	if err != nil {
		panic(err)
	}
	found := []*survey.SpaceHex{}
	for world_key, worldData := range entries {
		switch so.caseSensitive {
		case true:
			if strings.Contains(world_key, key) {
				found = append(found, worldData)
			}
		case false:
			if strings.Contains(strings.ToLower(world_key), strings.ToLower(key)) {
				found = append(found, worldData)
			}

		}
	}
	return found
}
