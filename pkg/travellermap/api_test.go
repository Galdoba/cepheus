package travellermap

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"
)

func TestGet(t *testing.T) {
	// worlds, err := GetWorldData(coordinates.NewSpaceCoordinates(-107, -17), 12)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for _, w := range worlds {
	// 	fmt.Println(w)
	// }
	// FullMapUpdate("/home/galdoba/worlds3.json")
	//
	db := Database{}
	db.Worlds = make(map[string]WorldData)
	bt, _ := os.ReadFile("/home/galdoba/worlds3.json")
	json.Unmarshal(bt, &db)
	fmt.Println(len(db.Worlds))
	keys := []string{}
	for k := range db.Worlds {
		keys = append(keys, k)

	}
	slices.Sort(keys)
	db2 := Database{}
	db2.Worlds = make(map[string]WorldData)
	for _, k := range keys {
		db2.Worlds[k] = db.Worlds[k]
	}
	f, _ := os.Create("/home/galdoba/worlds_condenced.json")
	data, _ := json.Marshal(db2)
	f.Write(data)
	//
	// cpN := calibrationPoints(13)
	// maxX := 0
	// maxY := 0
	// minX := 0
	// minY := 0
	// for i, cp := range cpN {
	// 	fmt.Println(i, cp)
	// 	x, y := cp.GlobalValues()
	// 	if x > maxX {
	// 		maxX = x
	// 	}
	// 	if x < minX {
	// 		minX = x
	// 	}
	// 	if y > maxY {
	// 		maxY = y
	// 	}
	// 	if y < minY {
	// 		minY = y
	// 	}
	// }
	// fmt.Println(maxX, maxY)
	// fmt.Println(minX, minY)
}
