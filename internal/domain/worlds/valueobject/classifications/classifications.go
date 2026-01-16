package classifications

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/support/valueobject/ehex"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

type Classification string

const (
	Ab Classification = "Ab"
	Ag Classification = "Ag"
	An Classification = "An"
	As Classification = "As"
	Ba Classification = "Ba"
	Bo Classification = "Bo"
	Co Classification = "Co"
	Cp Classification = "Cp"
	Cs Classification = "Cs"
	Cx Classification = "Cx"
	Cy Classification = "Cy"
	Da Classification = "Da"
	De Classification = "De"
	Di Classification = "Di"
	Fa Classification = "Fa"
	Fl Classification = "Fl"
	Fo Classification = "Fo"
	Fr Classification = "Fr"
	Ga Classification = "Ga"
	He Classification = "He"
	Hi Classification = "Hi"
	Ho Classification = "Ho"
	Ht Classification = "Ht"
	Ic Classification = "Ic"
	In Classification = "In"
	Lk Classification = "Lk"
	Lo Classification = "Lo"
	Lt Classification = "Lt"
	Mi Classification = "Mi"
	Mr Classification = "Mr"
	Na Classification = "Na"
	Ni Classification = "Ni"
	Oc Classification = "Oc"
	Pa Classification = "Pa"
	Pe Classification = "Pe"
	Ph Classification = "Ph"
	Pi Classification = "Pi"
	Po Classification = "Po"
	Pr Classification = "Pr"
	Px Classification = "Px"
	Pz Classification = "Pz"
	Re Classification = "Re"
	Ri Classification = "Ri"
	Sa Classification = "Sa"
	Tr Classification = "Tr"
	Tu Classification = "Tu"
	Tz Classification = "Tz"
	Va Classification = "Va"
	Wa Classification = "Wa"

	amberZoneGL = 20
	redZoneGL   = 22
)

