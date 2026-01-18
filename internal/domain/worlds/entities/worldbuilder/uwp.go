package worldbuilder

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func GenerateMissingFields_UWP(wb *WorldBuilder, u uwp.UWP) (uwp.UWP, error) {
	populated, err := uwp.Populate(u, wb.dice)
	if err != nil {
		return u, fmt.Errorf("failed to populate old uwp '%v': %v", string(u), err)
	}
	return populated, nil
}
