package starsystem

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/worlds/orbit"
)

// Star represents the physical and astronomical properties of a star, including its spectral
// classification, temperature, mass, luminosity, and key orbital zone boundaries
// (habitable zone, snow line, etc.) used for stellar characterization and exoplanetary system modeling.
type Star struct {
	Designation string `json:"designation"`

	Type string `json:"type"`

	//The stellar class (also called spectral type) is a letter code (O, B, A, F, G, K, M, L, T, Y)
	//that categorizes a star primarily by its effective temperature and color.
	//This field might contain 'BD' value as placeholder for L/T/Y classes until details are provided.
	StellarClass string `json:"stellar_class"`

	//NumericalSubClass is a digit from 0 to 9 that provides a finer temperature subdivision within
	//a given StellarClass (spectral type).
	NumeriacalSubClass string `json:"numeriacal_sub_class"`

	//The luminosity class indicates a star’s size and evolutionary stage. It is based on the width
	//of spectral lines (which correspond to surface gravity and pressure).
	// Common codes:
	// Ia – bright supergiant
	// Ib – supergiant
	// II – bright giant
	// III – normal giant
	// IV – subgiant
	// V – main sequence dwarf (the “hydrogen‑burning” stage)
	// VI – subdwarf (low luminosity, metal‑poor)
	// D – white dwarf
	LuminocityClass string `json:"luminocity_class"`

	//The effective surface temperature of the star measured in Kelvin. It determines the star’s color,
	//spectral features, and the wavelength distribution of its emitted radiation.
	Temperature_K int `json:"temperature_k"`

	//The mass of the star expressed in solar masses (SM). One solar mass (M☉) is
	//the mass of the Sun (~1.989×10³⁰ kg). This value strongly influences the star’s gravity,
	//internal pressure, fusion rate, and lifespan.
	Mass_SM float64 `json:"mass_sm"`

	//The total radiant power (energy output per second) of the star, expressed in solar units (SU).
	//One solar unit (L☉) equals the Sun’s luminosity (~3.828×10²⁶ W).
	//It determines how much energy planets receive.
	Luminosity_SU float64 `json:"luminosity_su"`

	//The inner boundary of the star’s planetary system or the region of interest for planet
	//formation/dynamics, measured in Astronomical Units (AU). One AU is the average
	//Earth–Sun distance (~149.6 million km). This limit might mark where temperatures become too high
	//for certain processes or where material cannot survive.
	InnerLimit_AU float64 `json:"inner_limit_au"`

	//The inner edge of the classical habitable zone (the “Goldilocks zone”) in AU. At this distance,
	//a rocky planet with a suitable atmosphere could theoretically support liquid water
	//on its surface before runaway greenhouse effects make it too hot.
	HabitableZone_AU_Low float64 `json:"habitable_zone_au_low"`

	//The outer edge of the classical habitable zone in AU. Beyond this distance, even with
	//a thick atmosphere and greenhouse warming, liquid water cannot persist because surface
	//temperatures fall below freezing.
	HabitableZone_AU_High float64 `json:"habitable_zone_au_high"`

	//The distance from the star where volatile compounds (primarily water, but also ammonia and
	//methane) condense into solid ice grains. Inside this line, these compounds remain gaseous;
	//outside, they form icy planetesimals. This line is crucial for planet formation models.
	SnowLine_AU float64 `json:"snow_line_au"`

	//The outer boundary of the star’s planetary system or the region of interest for stellar
	//dynamics, often marking where the star’s gravitational influence becomes weak, where planet
	//formation effectively stops, or where the protoplanetary disk fades.
	OuterLimit_AU float64 `json:"outer_limit_au"`

	//Orbit - field exist for every star that is not Primary
	Orbit *orbit.Orbit `json:"orbit,omitempty"`

	//children is a technical fields to describe sattelites
	Children []SystemBody `json:"children,omitempty"`
}

func NewStar() *Star {
	return &Star{}
}

func Import(stellar string) []*Star {
	stars := []*Star{}
	starKeys := Parse(stellar)
	for _, key := range starKeys {
		stel, num, lum := parseKey(key)
		st := &Star{
			StellarClass:       stel,
			NumeriacalSubClass: num,
			LuminocityClass:    lum,
		}
		stars = append(stars, st)
	}
	return stars
}

func (s *Star) MarshalJSON() ([]byte, error) {
	type Alias Star
	return json.MarshalIndent(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "star",
		Alias: (*Alias)(s),
	}, "", "  ")
}

