package starsystem

type StarSystem struct {
	Stars          []*Star          `json:"stars,omitempty"`            //maximum 4
	RoguePlanets   []*RoguePlanet   `json:"rogue_planets,omitempty"`    //maximum 5
	RogueGasGiants []*RogueGasGiant `json:"rogue_gas_giants,omitempty"` //maximum 2
	Commets        []*Commet        `json:"commets,omitempty"`          //maximum 8
	DataField      string           `json:"data_field,omitempty"`
}

type Star struct {
	Stars     []*Star     `json:"stars,omitempty"`     //maximum 1
	Planets   []*Planet   `json:"planets,omitempty"`   //maximum 15
	Belts     []*Belt     `json:"belts,omitempty"`     //maximum 3
	GasGiant  []*GasGiant `json:"gas_giant,omitempty"` //maximum 6
	DataField string      `json:"data_field,omitempty"`
}

type CompanionStar struct {
	Star
	Parent string `json:"parent,omitempty"`
}

type Planet struct {
	Moons     []*Satelite `json:"moons,omitempty"` //maximum 3
	DataField string      `json:"data_field,omitempty"`
}

type Belt struct {
	Significant   []*Satelite `json:"significant,omitempty"`   //maximum 3
	Insignificant []*Asteroid `json:"insignificant,omitempty"` //maximum 12
	DataField     string      `json:"data_field,omitempty"`
}

type GasGiant struct {
	Moons     []*Satelite `json:"moons,omitempty"` //maximum 6
	DataField string      `json:"data_field,omitempty"`
}

type Satelite struct {
	Planet
	Parent string `json:"parent,omitempty"`
}

type Asteroid struct {
}

type Commet struct {
}

type RoguePlanet struct {
	Planet
}

type RogueGasGiant struct {
	GasGiant
}
