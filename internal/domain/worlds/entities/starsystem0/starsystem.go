package starsystem0

import "github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/orbit"

type StarSystem struct {
	Sector          string                 `json:"sector"`
	Grid            string                 `json:"grid"`
	IISSDesignation string                 `json:"iiss_designation"`
	Location        string                 `json:"location"`
	InitialSurvey   string                 `json:"initial_survey"`
	LastUpdated     string                 `json:"last_updated"`
	SystemAge       float64                `json:"system_age"`
	TravelZone      string                 `json:"travel_zone"`
	Comments        string                 `json:"comments"`
	Stars           map[orbit.Orbit]*Star  `json:"stars"`
	Worlds          map[orbit.Orbit]*World `json:"worlds"`
	SurveyClass     int                    `json:"survey_class"`
	PrimaryStar     *Star                  `json:"primary_star"`
	Mainworld       *World                 `json:"mainworld"`
}

type Star struct {
	Type         string   `json:"type"`
	SubType      string   `json:"sub_type"`
	Class        string   `json:"class"`
	Mass         float64  `json:"mass"`
	Temperature  int      `json:"temperature"`
	Diameter     float64  `json:"diameter"`
	Luminocity   float64  `json:"luminocity"`
	Orbit        float64  `json:"orbit"`
	AU           float64  `json:"au"`
	Eccentricity float64  `json:"eccentricity"`
	Period       float64  `json:"period"`
	HZCO         float64  `json:"hzco"`
	MAO          float64  `json:"mao"`
	IsPrimary    bool     `json:"is_primary"`
	Notes        []string `json:"notes"`
}

type World struct {
	Name            string                `json:"name"`
	UWP             string                `json:"uwp"`
	Orbit           float64               `json:"orbit"`
	Eccentricity    float64               `json:"eccentricity"`
	Period          float64               `json:"period"`
	Type            string                `json:"type"`
	Primary         orbit.Orbit           `json:"primary"`
	BeltComposition *BeltComposition      `json:"belt_composition"`
	Size            *Size                 `json:"size"`
	Atmosphere      *Atmosphere           `json:"atmosphere"`
	Hydrographics   *Hydrographics        `json:"hydrographics"`
	Rotation        *Rotation             `json:"rotation"`
	Temperature     *Temperature          `json:"temperature"`
	Life            *Life                 `json:"life"`
	Resources       *Resources            `json:"resources"`
	Habitability    *Habitability         `json:"habitability"`
	Population      *Population           `json:"population"`
	Government      *Government           `json:"government"`
	LawLevel        *LawLevel             `json:"law_level"`
	Technology      *Technology           `json:"technology"`
	Culture         *Culture              `json:"culture"`
	Economics       *Economics            `json:"economics"`
	Starport        *Starport             `json:"starport"`
	Military        *Military             `json:"military"`
	Satellites      map[orbit.Orbit]World `json:"satellites"`
	IsMain          bool                  `json:"is_main"`
	Notes           []string              `json:"notes"`
}

type Size struct {
	Code           string  `json:"code"`
	Composition    string  `json:"composition"`
	Diameter       string  `json:"diameter"`
	Density        string  `json:"density"`
	Gravity        float64 `json:"gravity"`
	Mass           float64 `json:"mass"`
	EscapeVelocity float64 `json:"escape_velocity"`
}

type Atmosphere struct {
	Pressure    float64  `json:"pressure"`
	Composition string   `json:"composition"`
	OxygenBar   float64  `json:"oxygen_bar"`
	Taints      string   `json:"taints"`
	ScaleHeight float64  `json:"scale_height"`
	Notes       []string `json:"notes"`
}

type Hydrographics struct {
	Coverage     float64  `json:"coverage"`
	Composition  string   `json:"composition"`
	Distribution string   `json:"distribution"`
	MaiorBodies  string   `json:"maior_bodies"`
	MinorBodies  string   `json:"minor_bodies"`
	Other        string   `json:"other"`
	Notes        []string `json:"notes"`
}
type Rotation struct{}
type Temperature struct{}
type Life struct{}

type Resources struct {
	Rating float64  `json:"rating"`
	Notes  []string `json:"notes"`
}

type Habitability struct{}
type Population struct{}
type Government struct{}
type LawLevel struct{}
type Technology struct{}
type Culture struct{}
type Economics struct{}
type Starport struct{}
type Military struct{}

type BeltComposition struct {
	BeltSpan float64 `json:"belt_span"`
	Mtype    float64 `json:"mtype"`
	Stype    float64 `json:"stype"`
	Ctype    float64 `json:"ctype"`
	Other    float64 `json:"other"`
	Bulk     int     `json:"bulk"`
}
