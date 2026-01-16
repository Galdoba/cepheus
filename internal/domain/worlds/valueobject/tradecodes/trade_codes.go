package tradecodes

import "fmt"

type TradeCode string

const (
	Ab TradeCode = "Ab"
	Ag TradeCode = "Ag"
	An TradeCode = "An"
	As TradeCode = "As"
	Ba TradeCode = "Ba"
	Co TradeCode = "Co"
	Cp TradeCode = "Cp"
	Cs TradeCode = "Cs"
	Cx TradeCode = "Cx"
	Cy TradeCode = "Cy"
	Da TradeCode = "Da"
	De TradeCode = "De"
	Di TradeCode = "Di"
	Fa TradeCode = "Fa"
	Fl TradeCode = "Fl"
	Fo TradeCode = "Fo"
	Fr TradeCode = "Fr"
	Ga TradeCode = "Ga"
	He TradeCode = "He"
	Hi TradeCode = "Hi"
	Ho TradeCode = "Ho"
	Ht TradeCode = "Ht"
	Ic TradeCode = "Ic"
	In TradeCode = "In"
	Lk TradeCode = "Lk"
	Lo TradeCode = "Lo"
	Lt TradeCode = "Lt"
	Mi TradeCode = "Mi"
	Mr TradeCode = "Mr"
	Na TradeCode = "Na"
	Ni TradeCode = "Ni"
	Oc TradeCode = "Oc"
	Pa TradeCode = "Pa"
	Pe TradeCode = "Pe"
	Ph TradeCode = "Ph"
	Pi TradeCode = "Pi"
	Po TradeCode = "Po"
	Pr TradeCode = "Pr"
	Px TradeCode = "Px"
	Pz TradeCode = "Pz"
	Re TradeCode = "Re"
	Ri TradeCode = "Ri"
	Sa TradeCode = "Sa"
	Tr TradeCode = "Tr"
	Tu TradeCode = "Tu"
	Tz TradeCode = "Tz"
	Tw TradeCode = "Tw"
	Va TradeCode = "Va"
	Wa TradeCode = "Wa"
)

func (tc TradeCode) Description() string {
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
		return "Hot World\nThe world is at the upper temperature range of human endurance; typically in HZ-1."
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
