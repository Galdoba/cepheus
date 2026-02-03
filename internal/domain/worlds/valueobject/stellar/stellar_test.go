package stellar

import (
	"fmt"
	"testing"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
	"github.com/Galdoba/cepheus/pkg/dice"
)

func List() {
	db, err := jsonstorage.OpenStorage[t5ss.WorldData](`/home/galdoba/.local/share/trv_worlds/second_survey_data.json`)
	if err != nil {
		panic(err)
	}
	list := make(map[string]int)
	for _, e := range db.ReadAll() {
		// fmt.Println(e.Stellar)
		separated := ExtractStars(e.Stellar)
		fmt.Println(Stellar(e.Stellar))
		fmt.Println(Stellar(e.Stellar).Split())
		for _, str := range separated {
			// fmt.Printf("star='%v':%v", str, validateStar(str))
			if !validateStar(str) {

				fmt.Println(e)
				list[e.Stellar]++
			}
		}
	}
	max := len(list)
	fmt.Println(len(list))
	for max > 0 {
		for k, v := range list {
			if v == max {
				fmt.Println(k, v)
			}
		}
		max--
	}
	r := dice.New("")
	for i := 0; i < 20; i++ {
		fmt.Println(RollDesignations(r))
		time.Sleep(time.Second)
	}

}

func TestList(t *testing.T) {
	List()
}
