package uwp

import (
	"fmt"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/support/valueobject/ehex"
)

const (
	port      = 0
	size      = 1
	atmo      = 2
	hydr      = 3
	pops      = 4
	govr      = 5
	laws      = 6
	separator = 7
	tl        = 8
)

type UWP string

func (u UWP) Profile() map[int]ehex.Ehex {
	sl := strings.Split(string(u), "")
	ehexMap := make(map[int]ehex.Ehex)
	for i, code := range sl {
		eh := ehex.FromCode(code)
		switch i {
		case port:
			eh = ehex.SetDescription(eh, describePort(code))
		case size:
			eh = ehex.SetDescription(eh, describeSize(code))
		case atmo:
			eh = ehex.SetDescription(eh, describeAtmo(code))
		case hydr:
			eh = ehex.SetDescription(eh, describeHydr(code))
		case pops:
			eh = ehex.SetDescription(eh, describePops(code))
		case govr:
			eh = ehex.SetDescription(eh, describeGovr(code))
		case laws:
			eh = ehex.SetDescription(eh, describeLaws(code))
		case tl:
			eh = ehex.SetDescription(eh, describeTL(code))
		}
		ehexMap[i] = eh
	}
	return ehexMap
}

func (u UWP) String() string {
	return string(u)
}

func (u UWP) Description() string {
	s := fmt.Sprintf(" %v \n", u.String())
	s += " ||||||| | \n"
	s += fmt.Sprintf(" ||||||| +-Tech Level : %v\n", u.TL().Description())
	s += fmt.Sprintf(" ||||||+---Laws Level : %v\n", u.Laws().Description())
	s += fmt.Sprintf(" |||||+----Govenment  : %v\n", u.Govr().Description())
	s += fmt.Sprintf(" ||||+-----Population : %v\n", u.Pops().Description())
	s += fmt.Sprintf(" |||+------Hydrosphere: %v\n", u.Hydr().Description())
	s += fmt.Sprintf(" ||+-------Atmosphere : %v\n", u.Atmo().Description())
	s += fmt.Sprintf(" |+--------Size       : %v\n", u.Size().Description())
	s += fmt.Sprintf(" +---------Starpoprt  : %v\n", u.Port().Description())
	return s
}

func New(s string) (UWP, error) {
	issues := []string{}
	codes := strings.Split(s, "")
	if len(codes) < 9 {
		issues = append(issues, fmt.Sprintf("string contains less than 9 characters (%v)", len(codes)))
	}
	if len(codes) > 9 {
		issues = append(issues, fmt.Sprintf("string contains more than 9 characters (%v)", len(codes)))
	}
	for i, code := range codes {
		switch i {
		case separator:
			if eh := ehex.FromCode(code); eh.Code() != "-" {
				issues = append(issues, fmt.Sprintf("code %v is not a separator: %v", i, eh))
			}
		default:
			if eh := ehex.FromCode(code); eh.Value() < 0 {
				issues = append(issues, fmt.Sprintf("code %v is not defined: %v", i, eh.Description()))
			}
		}
	}
	u := UWP(s)
	if len(issues) != 0 {
		errText := "uwp contains issues:"
		for i, iss := range issues {
			errText += fmt.Sprintf("\n%v: %v", i+1, iss)
		}
		return u, fmt.Errorf("%s", errText)
	}
	return u, nil
}

func (u UWP) Size() ehex.Ehex {
	return u.Profile()[size]
}

func (u UWP) Atmo() ehex.Ehex {
	return u.Profile()[atmo]
}
func (u UWP) Hydr() ehex.Ehex {
	return u.Profile()[hydr]
}
func (u UWP) Pops() ehex.Ehex {
	return u.Profile()[pops]
}
func (u UWP) Govr() ehex.Ehex {
	return u.Profile()[govr]
}
func (u UWP) Laws() ehex.Ehex {
	return u.Profile()[laws]
}
func (u UWP) TL() ehex.Ehex {
	return u.Profile()[tl]
}
func (u UWP) Port() ehex.Ehex {
	return u.Profile()[port]
}

