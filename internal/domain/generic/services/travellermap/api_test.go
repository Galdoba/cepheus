package travellermap

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"slices"
// 	"testing"
// )

// func TestGet(t *testing.T) {
// 	FullMapUpdate("/home/galdoba/worlds3.json")
// 	db := Database{}
// 	db.Worlds = make(map[string]WorldData)
// 	bt, _ := os.ReadFile("/home/galdoba/worlds3.json")
// 	json.Unmarshal(bt, &db)
// 	fmt.Println(len(db.Worlds))
// 	keys := []string{}
// 	for k := range db.Worlds {
// 		keys = append(keys, k)

// 	}
// 	slices.Sort(keys)
// 	db2 := Database{}
// 	db2.Worlds = make(map[string]WorldData)
// 	for _, k := range keys {
// 		db2.Worlds[k] = db.Worlds[k]
// 	}
// 	f, _ := os.Create("/home/galdoba/worlds_condenced.json")
// 	data, _ := json.Marshal(db2)
// 	f.Write(data)
// }
