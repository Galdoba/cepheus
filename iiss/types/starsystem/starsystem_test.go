package starsystem

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Galdoba/cepheus/pkg/travellermap"
)

func TestPositions(t *testing.T) {
	data, err := os.ReadFile(`c:\Users\pemaltynov\.local\share\travellermap\data\default.json`)
	if err != nil {
		fmt.Println(err)
		return
	}
	wl := travellermap.SpaceMap{}
	wl.Worlds = make(map[string]travellermap.WorldData)
	fmt.Println(json.Unmarshal(data, &wl))
	i := 1
	for _, w := range wl.Worlds {
		switch w.Stellar {
		default:
			fmt.Printf("%v\t%v\n", i, w.Stellar)
			ssg := NewGenerator(WithStellar(w.Stellar))
			ss := ssg.GenerateSystem()
			fmt.Println("")
			for k, v := range ss.Stars {
				fmt.Println(k, v)
			}
			for k, v := range ss.Orbits {
				fmt.Println(k, v)
			}
			//time.Sleep(time.Second)

			i++
		}

	}
}
