package uwp

import (
	"fmt"
	"slices"
	"strings"
)

const (
	Port = "Port"
	Size = "Size"
	Atmo = "Atmosphere"
	Hydr = "Hydrosphere"
	Pops = "Population"
	Govr = "Government"
	Laws = "Laws"
	TL   = "TL"
)

type ProfileValue struct {
	Index       string            `json:"index,omitempty"`
	Category    string            `json:"category,omitempty"`
	Code        string            `json:"code"`
	Numerical   int               `json:"value"`
	Description map[string]string `json:"description,omitempty"`
}

func Description(category, code string) map[string]string {
	descriptionMap := make(map[string]string)
	//PORT
	if category == Port {
		switch code {
		case "A":
			descriptionMap["en"] = "An 'A' code is assigned to excellent-quality installations, usually with construction yards attached. The size of a port’s shipyards can vary considerably. Small-ship construction is commonplace, building vessels of common types up to 1,000 tons or so. Some ports can build larger vessels but many such craft are built at specialist yards. These may or may not be co-located with a major starport but the facility will need a port of some kind.\nA Class A port can conduct starship overhauls and modifications, and sells both refined and unrefined fuel."
		case "B":
			descriptionMap["en"] = "Code B is assigned to large or high-quality (sometimes but not always both) installations. Ports serving systems that see a lot of interstellar traffic usually grow into a Class B installation. To qualify for this code, a starport must be capable of conducting maintenance and overhauls, and must have facilities for refining starship fuel. It will usually have yards capable of constructing non-starships such as shuttles and launches. Class B ports are often important trade centres."
		case "C":
			descriptionMap["en"] = "Code C is assigned to 'average' starports, which can vary considerably. To qualify, a port must have at least a basic orbital facility, although this might be nothing more than a docking station with a shuttle service planetside. Both refined and unrefined fuel are usually available but starship building facilities are highly unlikely to be present, and only fairly basic maintenance or repairs are possible."
		case "D":
			descriptionMap["en"] = "Code D is assigned to a very basic port. A Class D port is almost always a downport with no orbital component at all. Only unrefined fuel is available and no repair facilities are on offer; what facilities there are may be quite primitive."
		case "E":
			descriptionMap["en"] = "Code E indicates what is often termed a 'frontier installation'. A Class E port can be nothing more than a ‘usual place to land’ or a known safe area. There are no facilities as such, although there might be a settlement close by where ad-hoc services can be obtained."
		case "X":
			descriptionMap["en"] = "Code X applies to worlds with no port at all, such as unexplored, interdicted and uninhabited worlds. Not only are there no port facilities, there is no guarantee that a landing site will be safe. Ground subsidence, flooding and other problems can make an apparently useable site hazardous."
		case "F":
			descriptionMap["en"] = "This is assigned to a high-quality port capable of handling a large volume of traffic. Class F spaceports usually have facilities equivalent to a Class B starport or perhaps a high-end C. A few are as large and capable as Class A ports and may have additional facilities such as large shipyards or a naval base. Class F ports often serve large cities on a planet with a central starport, or important offworld installations, such as colonies and the like."
		case "G":
			descriptionMap["en"] = "A G-code is assigned to a basic or fairly poor spaceport, equivalent to the lower end of Class C or better than average Class D starports. A Class G port is likely to be found at a minor city or typical offworld asset such as a mining base on a gas giant moon."
		case "H":
			descriptionMap["en"] = "An H-code is typically assigned to a very basic or improvised facility equivalent to a low-end Class D port, or simply a landing area next to a small or temporary installation. A scientific outpost with a dozen personnel, which receives a supply ship once every three months, would probably have Class H port."
		case "Y":
			descriptionMap["en"] = "Code Y indicates no spaceport is present. The starport equivalent code, Class X, can also be used."
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	//SIZE
	if category == Size {
		switch code {
		case "0":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 0-400km; %v; One or more small bodies, an asteroid or planetoid belt", gravity(0))
		case "R":
			descriptionMap["en"] = "This is a special code for planetary rings"
		case "S":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 400–799km; %v; These small bodies are dwarf planets or significant moons", gravity(0))
		case "1":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 800–2,399km; %v bar; Small planets which may also exist in asteroid or planetoid belts", gravity(1))
		case "2":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 2,400-3,999km; %v; Example: Luna", gravity(2))
		case "3":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 4,000-5,599km; %v; Examples: Mercury, Ganymede, Titan", gravity(3))
		case "4":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 5,600–7,199km; %v; Example: Mars", gravity(4))
		case "5":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 7,200– 8,799km; %v", gravity(5))
		case "6":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 8,800–10,399km; %v", gravity(6))
		case "7":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 10,400–11,999km; %v", gravity(7))
		case "8":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 12,000–13,599km; %v; Examples: Venus, Terra", gravity(8))
		case "9":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 13,600–15,199km; %v bar; Super-Earth", gravity(9))
		case "A":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 15,200–16,799km; %v", gravity(10))
		case "B":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 16,800–18,399km; %v", gravity(11))
		case "C":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 18,400–19,999km; %v", gravity(12))
		case "D":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 20,000–21,599km; %v", gravity(13))
		case "E":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 21,600–23,199km; %v", gravity(14))
		case "F":
			descriptionMap["en"] = fmt.Sprintf("Diameter: 23,200– 25,799km; %v; Maximum Super-Earth", gravity(15))
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	if category == Atmo {
		switch code {
		case "0":
			descriptionMap["en"] = "None; Pressure: 0.00–0.0009 bar; Survival Gear: Vacc Suit; Examples: Mercury, Luna"
		case "1":
			descriptionMap["en"] = "Trace; Pressure: 0.001–0.09 bar; Survival Gear: Vacc Suit; Example: Mars"
		case "2":
			descriptionMap["en"] = "Very Thin, Tainted; Pressure: 0.1–0.42 bar; Survival Gear: Respirator and Filter"
		case "3":
			descriptionMap["en"] = "Very Thin; Pressure: 0.1–0.42 bar; Survival Gear: Respirator"
		case "4":
			descriptionMap["en"] = "Thin, Tainted 0.43–0.70 0.27 bar; Survival Gear: Filter"
		case "5":
			descriptionMap["en"] = "Thin; Pressure: 0.43–0.70 bar; Survival Gear: None"
		case "6":
			descriptionMap["en"] = "Standard; Pressure: 0.70–1.49 bar; Survival Gear: None; Example: Terra"
		case "7":
			descriptionMap["en"] = "Standard, Tainted; Pressure: 0.70–1.49 bar; Survival Gear: Filter"
		case "8":
			descriptionMap["en"] = "Dense; Pressure: 1.50–2.49 bar; Survival Gear: None"
		case "9":
			descriptionMap["en"] = "Dense, Tainted; Pressure: 1.50–2.49 bar; Survival Gear: Filter"
		case "A":
			descriptionMap["en"] = "Exotic; Pressure: Varies bar; Survival Gear: Air Supply; Example: Titan"
		case "B":
			descriptionMap["en"] = "Corrosive; Pressure: Varies bar; Survival Gear: Vacc Suit; Example: Venus"
		case "C":
			descriptionMap["en"] = "Insidious; Pressure: Varies bar; Survival Gear: Vacc Suit"
		case "D":
			descriptionMap["en"] = "Very Dense; Pressure: 2.50–10.0 bar; Survival Gear: Varies by altitude"
		case "E":
			descriptionMap["en"] = "Low; Pressure: 0.10–0.42 bar; Survival Gear: Varies by altitude"
		case "F":
			descriptionMap["en"] = "Unusual; Pressure: Varies bar; Survival Gear: Varies"
		case "G":
			descriptionMap["en"] = "Gas, Helium; Pressure: 100+ HEV bar; Survival Gear: Suit; Notes: Dense helium-dominated gas"
		case "H":
			descriptionMap["en"] = "Gas, Hydrogen; Pressure: 1,000+ bar; Survival Gear: Not; Notes: Survivable Gas Dwarf"
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	if category == Hydr {
		switch code {
		case "0":
			descriptionMap["en"] = "Surface Cover: 0%–5%"
		case "1":
			descriptionMap["en"] = "Surface Cover: 6%–15%"
		case "2":
			descriptionMap["en"] = "Surface Cover: 16%–25%"
		case "3":
			descriptionMap["en"] = "Surface Cover: 26%–35%"
		case "4":
			descriptionMap["en"] = "Surface Cover: 36%–45%"
		case "5":
			descriptionMap["en"] = "Surface Cover: 46%–55%"
		case "6":
			descriptionMap["en"] = "Surface Cover: 56%–65%"
		case "7":
			descriptionMap["en"] = "Surface Cover: 66%–75%"
		case "8":
			descriptionMap["en"] = "Surface Cover: 76%–85%"
		case "9":
			descriptionMap["en"] = "Surface Cover: 86%–95%"
		case "A":
			descriptionMap["en"] = "Surface Cover: 96%–100%"
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	if category == Pops {
		switch code {
		case "0":
			descriptionMap["en"] = "Few or no registered inhabitants"
		case "1":
			descriptionMap["en"] = "Population Range: 1 – 99"
		case "2":
			descriptionMap["en"] = "Population Range: 100 – 999"
		case "3":
			descriptionMap["en"] = "Population Range: 1,000 – 9,999"
		case "4":
			descriptionMap["en"] = "Population Range: 10,000 – 99,999"
		case "5":
			descriptionMap["en"] = "Population Range: 100,000 – 999,999"
		case "6":
			descriptionMap["en"] = "Population Range: 1,000,000 – 9,999,999"
		case "7":
			descriptionMap["en"] = "Population Range: 10,000,000 – 99,999,999"
		case "8":
			descriptionMap["en"] = "Population Range: 100,000,000 – 999,999,999"
		case "9":
			descriptionMap["en"] = "Population Range: 1,000,000,000 – 9,999,999,999"
		case "A":
			descriptionMap["en"] = "Population Range: 10,000,000,000 – 99,999,999,999"
		case "B":
			descriptionMap["en"] = "Population Range: 100,000,000,000 – 999,999,999,999"
		case "C":
			descriptionMap["en"] = "Population Range: 1,000,000,000,000 – 9,999,999,999,999"
		case "D":
			descriptionMap["en"] = "Population Range: 10,000,000,000,000 – 99,999,999,999,999"
		case "E":
			descriptionMap["en"] = "Population Range: 100,000,000,000,000 – 999,999,999,999,999"
		case "F":
			descriptionMap["en"] = "Population Range: 1,000,000,000,000,000 – 9,999,999,999,999,999"
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	if category == Govr {
		switch code {
		case "0":
			descriptionMap["en"] = "No Government Structure"
		case "1":
			descriptionMap["en"] = "Company / Corporation"
		case "2":
			descriptionMap["en"] = "Participating Democracy"
		case "3":
			descriptionMap["en"] = "Self-Perpetuating Oligarhy"
		case "4":
			descriptionMap["en"] = "Represintative Democracy"
		case "5":
			descriptionMap["en"] = "Feudal Technocracy"
		case "6":
			descriptionMap["en"] = "Captive Government / Colony"
		case "7":
			descriptionMap["en"] = "Balkanization"
		case "8":
			descriptionMap["en"] = "Civil Service Bureaucracy"
		case "9":
			descriptionMap["en"] = "Imperial Bureaucracy"
		case "A":
			descriptionMap["en"] = "Charismatic Dictator"
		case "B":
			descriptionMap["en"] = "Non-Charismatic Dictator"
		case "C":
			descriptionMap["en"] = "Charismatic Oligarhy"
		case "D":
			descriptionMap["en"] = "Religious Dictatirship"
		case "E":
			descriptionMap["en"] = "Religious Autocracy"
		case "F":
			descriptionMap["en"] = "Totalitarian Oligarhy"
		case "G":
			descriptionMap["en"] = "Small Station or Facility (Aslan)"
		case "H":
			descriptionMap["en"] = "Split Clan Control (Aslan)"
		case "J":
			descriptionMap["en"] = "Single On-world Clan Control (Aslan)"
		case "K":
			descriptionMap["en"] = "Single Multi-world Clan Control (Aslan)"
		case "L":
			descriptionMap["en"] = "Major Clan Control (Aslan)"
		case "M":
			descriptionMap["en"] = "Vassal Clan Control (Aslan)"
		case "N":
			descriptionMap["en"] = "Major Vassal Clan Control (Aslan)"
		case "P":
			descriptionMap["en"] = "Small Station or Facility (K'kree)"
		case "Q":
			descriptionMap["en"] = "Krurruna or Krumanak Rule for Off-world Steppelord (K'kree)"
		case "R":
			descriptionMap["en"] = "Steppelord On-world Rule (K'kree)"
		case "S":
			descriptionMap["en"] = "Sept (Hiver)"
		case "T":
			descriptionMap["en"] = "Unsupervised Anarchy (Hiver)"
		case "U":
			descriptionMap["en"] = "Supervised Anarchy (Hiver)"
		case "W":
			descriptionMap["en"] = "Committie (Hiver)"
		case "X":
			descriptionMap["en"] = "Drone Hierarchy (Droyne)"
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	if category == Laws {
		switch code {
		case "0":
			descriptionMap["en"] = "No restrictions; No contract law or licenses required"
		case "1":
			descriptionMap["en"] = "Phobited: Poison gas, explosives, undetectable weapons, weapons of mass destruction, battle dress (C5); Optional registration of private agreements, claim registration"
		case "2":
			descriptionMap["en"] = "Phobited: Portable energy and laser weapons, combat armour; Registration of corporations, enforcement of claims "
		case "3":
			descriptionMap["en"] = "Phobited: Military weapons (all portable heavy weapons), flak jackets and obvious armour (C4); Basic permitting and zoning laws, required licensing of corporations and tax reporting, bankruptcy law "
		case "4":
			descriptionMap["en"] = "Phobited: Light assault weapons and submachine guns (all fully automatic weapons), cloth armour (C3); Registration of professional licenses, periodic random auditing of major financial transactions"
		case "5":
			descriptionMap["en"] = "Phobited: Personal concealable ranged weapons (auto pistols and revolvers), mesh armour; Required professional licenses for most skilled professions"
		case "6":
			descriptionMap["en"] = "Phobited: All firearms except shotguns and stunners; carrying weapons discouraged; Moderate permitting and zoning laws, registration fees required for professional licenses"
		case "7":
			descriptionMap["en"] = "Phobited: Shotguns and all other ranged firearms (C2); Professional licenses required for all skilled labour, periodic auditing of major financial transactions "
		case "8":
			descriptionMap["en"] = "Phobited: All bladed weapons, stunners, all visible armour (C1); Restrictive zoning and permitting laws"
		case "9":
			descriptionMap["en"] = "Phobited: All weapons, including knives longer than 10cm, all armour; Active auditing of all financial transactions"
		case "A":
			descriptionMap["en"] = "All Weapons phobited. Violations are treated as Serious crimes; Arduous permitting and zoning laws"
		case "B":
			descriptionMap["en"] = "Random sweeps for weapons violations; Continuous auditing of all financial transactions"
		case "C":
			descriptionMap["en"] = "Active monitoring for ownership violations; All economic regulation enforcement transferred to criminal justice system; Rigid control of civilian movement"
		case "D":
			descriptionMap["en"] = "Paramilitary law enforcement"
		case "E":
			descriptionMap["en"] = "Full-fledged police state"
		case "F":
			descriptionMap["en"] = "All facets of daily life regularly ligislated and regulated"
		case "G":
			descriptionMap["en"] = "Severe punishments for petty infractions"
		case "H":
			descriptionMap["en"] = "Legalized oppressive practices"
		case "J":
			descriptionMap["en"] = "Routinley oppressive and restrictive"
		case "K":
			descriptionMap["en"] = "Excessively oppressive and restrictive"
		case "L":
			descriptionMap["en"] = "Totaly oppressive and restrictive"
		case "S":
			descriptionMap["en"] = "Special / Variable situation"
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	if category == TL {
		switch code {
		case "0":
			descriptionMap["en"] = "Stone Age; Primitive"
		case "1":
			descriptionMap["en"] = "Bronze, Iron; Bronze Age to Middle Ages"
		case "2":
			descriptionMap["en"] = "Printing Press; circa 1400 to 1700"
		case "3":
			descriptionMap["en"] = "Basic Science; circa 1700 to 1860"
		case "4":
			descriptionMap["en"] = "External Combustion; circa 1860 to 1910"
		case "5":
			descriptionMap["en"] = "Mass Production; circa 1910 to 1940"
		case "6":
			descriptionMap["en"] = "Nuclear Power; circa 1940 to 1970"
		case "7":
			descriptionMap["en"] = "Minitaturized Electronics; circa 1970 to 2000"
		case "8":
			descriptionMap["en"] = "Quality Computer; Modern Age Terra"
		case "9":
			descriptionMap["en"] = "Anti-Gravity"
		case "A":
			descriptionMap["en"] = "Interstellar Community; Reliable Jump-1 Drive"
		case "B":
			descriptionMap["en"] = "Lower Average Imperial; Reliable Jump-2 Drive"
		case "C":
			descriptionMap["en"] = "Average Imperial; Reliable Jump-3 Drive"
		case "D":
			descriptionMap["en"] = "Above Imperial Average; Reliable Jump-4 Drive"
		case "E":
			descriptionMap["en"] = "Above Imperial Average; Reliable Jump-5 Drive"
		case "F":
			descriptionMap["en"] = "Technical Imperial Average; Reliable Jump-6 Drive"
		case "G":
			descriptionMap["en"] = "Self-Aware Robots"
		case "H":
			descriptionMap["en"] = "Artificial Intelligence"
		case "J":
			descriptionMap["en"] = "Personal Desintegrators"
		case "K":
			descriptionMap["en"] = "Plastic Metals"
		case "L":
			descriptionMap["en"] = "Comprehensible only as technological magic"
		case "?":
			descriptionMap["en"] = `<<UNKNOWN>>`
		}
	}
	descriptionMap["ru"] = descriptionMap["en"]
	return descriptionMap
	/*

	 */

}

func gravity(i int) string {
	minG := (lowsizes[i] * 0.82) / 12742.0
	maxG := (lowsizes[i+1] * 1.12) / 12742.0
	return fmt.Sprintf("Gravity: %0.2f-%0.2fg", minG, maxG)
	// return fmt.Sprintf("Average Gravity: %0.2fg", (minG+maxG)/2)
}

var lowsizes = map[int]float64{
	0:  0,
	1:  800,
	2:  2400,
	3:  4000,
	4:  5600,
	5:  7200,
	6:  8800,
	7:  10400,
	8:  12000,
	9:  13600,
	10: 15200,
	11: 16800,
	12: 18400,
	13: 20000,
	14: 21600,
	15: 23200,
	16: 150000,
}

func StringValid(s string) bool {
	data := strings.Split(s, "")
	if len(data) != 9 {
		return false
	}
	for i, value := range data {
		switch i {
		case 7:
			if value != "-" {
				return false
			}
		default:
			if !slices.Contains(uwpValues(), value) {
				return false
			}
		}
	}
	return true
}

func uwpValues() []string {
	return []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "?"}
}
