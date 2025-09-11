package size

import "github.com/Galdoba/cepheus/pkg/dice"

type SizeDetails struct {
	Code            string
	Diameter        int
	Dencity         float64
	Gravity         float64
	Mass            float64
	EscapeVelocity  float64
	OrbitalVelocity float64
	Composition     string
	hzco            float64
	gyrs            float64
}

func GenerateDetails(dp *dice.Dicepool, code string, opts ...SizeDetailOption) *SizeDetails {

}

type SizeDetailOption func(*SizeDetails)
