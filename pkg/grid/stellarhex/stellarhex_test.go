package stellarhex

import (
	"fmt"
	"testing"
	"time"

	"github.com/Galdoba/t5/internal/grid/coordinates"
	"github.com/Galdoba/t5/pkg/dice"
	"github.com/Galdoba/t5/pkg/ehex"
)

func TestNew(t *testing.T) {
	star := 0

	for sc := 0; sc <= 65; sc++ {
		for ss := 0; ss <= 15; ss++ {
			for x := 10; x <= 18; x++ {
				for y := 10; y <= 20; y++ {
					sh := New(coordinates.NewSpaceCoordinates(x, y), Density(Density_Standard))
					err := sh.GenerateMissingDetails(dice.NewDicepool())
					if err != nil {
						t.Errorf("err: %v", err)
						panic(0)
					}
					fmt.Printf("star %v: sector %v", star, sc)
					fmt.Print(ehex.FromInt(ss+10).Code(), " ")
					fmt.Print(sh)
					if sh.StarSystem == "+" {
						fmt.Print(" Normal system")
					}
					if sh.BH == "+" {
						fmt.Print(" Black Hole")
					}
					if sh.NS == "+" {
						fmt.Print(" Neitron Star")
					}
					if sh.BD == "+" {
						fmt.Print(" hidden Brown dwarf")
					}
					if sh.D == "+" {
						fmt.Print(" hidden White dwarf")
					}
					fmt.Println("")
					time.Sleep(time.Millisecond * 1)
					star++
				}
			}
		}
	}
}
