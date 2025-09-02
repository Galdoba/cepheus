package travellermap

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Galdoba/cepheus/pkg/grid/coordinates"
)

func TestGet(t *testing.T) {
	// tests := []struct {
	// 	name string // description of this test case
	// 	// Named input parameters for target function.
	// 	url     string
	// 	want    []byte
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, gotErr := travellermap.Get(tt.url)
	// 		if gotErr != nil {
	// 			if !tt.wantErr {
	// 				t.Errorf("Get() failed: %v", gotErr)
	// 			}
	// 			return
	// 		}
	// 		if tt.wantErr {
	// 			t.Fatal("Get() succeeded unexpectedly")
	// 		}
	// 		// TODO: update the condition below to compare got with tt.want.
	// 		if true {
	// 			t.Errorf("Get() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
	fmt.Println("start test")
	data, err := Get("https://www.travellermap.com/data")
	fmt.Println(err)
	fmt.Println(len(data), "bytes received")
	sectors := SectorList{}
	err = json.Unmarshal(data, &sectors)
	fmt.Println(err)
	// for _, sect := range sectors.Sectors {
	// 	if !strings.Contains(sect.Tags, "OTU") {
	// 		continue
	// 	}
	// 	fmt.Println(sect)
	// 	// time.Sleep(time.Second)
	// }
	coord := coordinates.NewSpaceCoordinates(-107, -17)
	sx, sy, _, _ := coord.LocalValues()
	abb := ""
	for _, sect := range sectors.Sectors {
		if !(sect.X == sx && sect.Y == sy && sect.Milieu == "M1105") {
			continue
		}
		abb = sect.Abbreviation
	}
	url := fmt.Sprintf("https://www.travellermap.com/data/%v/%v/jump/1", abb, coord.SectorHex())
	data2, err := Get(url)
	worlds := WorldList{}
	err = json.Unmarshal(data2, &worlds)
	fmt.Println(err)
	for _, w := range worlds.Worlds {
		fmt.Println(w)
		time.Sleep(time.Second)
	}
}