func (s *Star) UnmarshalJSON(data []byte) error {
	type Alias Star
	aux := &struct {
		Children json.RawMessage `json:"children"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Children) > 0 {
		children, err := UnmarshalBodies(aux.Children)
		if err != nil {
			return err
		}
		s.Children = children
	}
	return nil
}

func (s *Star) GetDesignation() string {
	return s.Designation
}

func (s *Star) GetType() string {
	return "star"
}

func (s *Star) GetOrbit() *orbit.Orbit {
	return s.Orbit
}

func (s *Star) GetChildren() []SystemBody {
	return s.Children
}

func (s *Star) SetChildren(children []SystemBody) {
	s.Children = children
}

func Unmarshal(raw json.RawMessage) (*Star, error) {
	return &Star{}, nil
}

type StarKey string

func defineStarKeys() map[StarKey]bool {
	stars := make(map[StarKey]bool)
	for _, stellarClass := range []string{"O", "B", "A", "F", "G", "K", "M", "L", "T", "Y", "D"} {
		for _, numericalSubclass := range []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"} {
			for _, luminocityClass := range []string{"Ia", "Ib", "II", "III", "IV", "V", "VI", "BD"} {
				key := fmt.Sprintf("%s%s %s", stellarClass, numericalSubclass, luminocityClass)
				switch stellarClass {
				case "L", "T", "Y", "D":
					key = fmt.Sprintf("%s%s", stellarClass, numericalSubclass)
					stars[StarKey(key)] = true
					continue
				case "O", "B", "A", "F":
					if luminocityClass == "VI" {
						key = strings.ReplaceAll(key, "VI", "V")
					}
				}
				if luminocityClass == "BD" {
					key = luminocityClass
					stars[StarKey(key)] = true
					continue
				}
				stars[StarKey(key)] = true
			}
		}
	}
	return stars
}

func Parse(s string) []StarKey {
	validKeys := defineStarKeys()

	keys := make([]string, 0, len(validKeys))
	for k := range validKeys {
		keys = append(keys, string(k))
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	result := make([]StarKey, 0)
	n := len(s)
	i := 0
	for i < n {
		found := false
		for _, key := range keys {
			l := len(key)
			if i+l <= n && s[i:i+l] == key {
				result = append(result, StarKey(key))
				i += l
				found = true
				break
			}
		}
		if !found {
			i++
		}
	}
	return result
}

// parseKey разбирает StarKey на составляющие:
// stellarClass (O,B,A,F,G,K,M,L,T,Y,D или пустая строка для BD),
// numericalSubclass (0-9 или пустая строка),
// luminocityClass (Ia,Ib,II,III,IV,V,VI,D,BD или пустая строка).
func parseKey(sk StarKey) (string, string, string) {
	s := string(sk)

	// Специальный случай: "BD" (коричневый карлик)
	if s == "BD" {
		return "", "", "BD"
	}

	if len(s) < 2 {
		return "", "", ""
	}

	// Первый символ должен быть буквой класса
	stellarClass := s[0:1]
	if !isValidStellarClass(stellarClass) {
		return "", "", ""
	}

	// Второй символ - цифра подкласса
	numericalSubclass := s[1:2]
	if numericalSubclass < "0" || numericalSubclass > "9" {
		return "", "", ""
	}

	// Если длина ровно 2 -> формат "Xy" (L,T,Y,D)
	if len(s) == 2 {
		return stellarClass, numericalSubclass, ""
	}

	// Иначе ожидаем пробел и luminosityClass
	if len(s) >= 4 && s[2] == ' ' {
		luminocityClass := s[3:]
		if isValidLuminosityClass(luminocityClass) {
			return stellarClass, numericalSubclass, luminocityClass
		}
		// даже если невалиден, возвращаем как есть (на случай расширения)
		return stellarClass, numericalSubclass, luminocityClass
	}

	// Неизвестный формат
	return "", "", ""
}

// isValidStellarClass проверяет допустимые спектральные классы (включая D)
func isValidStellarClass(c string) bool {
	switch c {
	case "O", "B", "A", "F", "G", "K", "M", "L", "T", "Y", "D":
		return true
	default:
		return false
	}
}

// isValidLuminosityClass проверяет допустимые классы светимости (включая BD)
func isValidLuminosityClass(lc string) bool {
	switch lc {
	case "Ia", "Ib", "II", "III", "IV", "V", "VI", "D", "BD":
		return true
	default:
		return false
	}
}
