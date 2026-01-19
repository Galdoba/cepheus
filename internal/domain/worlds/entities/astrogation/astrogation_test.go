package astrogation

import (
	"fmt"
	"testing"
)

func TestAstrogationBasic_JumpMap(t *testing.T) {
	as, err := New(`c:\Users\pemaltynov\.local\share\trv_worlds\data\second_survey_data.json`)
	if err != nil {
		fmt.Println(err)
		return
	}
	source, err := as.World("{-107,-17}")
	worlds := as.JumpMap(source, 6)
	for _, w := range worlds {
		fmt.Println(w.DatabaseKey(), w)
	}
	fmt.Println("cloasest navy base:", as.ClosestNavyBase(source.Coordinates()))
}
