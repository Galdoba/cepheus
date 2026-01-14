package astrogation

import (
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
)

type SpaceHex struct {
	crd              coordinates.Cube
	jumpspaceMod_IN  int
	jumpspaceMod_OUT int
	hasGasGigant     bool
	attraction       int
}