func (tc Classification) Description() string {
	switch tc {
	case Ab:
		return "Data Repository\nThe world has a centralized collection point for information and data. Organizations and governments deposit records of their transactions and output in this collection point.\nThe TC refers to AAB, the Imperial designation for data repositories."
	case Ag:
		return "Agricultural\nThe world has climate and conditions which promote farming and ranching. It is a producer of inexpensive foodstuffs. It also is a source of unusual, exotic, or strange delicacies."
	case An:
		return "Ancient Site\nThe world (or the star system) includes one or more locations identified as the ruins of the long-dead race called the Ancients. Ancient Sites are exploited for the Artifact remains of this long dead technological civilization."
	case As:
		return "Asteroid Belt\nThe world is an asteroid belt which is the primary world or mainworld in the system. It is a producer of raw materials and semi-finished goods, especially ores, metals, and minerals."
	case Ba:
		return "Barren World\nThe world has no population, government, or law level. It has never been developed; it has no local infrastructure beyond the starport (if that).\nA Barren world UWP has a default zero Tech Level."
	case Bo:
		return "Boiling World\nNo ice caps, little liquid water."
	case Co:
		return "Cold World\nThe world is at the lower temperature range of human endurance; typically in HZ+1."
	case Cp:
		return "Subsector Capital\nThe world is the political center of a group of tens or dozens of star systems (typically a subsector)."
	case Cs:
		return "Sector Capital\nThe world is the political center of a group of hundreds of star systems (typically a sector)."
	case Cx:
		return "Imperial Capital\nThe world is the overall political center of an interstellar government controlling thousands of star systems."
	case Cy:
		return "Colony\nThe world is a colony Owned by the Most Important, Highest Population, Highest TL world within 6 hexes.\nAdd the remark O:nnnn (=hex of owning world)."
	case Da:
		return "Dangerous\nSome aspect of the world (conditions, customs, laws, life forms, climate, economics, or other) is not well understood or easily understood by typical visitors, and it presents a danger.\nThe world is a TAS Amber Zone."
	case De:
		return "Desert World\nThe world has no open or standing water. This lack of water significantly reduces the level of agricultural development."
	case Di:
		return "Die-Back\nThe world was once extensively settled and developed, but at some time in the last thousand years its inhabiting sophonts died out leaving behind the remnants of their civilization.\nA Die-Back world UWP has a non-zero Tech Level."
	case Fa:
		return "Farming\nThe world has climate and conditions which promote farming and ranching. In addition, it is in the Habitable Zone and not a Mainworld."
	case Fl:
		return "Fluid Oceans\nThe world's oceans are not composed of water. Non-water oceans may be valuable sources of raw materials for industry."
	case Fo:
		return "Forbidden\nSome conditions, customs, laws, life forms, climate, economics, or other circumstance presents an active threat to the health and well-being of individuals. The world is a TAS Red Zone."
	case Fr:
		return "Frozen\nThe world lies substantially beyond the Habitable Zone of the system (HZ+2 or greater) and environmental temperatures are well below the freezing point of many gases."
	case Ga:
		return "Garden World\nThe world is hospitable to most sophonts. Its size, atmosphere, and hydrographic make it an extremely attractive world. A Garden World has a safe environment which does not require protective equipment for humans and sophonts which share the human environment."
	case He:
		return "Hellworld\nThe world is inhospitable to most sophonts. Its size, atmosphere, and hydrographic make it an extremely unattractive world."
	case Hi:
		return "High Population\nThe world's population is one billion or more (Pop = 9 or A or more). High population worlds, because of the economy of scale for production, produce quality inexpensive trade goods."
	case Ho:
		return "Hot World\nThe world is at the upper temperature range of human endurance. No ice caps, little liquid water; typically in HZ-1."
	case Ic:
		return "Ice-Capped\nThe world's water is locked in ice-caps."
	case In:
		return "Industrial\nThe world has a strong manufacturing infrastructure and is a producer of many types of goods."
	case Lk:
		return "Locked\nThe world is a satellite (in orbits Ay through Em) which is locked to the planet it orbits. A Locked satellite does not have a Twilight Zone; its day length equals the time it takes to orbit its planet."
	case Lo:
		return "Low Population\nThe world has a non-zero population less than 10,000. Low Population fluctuates wildly and may change significantly on a yearly (or less) basis.\nLocals are Transients: merchants, corporate employees, military, security, or research personnel."
	case Mi:
		return "Mining\nThe world is the site of extensive mineral resource exploitation. It is not a Mainworld and is located in a star system with an Industrial Mainworld."
	case Mr:
		return "Military Rule\nThe non-Mainworld is ruled by the military from a nearby world."
	case Na:
		return "Non-Agricultural\nThe world is unable to produce enough food agriculturally to feed its population; synthetic food production generally meets basic food needs."
	case Ni:
		return "Non-Industrial\nThe world has a non-zero population (more than 10,000 and less than one million). The TC Non-Industrial remains constant and reflects an expected population level.\nInhabitants of a Non-Industrial world are Settlers: part of a permanent settlement not yet a Colony."
	case Oc:
		return "Ocean World\nThe world surface is covered with very deep seas. There is no (less than a hundredth) land above sea level."
	case Pa:
		return "Pre-Agricultural\nThe world is a candidate for the Agricultural trade classification; its population is just outside the requirement for Agricultural."
	case Pe:
		return "Penal Colony\nThe world is a dumping ground for individuals who will not / do not / cannot conform to standards of behavior."
	case Ph:
		return "Pre-High\nThe world is a candidate for elevation to the High Population trade classification; its population level is just below the requirements for High."
	case Pi:
		return "Pre-Industrial\nThe world is a candidate for the Industrial trade classification; its population is just below the requirements."
	case Po:
		return "Poor\nThe world has poor grade living conditions: a scarcity of water and a relatively sparse atmosphere."
	case Pr:
		return "Pre-Rich\nThe world is a candidate for the Rich trade classification; its population is just outside the criteria for Rich."
	case Px:
		return "Prison, Exile Camp\nThe non-mainworld population consists of criminals or undesirables transported here from other worlds."
	case Pz:
		return "Puzzle\nSome aspect of the world (conditions, customs, laws, life forms, climate, economics, or other) is not well or easily understood by typical visitors. The world is a TAS Amber Zone."
	case Re:
		return "Reserve\nThe world has been set aside (by the highest levels of government) to preserve indigenous life forms, to delay resource development, or to frustrate inquiry into local conditions."
	case Ri:
		return "Rich\nThe world has an untainted atmosphere which is comfortable and attractive for most sophonts, and has a population suitable as a workforce."
	case Sa:
		return "Satellite\nThe world is the satellite of a planet (or gas giant) in the system."
	case Tr:
		return "Tropic\nThe world is relatively warmer than normal (although it is considered habitable). Its orbit is at the inner (warmer) edge of the Habitable Zone.\nThe world has a Hot climate (at the upper limits of human temperature endurance)."
	case Tu:
		return "Tundra\nThe world is relatively colder than normal (although it is considered habitable). Its orbit is at the outer (colder) edge of the Habitable Zone. The world has a Cold climate (at the lower limits of human temperature endurance)."
	case Tz:
		return "Twilight Zone\nThe world is tidally locked with a Temperate band at the Twilight Zone, plus a Hot region (hemisphere) facing the Primary and a Cold region (hemisphere) away from the Primary."
	case Va:
		return "Vacuum World\nThe world has no atmosphere."
	case Wa:
		return "Water World\nThe world surface is covered with water; there is very little land (= less than 10%) above the water surface."
	default:
		return fmt.Sprintf("<undefined trade code '%v'>", string(tc))
	}
}