func describeSize(code string) string {
	d := "<invalid size code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "Diameter N/A (One or more small bodies, an asteroid or planetoid belt)",
		"R": "Diameter N/A (Planetary rings)",
		"S": "Diameter 400-799km (These small bodies are dwarf planets or significant moons)",
		"1": "Diameter 800–2,399km  (Small planets which may also exist in asteroid or planetoid belts)",
		"2": "Diameter 2,400–3,999km (Luna sized body)",
		"3": "Diameter 4,000–5,599km (Mercury, Ganymede, Titan sized body)",
		"4": "Diameter 5,600–7,199km (Mars sized body)",
		"5": "Diameter 7,200–8,799km (Larger than Mars)",
		"6": "Diameter 8,800–10,399km (Average between Terra and Mars)",
		"7": "Diameter 10,400–11,999km (Smaller than Terra)",
		"8": "Diameter 12,000–13,599km (Venus, Terra sized body)",
		"9": "Diameter 13,600–15,199km (Larger than Earth)",
		"A": "Diameter 15,200–16,799km (Much larger than Earth)",
		"B": "Diameter 16,800–18,399km (Small Super-Earth)",
		"C": "Diameter 18,400–19,999km (Super-Earth)",
		"D": "Diameter 20,000–21,599km (Large Super-Earth)",
		"E": "Diameter 21,600–23,199km (Huge Super-Earth)",
		"F": "Diameter 23,200–24,799km (Maximum Super-Earth)",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}

func describeAtmo(code string) string {
	d := "<invalid atmosphere code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "None (Vacc Suite required)",
		"1": "Trace (Vacc Suite required)",
		"2": "Very Thin, Tainted (Respirator and Filter required)",
		"3": "Very Thin (Respirator required)",
		"4": "Thin, Tainted (Filter required)",
		"5": "Thin",
		"6": "Standard",
		"7": "Standard, Tainted (Filter required)",
		"8": "Dense",
		"9": "Dense, Tainted (Filter required)",
		"A": "Exotic (Air Supply required)",
		"B": "Corrosive (Vacc Suit required)",
		"C": "Insidious (Vacc Suit required)",
		"D": "Very Dense (Varies by altitude)",
		"E": "Low (Varies by altitude)",
		"F": "Unusual (Varies)",
		"G": "Gas, Helium (HEV Suit required)",
		"H": "Gas, Hydrogen (Not Survivable)",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}

func describeHydr(code string) string {
	d := "<invalid hydrosphere code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "0-5% liquid surface cover",
		"1": "6-15% liquid surface cover",
		"2": "16-25% liquid surface cover",
		"3": "26-35% liquid surface cover",
		"4": "36-45% liquid surface cover",
		"5": "46-55% liquid surface cover",
		"6": "56-65% liquid surface cover",
		"7": "66-75% liquid surface cover",
		"8": "76-85% liquid surface cover",
		"9": "86-95% liquid surface cover",
		"A": "96-100% liquid surface cover",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}
func describePops(code string) string {
	d := "<invalid population code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "No Registrered or Permanent population",
		"1": "Registred population 1-99",
		"2": "Registred population 100-999",
		"3": "Registred population 1,000-9,999",
		"4": "Registred population 10,000-99,999",
		"5": "Registred population 100,000-999,999",
		"6": "Registred population 1,000,000-9,999,999",
		"7": "Registred population 10,000,000-99,999,999",
		"8": "Registred population 100,000,000-999,999,999",
		"9": "Registred population 1,000,000,000-9,999,999,999",
		"A": "Registred population 10,000,000,000-99,999,999,999",
		"B": "Registred population 100,000,000,000-999,999,999,999",
		"C": "Registred population above 1,000,000,000,000",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}
