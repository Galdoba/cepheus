package star

// Star represents the physical and astronomical properties of a star, including its spectral
// classification, temperature, mass, luminosity, and key orbital zone boundaries
// (habitable zone, snow line, etc.) used for stellar characterization and exoplanetary system modeling.
type Star struct {
	//The stellar class (also called spectral type) is a letter code (O, B, A, F, G, K, M, L, T, Y)
	//that categorizes a star primarily by its effective temperature and color.
	//This field might contain 'BD' value as placeholder for L/T/Y classes until details are provided.
	StellarClass string

	//NumericalSubClass is a digit from 0 to 9 that provides a finer temperature subdivision within
	//a given StellarClass (spectral type).
	NumeriacalSubClass string

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
	LuminocityClass string

	//The effective surface temperature of the star measured in Kelvin. It determines the star’s color,
	//spectral features, and the wavelength distribution of its emitted radiation.
	Temperature_K int

	//The mass of the star expressed in solar masses (SM). One solar mass (M☉) is
	//the mass of the Sun (~1.989×10³⁰ kg). This value strongly influences the star’s gravity,
	//internal pressure, fusion rate, and lifespan.
	Mass_SM float64

	//The total radiant power (energy output per second) of the star, expressed in solar units (SU).
	//One solar unit (L☉) equals the Sun’s luminosity (~3.828×10²⁶ W).
	//It determines how much energy planets receive.
	Luminosity_SU float64

	//The inner boundary of the star’s planetary system or the region of interest for planet
	//formation/dynamics, measured in Astronomical Units (AU). One AU is the average
	//Earth–Sun distance (~149.6 million km). This limit might mark where temperatures become too high
	//for certain processes or where material cannot survive.
	InnerLimit_AU float64

	//The inner edge of the classical habitable zone (the “Goldilocks zone”) in AU. At this distance,
	//a rocky planet with a suitable atmosphere could theoretically support liquid water
	//on its surface before runaway greenhouse effects make it too hot.
	HabitableZone_AU_Low float64

	//The outer edge of the classical habitable zone in AU. Beyond this distance, even with
	//a thick atmosphere and greenhouse warming, liquid water cannot persist because surface
	//temperatures fall below freezing.
	HabitableZone_AU_High float64

	//The distance from the star where volatile compounds (primarily water, but also ammonia and
	//methane) condense into solid ice grains. Inside this line, these compounds remain gaseous;
	//outside, they form icy planetesimals. This line is crucial for planet formation models.
	SnowLine_AU float64

	//The outer boundary of the star’s planetary system or the region of interest for stellar
	//dynamics, often marking where the star’s gravitational influence becomes weak, where planet
	//formation effectively stops, or where the protoplanetary disk fades.
	OuterLimit_AU float64
}

func New() *Star {
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