type classificationOptions struct {
	isSecondaryWorld   bool
	mainworldUWP       uwp.UWP
	hzco               float64
	sateliteOrbit      string
	planetaryOrbit     int
	starVaporationZone int
	requestedCodes     []string
}

func defaultClassificationOptions() classificationOptions {
	return classificationOptions{
		isSecondaryWorld:   false,
		mainworldUWP:       "",
		hzco:               0,
		sateliteOrbit:      "",
		requestedCodes:     []string{},
		planetaryOrbit:     -1,
		starVaporationZone: -10, //TODO: think about inplementing MAO mechanics
	}
}

type classificationOptionFunc func(*classificationOptions)

func WithMainworldUWP(u uwp.UWP) classificationOptionFunc {
	return func(co *classificationOptions) {
		co.isSecondaryWorld = true
		co.mainworldUWP = u
	}
}

func WithSateliteOrbit(orbit string) classificationOptionFunc {
	return func(co *classificationOptions) {
		co.sateliteOrbit = orbit
	}
}

func WithHZCO(hzco float64) classificationOptionFunc {
	return func(co *classificationOptions) {
		co.hzco = hzco
	}
}

func WithCodesRequested(codes ...string) classificationOptionFunc {
	return func(co *classificationOptions) {
		for _, code := range codes {
			co.requestedCodes = append(co.requestedCodes, code)
		}
	}
}

func WithPlanetOrbit(starVaporationZone, planetOrbit int) classificationOptionFunc {
	return func(co *classificationOptions) {
		co.starVaporationZone = starVaporationZone
		co.planetaryOrbit = planetOrbit
	}
}

