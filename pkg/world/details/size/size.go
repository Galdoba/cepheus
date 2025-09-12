package size

import "github.com/Galdoba/cepheus/pkg/dice"

type SizeDetails struct {
	Code            string  `json:"code,omitempty"`
	Diameter        int     `json:"diameter,omitempty"`
	Dencity         float64 `json:"dencity,omitempty"`
	Gravity         float64 `json:"gravity,omitempty"`
	Mass            float64 `json:"mass,omitempty"`
	EscapeVelocity  float64 `json:"escape_velocity,omitempty"`
	OrbitalVelocity float64 `json:"orbital_velocity,omitempty"`
	Composition     string  `json:"core_composition,omitempty"`
	worldType       string
	hzco            float64
	gyrs            float64
}

func GenerateDetails(dp *dice.Dicepool, code string) *SizeDetails {
	return &SizeDetails{}
}
