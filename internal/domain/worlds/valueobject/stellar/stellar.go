package stellar

import (
	"fmt"
	"strings"
)

const (
	Primary     StarDesignation = "Aa"
	Close       StarDesignation = "Ba"
	Near        StarDesignation = "Ca"
	Far         StarDesignation = "Da"
	PrimaryComp StarDesignation = "Ab"
	CloseComp   StarDesignation = "Bb"
	NearComp    StarDesignation = "Cb"
	FarComp     StarDesignation = "Db"
)

// Stellar is an alias for sting describing system stars
type Stellar string

func New(s string) (Stellar, error) {
	if err := Stellar(s).Validate(); err != nil {
		return "", err
	}
	return Stellar(s), nil
}

func (s Stellar) Validate() error {
	for _, str := range ExtractStars(string(s)) {
		if !validateStar(str) {
			return fmt.Errorf("invalid code '%v' in '%v'", str, s)
		}
	}
	return nil
}

func (s Stellar) String() string {
	return strings.Join(ExtractStars(string(s)), ":")
}

func (s Stellar) Split() []string {
	return strings.Split(s.String(), ":")
}

func ExtractStars(stellar string) []string {
	stars := []string{}
	parts := strings.Split(stellar, " ")
	last := len(parts) - 1
	for i := 0; i <= last; i++ {
		switch i == last {
		case false:
			if isClass(parts[i]) {
				stars = append(stars, parts[i])
				continue
			}
			stars = append(stars, parts[i]+" "+parts[i+1])
			i++
		case true:
			stars = append(stars, parts[i])
		}
	}
	return stars
}

func isClass(part string) bool {
	switch part {
	case "Ia", "Ib", "II", "III", "IV", "V", "VI":
		return true
	case "BD", "D", "BH", "NS", "PSR":
		return true
	}
	return false
}

func validateStar(s string) bool {
	switch s {
	case "", "BD", "D", "BH", "NS":
		return true
	}
	validType := false
	validClass := false
	for _, tp := range []string{"O", "B", "A", "F", "G", "K", "M", "L", "T", "Y"} {
		if strings.HasPrefix(s, tp) {
			validType = true
		}
	}
	for _, cl := range []string{"Ia", "Ib", "II", "III", "IV", "V", "VI"} {
		if strings.HasSuffix(s, cl) {
			validClass = true
		}
	}
	return validType && validClass
}

func RollStellarDesignations(r Roller, s Stellar) []StarDesignation {
	stars := ExtractStars(string(s))
	listing := []StarDesignation{}
	for len(listing) != len(stars) {
		listing = RollDesignations(r)
	}
	return listing
}

type StarDesignation string

type Roller interface {
	Roll(string) int
}

func RollDesignations(r Roller) []StarDesignation {
	stars := make(map[StarDesignation]bool)
	stars[Primary] = true
	for _, des := range []StarDesignation{Primary, Close, Near, Far} {
		if des == Primary {
			stars[des] = true
			continue
		}
		if r.Roll("2d6") >= 10 {
			stars[des] = true
		}
	}
	result := []StarDesignation{}
	for _, des := range []StarDesignation{Primary, Close, Near, Far} {
		if stars[des] {
			result = append(result, des)
			code := strings.TrimSuffix(string(des), "a")
			if r.Roll("2d6") >= 10 {
				result = append(result, StarDesignation(code+"b"))
			}
		}
	}
	return result
}