func Classify(u uwp.UWP, opts ...classificationOptionFunc) []Classification {
	options := defaultClassificationOptions()
	for _, set := range opts {
		set(&options)
	}
	confirmed := []Classification{}

	//Primary Codes
	if match(u.Size(), "0") && match(u.Atmo(), "0") && match(u.Hydr(), "0") {
		confirmed = append(confirmed, As)
	}
	if match(u.Atmo(), "23456789") && match(u.Hydr(), "0") {
		confirmed = append(confirmed, De)
	}
	if match(u.Atmo(), "ABC") && match(u.Hydr(), "123456789A") {
		confirmed = append(confirmed, Fl)
	}
	if match(u.Size(), "678") && match(u.Atmo(), "568") && match(u.Hydr(), "567") {
		confirmed = append(confirmed, Ga)
	}
	if match(u.Size(), "3456789ABC") && match(u.Atmo(), "2479ABC") && match(u.Hydr(), "012") {
		confirmed = append(confirmed, He)
	}
	if match(u.Atmo(), "01") && match(u.Hydr(), "123456789A") {
		confirmed = append(confirmed, Ic)
	}
	if match(u.Size(), "ABCDEF") && match(u.Atmo(), "3456789DEF") && match(u.Hydr(), "123456789A") {
		confirmed = append(confirmed, Oc)
	}
	if match(u.Atmo(), "0") {
		confirmed = append(confirmed, Va)
	}
	if match(u.Size(), "3456789") && match(u.Atmo(), "3456789DEF") && match(u.Hydr(), "A") {
		confirmed = append(confirmed, Wa)
	}

	//Population Codes
	if match(u.Pops(), "0") && match(u.Govr(), "0") && match(u.Laws(), "0") && match(u.TL(), "1+") {
		confirmed = append(confirmed, Di)
	}
	if match(u.Pops(), "0") && match(u.Govr(), "0") && match(u.Laws(), "0") && match(u.TL(), "0") && match(u.Port(), "XEYH") {
		confirmed = append(confirmed, Ba)
	}
	if match(u.Pops(), "123") {
		confirmed = append(confirmed, Lo)
	}
	if match(u.Pops(), "1+") && match(u.TL(), "5-") {
		confirmed = append(confirmed, Lt)
	}
	if match(u.Pops(), "456") {
		confirmed = append(confirmed, Ni)
	}
	if match(u.Pops(), "8") {
		confirmed = append(confirmed, Ph)
	}
	if match(u.Pops(), "9ABCDEF") {
		confirmed = append(confirmed, Hi)
	}
	if match(u.Pops(), "1+") && match(u.TL(), "C+") {
		confirmed = append(confirmed, Ht)
	}

	//Ecomomic Codes
	if match(u.Atmo(), "456789") && match(u.Hydr(), "45678") && match(u.Pops(), "48") {
		confirmed = append(confirmed, Pa)
	}
	if match(u.Atmo(), "456789") && match(u.Hydr(), "45678") && match(u.Pops(), "567") {
		confirmed = append(confirmed, Ag)
	}
	if match(u.Atmo(), "0123") && match(u.Hydr(), "0123") && match(u.Pops(), "6789ABCDEF") {
		confirmed = append(confirmed, Na)
	}
	if match(u.Atmo(), "23AB") && match(u.Hydr(), "12345") && match(u.Pops(), "3456") && match(u.Laws(), "6789") {
		confirmed = append(confirmed, Px)
	}
	if match(u.Atmo(), "012479") && match(u.Pops(), "78") {
		confirmed = append(confirmed, Pi)
	}
	if match(u.Atmo(), "012479ABC") && match(u.Pops(), "9ABCEDF") {
		confirmed = append(confirmed, In)
	}
	if match(u.Atmo(), "2345") && match(u.Hydr(), "0123") {
		confirmed = append(confirmed, Po)
	}
	if match(u.Atmo(), "68") && match(u.Pops(), "59") {
		confirmed = append(confirmed, Pr)
	}
	if match(u.Atmo(), "68") && match(u.Pops(), "678") {
		confirmed = append(confirmed, Ri)
	}

	//Climate Codes TODO: evaluate hzco numbers
	if match(u.Size(), "2+") && match(u.Hydr(), "123456789A") && options.hzco >= 1.1 {
		confirmed = append(confirmed, Fr)
	}
	if match(u.Size(), "2+") && match(u.Hydr(), "123456789A") && options.hzco <= -1.1 {
		confirmed = append(confirmed, Bo)
	}
	if match(u.Size(), "2345A+") && match(u.Atmo(), "23ABCDEF") && options.hzco <= -0.5 && options.hzco > -1.1 {
		confirmed = append(confirmed, Ho)
	}
	if match(u.Size(), "2345A+") && match(u.Atmo(), "23ABCDEF") && options.hzco >= 0.5 && options.hzco < 1.1 {
		confirmed = append(confirmed, Co)
	}
	if slices.Contains([]string{"Ay", "Bee", "Cee", "Dee", "Ee", "Eff", "Gee", "Aitch", "Eye", "Jay", "Kay", "Ell", "Em"}, options.sateliteOrbit) {
		confirmed = append(confirmed, Lk)
	}
	if match(u.Size(), "6789") && match(u.Atmo(), "456789") && match(u.Hydr(), "34567") && options.hzco < -0.5 && options.hzco > -1.1 {
		confirmed = append(confirmed, Tr)
	}
	if match(u.Size(), "6789") && match(u.Atmo(), "456789") && match(u.Hydr(), "34567") && options.hzco > 0.5 && options.hzco < 1.1 {
		confirmed = append(confirmed, Tu)
	}
	if options.sateliteOrbit == "" && (options.planetaryOrbit == options.starVaporationZone || options.planetaryOrbit == options.starVaporationZone+1) {
		confirmed = append(confirmed, Tz)
	}

	//Secondary Codes
	if options.isSecondaryWorld {
		if match(u.Atmo(), "456789") && match(u.Hydr(), "45678") && match(u.Pops(), "23456") && options.hzco > -1.1 && options.hzco < 1.1 {
			confirmed = append(confirmed, Fa)
		}
		if match(u.Pops(), "23456") && options.hzco > -1.1 && options.hzco < 1.1 {
			if slices.Contains(Classify(options.mainworldUWP), In) {
				confirmed = append(confirmed, Mi)
			}
		}
		//TODO: Mr evaluation POP: 23 GOVR: Totalitarian? Base?
		if match(u.Atmo(), "23AB") && match(u.Hydr(), "12345") && match(u.Pops(), "3456") && match(u.Govr(), "6") && match(u.Laws(), "6789") {
			confirmed = append(confirmed, Pe)
		}
		if match(u.Pops(), "1234") && match(u.Govr(), "6") && match(u.Laws(), "45") {
			confirmed = append(confirmed, Re)
		}
	}

	//Political
	if match(u.Port(), "A") && slices.Contains(options.requestedCodes, "Cp") {
		confirmed = append(confirmed, Cp)
	}
	if match(u.Port(), "A") && slices.Contains(options.requestedCodes, "Cs") {
		confirmed = append(confirmed, Cs)
	}
	if match(u.Port(), "A") && slices.Contains(options.requestedCodes, "Cx") {
		confirmed = append(confirmed, Cx)
	}
	if match(u.Pops(), "56789A") && match(u.Govr(), "6") && match(u.Laws(), "0123") {
		confirmed = append(confirmed, Cy)
	}

	//Special Codes
	if slices.Contains([]string{"En", "Oh", "Pee", "Que", "Arr", "Ess", "Tee", "Yu", "Vee", "Dub", "Ex", "Wye", "Zee"}, options.sateliteOrbit) {
		confirmed = append(confirmed, Sa)
	}
	if u.Govr().Value()+u.Laws().Value() >= redZoneGL || slices.Contains(options.requestedCodes, "R") {
		confirmed = append(confirmed, Fo)
	}
	if match(u.Pops(), "789ABCDEF") && (u.Govr().Value()+u.Laws().Value() >= amberZoneGL) || slices.Contains(options.requestedCodes, "A") {
		confirmed = append(confirmed, Pz)
	}
	if match(u.Pops(), "0123456") && (u.Govr().Value()+u.Laws().Value() >= amberZoneGL) || slices.Contains(options.requestedCodes, "A") {
		confirmed = append(confirmed, Da)
	}

	//Direct Requested Codes:
	for _, req := range options.requestedCodes {
		switch req {
		case "Ab":
			confirmed = append(confirmed, Ab)
		case "An":
			confirmed = append(confirmed, An)
		case "Fo":
			confirmed = append(confirmed, Fo)
		case "Mr":
			confirmed = append(confirmed, Mr)
		}
	}

	return confirmed
}

