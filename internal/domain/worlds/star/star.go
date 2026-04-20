package star

type Star struct {
	StellarClass          string // O, B, A, F, G, K, M, L, T, Y
	NumeriacalSubClass    string // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	LuminocityClass       string // Ia, Ib, II, III, IV, V, VI, D
	Temperature_K         int
	Mass_SM               float64
	Luminosity_SU         float64
	InnerLimit_AU         float64
	HabitableZone_AU_Low  float64
	HabitableZone_AU_High float64
	SnowLine_AU           float64
	OuterLimit_AU         float64
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
