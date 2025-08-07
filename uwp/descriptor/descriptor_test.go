package descriptor_test

import (
	"testing"

	"github.com/Galdoba/cepheus/uwp"
	"github.com/Galdoba/cepheus/uwp/descriptor"
)

func TestNew(t *testing.T) {
	ds := descriptor.New(`c:\Users\pemaltynov\go\src\github.com\Galdoba\cepheus\uwp\descriptor\example.json`)
	for _, category := range []string{
		uwp.Port,
		uwp.Size,
		uwp.Atmo,
		uwp.Hydr,
		uwp.Pops,
		uwp.Govr,
		uwp.Laws,
		uwp.TL,
	} {
		codes := []string{}
		switch category {
		case uwp.Port:
			codes = []string{"A", "B", "C", "D", "E", "X", "F", "G", "H", "Y", "?"}
		case uwp.Size:
			codes = []string{"0", "R", "S", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "?"}
		case uwp.Atmo:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "?"}
		case uwp.Hydr:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "?"}
		case uwp.Pops:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "E", "F", "?"}
		case uwp.Govr:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "W", "X", "?"}
		case uwp.Laws:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "S", "?"}
		case uwp.TL:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "?"}
		}
		for _, code := range codes {
			ds.AddDescription(category, code, uwp.Description(category, code)["en"])
		}

	}
	ds.Save()
}
