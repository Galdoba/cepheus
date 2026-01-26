package importance

import (
	"github.com/Galdoba/cepheus/internal/domain/support/valueobject/uwp"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/classifications"
)

type Importance int

type ImportanceMaker interface {
	UWP() uwp.UWP
	TradeCodes() []classifications.Classification
	Bases() []string
}
