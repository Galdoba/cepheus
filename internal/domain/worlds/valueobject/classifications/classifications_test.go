package classifications

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func TestClassify(t *testing.T) {
	u, err := uwp.New("B100727-C")
	fmt.Println(err)
	cls := Classify(u, WithCodesRequested("Ab"))
	tc := TradeCodes(u)
	fmt.Println(cls)
	fmt.Println(tc)
}