func describeGovr(code string) string {
	d := "<invalid government code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "No Government (In many cases, family bonds predominate)",
		"1": "Company/Corporation (Company managerial elite and most citizenry are company employees or dependants)",
		"2": "Participating Democracy (Ruling by the advice and consent of the citizenry directly)",
		"3": "Self-Perpetuating Oligarchy (Restricted minority, with little or no input from the mass of citizenry)",
		"4": "Representative Democracy (Elected representatives)",
		"5": "Feudal Technocracy (Relationships are based on the performance of technical activities that are mutually beneficial)",
		"6": "Captive Government (Imposed leadership answerable to an outside group)",
		"7": "Balkanisation (No central authority exists. Law Level refers to the government nearest the starport)",
		"8": "Civil Service Bureaucracy (Government agencies employing individuals selected for their expertise)",
		"9": "Impersonal Bureaucracy (Ruling by agencies that have become insulated from the governed citizens)",
		"A": "Charismatic Dictatorship (Ruling by a single leader who enjoys overwhelming confidence of the citizens)",
		"B": "Non-Charismatic Dictatorship (A dictator who inherited power through normal channels)",
		"C": "Charismatic Oligarchy (A select group or class that enjoys the overwhelming confidence of the citizenry)",
		"D": "Religious Dictatorship (Ruling by religious organisation without regard to the individual needs of the citizenry)",
		"E": "Religious Autocracy (Government by a single religious leader having absolute power over the citizenry)",
		"F": "Totalitarian Oligarchy (Government by an all-powerful minority which maintains absolute control through widespread coercion and oppression)",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}
func describeLaws(code string) string {
	d := "<invalid law code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "No formal legal system",
		"1": "Grave and serious crimes prosecuted",
		"2": "Moderate crimes prosecuted",
		"3": "Minor crimes prosecuted",
		"4": "Petty crimes prosecuted",
		"5": "Trivial crimes prosecuted",
		"6": "Public surveillance",
		"7": "Insignificant crimes prosecuted",
		"8": "Indefinite detention allowed",
		"9": "No effective right to counsel",
		"A": "Pre-emptive detention allowed",
		"B": "Arbitrary indefinite detention allowed",
		"C": "Arbitrary verdicts without defendant participation",
		"D": "Paramilitary law enforcement, thought crimes prosecuted",
		"E": "Fully-fledged police state, arbitrary executions or 'disappearances'",
		"F": "Rigid control of daily life, gulag state",
		"G": "Thoughts controlled, disproportionate punishments",
		"H": "Legalised oppression",
		"J": "Routine oppression",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}
func describeTL(code string) string {
	d := "<invalid tech level code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"0": "Primitive (Stone Age)",
		"1": "Primitive (Bronze Age/Medieval)",
		"2": "Primitive (Age of Sail/Renesance: 16-17 Century)",
		"3": "Primitive (Industrial Revolution: 18-19 Century)",
		"4": "Industrial (Mechanization: 1900-1940)",
		"5": "Industrial (Circa 1940-1960)",
		"6": "Industrial (Nuclear Age: 1960+)",
		"7": "Pre-Stellar (Computerization: 1980+)",
		"8": "Pre-Stellar (Modern/Near Future Earth)",
		"9": "Pre-Stellar (Late 21 Century)",
		"A": "Early Stellar (Commonly available jump drives)",
		"B": "Early Stellar (True Artificial Intelligence/Jump-2)",
		"C": "Average Stellar (Common Imperial/Plasma/Jump-3)",
		"D": "Average Stellar (Battle Dress/Clonning/Jump-4)",
		"E": "Average Stellar (Aslan/Solomani High/Jump-5)",
		"F": "High Stellar (Imperial High/Jump-6)",
		"G": "High Stellar (Darrian High)",
		"H": "Advanced Stellar ()",
		"J": "Advanced Stellar ()",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}
func describePort(code string) string {
	d := "<invalid starport code>"
	if set := ehex.FromCode(code).Description(); set != "" {
		d = set
	}
	descr := map[string]string{
		"A": "Excelent Quality Starport",
		"B": "Good Quality Starport",
		"C": "Routine Quality Starport",
		"D": "Poor Quality Starport",
		"E": "Frontier Outpost",
		"X": "No Starport",
	}
	if descr[code] != "" {
		d = descr[code]
	}
	return d
}