func match(code ehex.Ehex, demand string) bool {
	codeValue := code.Value()
	if codeValue < 0 {
		return false
	}
	demandedCodes := []int{}
	orLess := false
	orMore := false
	for _, singleCode := range strings.Split(demand, "") {
		if singleCode == "+" {
			orMore = true
			continue
		}
		if singleCode == "-" {
			orLess = true
			continue
		}
		val := ehex.FromCode(singleCode).Value()
		if val >= 0 {
			demandedCodes = append(demandedCodes, val)
		}
	}
	if orLess {
		for m := sliceMin(demandedCodes); m > 0; m-- {
			demandedCodes = append(demandedCodes, m)
		}
	}
	if orMore {
		for m := sliceMax(demandedCodes); m < 33; m++ {
			demandedCodes = append(demandedCodes, m)
		}
	}
	for _, v := range demandedCodes {
		if v == codeValue {
			return true
		}
	}
	return false
}

func sliceMin(sl []int) int {
	if len(sl) == 0 {
		return 0
	}
	m := sl[0]
	for _, i := range sl {
		if i < m {
			m = i
		}
	}
	return m
}
func sliceMax(sl []int) int {
	if len(sl) == 0 {
		return 0
	}
	m := sl[0]
	for _, i := range sl {
		if i > m {
			m = i
		}
	}
	return m
}

func TradeCodes(u uwp.UWP) []Classification {
	clmap := make(map[Classification]bool)
	for _, cl := range Classify(u) {
		clmap[cl] = true
	}
	tc := []Classification{}
	for _, code := range []Classification{Ag, As, Ba, De, Fl, Ga, Hi, Ht, Ic, In, Lo, Lt, Na, Ni, Po, Ri, Va, Wa} {
		if clmap[code] {
			tc = append(tc, code)
		}
	}
	return tc
}
