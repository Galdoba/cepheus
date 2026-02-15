package stellar

import (
	"testing"
)

// func List() {
// 	db, err := jsonstorage.OpenStorage[t5ss.WorldData](`/home/galdoba/.local/share/trv_worlds/second_survey_data.json`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	list := make(map[string]int)
// 	for _, e := range db.ReadAll() {
// 		// fmt.Println(e.Stellar)
// 		separated := ExtractStars(e.Stellar)
// 		fmt.Println(Stellar(e.Stellar))
// 		fmt.Println(Stellar(e.Stellar).Split())
// 		for _, str := range separated {
// 			// fmt.Printf("star='%v':%v", str, validateStar(str))
// 			if !validateStar(str) {

// 				fmt.Println(e)
// 				list[e.Stellar]++
// 			}
// 		}
// 	}
// 	max := len(list)
// 	fmt.Println(len(list))
// 	for max > 0 {
// 		for k, v := range list {
// 			if v == max {
// 				fmt.Println(k, v)
// 			}
// 		}
// 		max--
// 	}
// 	r := dice.New("")
// 	for i := 0; i < 20; i++ {
// 		fmt.Println(RollDesignations(r))
// 		time.Sleep(time.Second)
// 	}

// }

// func TestList(t *testing.T) {
// 	List()
// }

func Test_parseStellar(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		stellar string
		want    []string
	}{
		{
			name:    "simple",
			stellar: "G2 V",
			want:    []string{"G2 V"},
		},
		{
			name:    "simple 2",
			stellar: "BD",
			want:    []string{"BD"},
		},
		{
			name:    "complex 1",
			stellar: "Anomaly with BD and M5 III",
			want:    []string{"BD", "M5 III"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStellar(tt.stellar)
			// TODO: update the condition below to compare got with tt.want.
			if !match(tt.want, got) {
				t.Errorf("parseStellar() = %v, want %v", got, tt.want)
			}
		})
	}

}

func match(sl1, sl2 []string) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for i := range sl1 {
		if sl1[i] != sl2[i] {
			return false
		}
	}
	return true
}
