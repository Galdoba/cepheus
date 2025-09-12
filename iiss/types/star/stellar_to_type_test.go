package star

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/travellermap"
)

func TestFromStellar(t *testing.T) {
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
			// fmt.Printf("%v\t%v\n", i, w.Stellar)
			stars := FromStellar(w.Stellar)
			for k, s := range stars {
				if k > 0 {
					continue
				}
				st, err := Generate(dice.NewDicepool(), KnownStellar(s.str))
				if err != nil {
					fmt.Println("bad gen")
					panic(err)
				}
				fmt.Println(i, k, st.String())
				if w.Stellar == "" {
					fmt.Println("created")
				}

			}

			i++
		}

	}

}
