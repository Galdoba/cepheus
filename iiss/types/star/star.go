package star

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/gametable"
)

type Star struct {
	Type           string   `json:"type,omitempty"`
	SubType        *int     `json:"subtype,omitempty"`
	Class          string   `json:"class,omitempty"`
	Mass           float64  `json:"mass,omitempty"`
	Temperature    int      `json:"temperature,omitempty"`
	Diameter       float64  `json:"diameter,omitempty"`
	Luminocity     float64  `json:"luminocity,omitempty"`
	Designation    *string  `json:"designation,omitempty"`
	OrbitN         *float64 `json:"orbit#,omitempty"`
	SystemAU       *float64 `json:"au primary,omitempty"`
	Eccentricity   *float64 `json:"eccentricity,omitempty"`
	Age            float64  `json:"age,omitempty"`
	realetdPrimary *Star
}

func Generate(dp *dice.Dicepool, knownData ...KnownStarData) (Star, error) {
	st := Star{}
	for _, add := range knownData {
		add(&st)
	}
	if st.Type+st.Class == "" && st.realetdPrimary == nil {
		tpe, cls, err := StarTypeAndClassDetermination(dp)
		if err != nil {
			return st, fmt.Errorf("failed to determine type and class of primary star: %v", err)
		}
		st.Type = tpe
		st.Class = cls

		st.SubType, err = StarSubTypeDetermination(dp, st)
		if err != nil {
			return st, fmt.Errorf("failed to determine subtype of the star: %v", err)
		}
	}
	mass := 0.0
	switch st.Class {
	case "BD":
		st.Type, st.SubType, mass = bdTypeDetails(dp)
	default:
		mass = adjust(dp, massByIndex(st.index()))
	}
	st.Mass = roundFloat(mass)
	diam := adjust(dp, diamByIndex(st.index()))
	st.Diameter = roundFloat(diam)
	temp := adjust(dp, tempByIndex(st.index()))
	st.Temperature = int(temp)
	st.Luminocity = roundFloat(calculateLuminosity(st.Diameter, temp))
	st.Age = generateAge(dp, st)
	return st, nil
}

func adjust(dp *dice.Dicepool, fl float64) float64 {
	adj := float64(dp.Flux())*0.01 + 1.0 + float64(dp.Flux())*0.001
	return fl * adj
}

func calculateLuminosity(diameter, temperature float64) float64 {
	diameterComponent := math.Pow(diameter, 2)
	temperatureComponent := math.Pow(temperature/5772, 4)
	luminosity := diameterComponent * temperatureComponent
	return luminosity
}

func generateAge(dp *dice.Dicepool, st Star) float64 {
	age := 0.0
	msls := mainSequanceLifespan(st.Mass)
	sgls := msls / (4.0 / st.Mass)
	glls := msls / (10.0 / math.Pow(st.Mass, 3))

	switch st.Class {
	case "Ia", "Ib", "II", "V", "VI":
		switch st.Mass > 0.9 {
		case false:
			age = smallStarAge(dp)
		case true:
			age = msls
		}
	case "IV":
		age = msls + sgls
	case "III":
		age = msls + sgls + glls
	case "BD":
		age = smallStarAge(dp)
	default:
		mass := originalMass(dp, st.Mass)
		msls = mainSequanceLifespan(mass)
		sgls = msls / (4.0 / st.Mass)
		glls = msls / (10.0 / math.Pow(mass, 3))
		finalAge := msls + sgls + glls
		switch st.Class {
		case "D", "NS", "BH":
			age = finalAge + smallStarAge(dp)
		case "PSR":
			age = 0.1/(float64(dp.Sum("2d10"))) + finalAge
		case "Protostar":
			age = 0.01 / float64(dp.Sum("2d10"))
		}

	}
	age = adjust(dp, age*variance(dp))
	if st.Mass < 4.7 && age < 0.01 {
		age = 0.01
	}
	return roundFloat(age)
}

func bdTypeDetails(dp *dice.Dicepool) (string, *int, float64) {
	mass := float64(dp.Sum("1d6"))/100.0 + float64(dp.Sum("4d6")-1)/1000.0
	vals := []float64{0.080, 0.076, 0.072, 0.068, 0.064, 0.060, 0.058, 0.056, 0.054, 0.052, 0.050, 0.048, 0.046, 0.044, 0.042, 0.040, 0.037, 0.034, 0.031, 0.028, 0.025, 0.022, 0.019, 0.016, 0.014, 0.013, 0.012, 0.011, 0.010, 0.009}
	index := -1
	for i, v := range vals {
		if v <= mass {
			index = i
			break
		}
	}
	tp := ""
	switch index / 10 {
	case 0:
		tp = "L"
	case 1:
		tp = "T"
	case 2:
		tp = "Y"
	}
	stp := index % 10
	return tp, &stp, mass
}

func originalMass(dp *dice.Dicepool, currentMass float64) float64 {
	return float64(dp.Sum("1d3")) / 2 * currentMass
}

func mainSequanceLifespan(mass float64) float64 {
	denominator := math.Pow(mass, 2.5)
	msa := 10.0 / denominator
	return msa
}

func smallStarAge(dp *dice.Dicepool) float64 {
	age := float64(dp.Sum("1d6"))*2 + float64(dp.Sum("1d3")-2) + variance(dp)
	age = adjust(dp, age)
	return roundFloat(age)
}

func largeStarAge(dp *dice.Dicepool, msa float64) float64 {
	age := msa * variance(dp)
	age = adjust(dp, age)
	return roundFloat(age)
}

func variance(dp *dice.Dicepool) float64 {
	return float64(dp.Sum("1d100")) / 100.0
}

func roundFloat(x float64) float64 {
	// Сначала округляем до миллионных (6 знаков)
	intermediate := math.Round(x*1e6) / 1e6

	n := 6
	if x > 0.001 {
		n = 3
	}
	if x > 10 {
		n = 1
	}

	// Затем округляем до требуемого количества знаков
	pow := math.Pow(10, float64(n))
	return math.Round(intermediate*pow) / pow
}

type KnownStarData func(*Star)

func KnownType(sType string) KnownStarData {
	return func(s *Star) {
		s.Type = sType
	}
}

func KnownClass(class string) KnownStarData {
	return func(s *Star) {
		s.Class = class
	}
}

func StarTypeAndClassDetermination(dp *dice.Dicepool) (string, string, error) {
	giantsTable, err := gametable.NewTable("Unusual", "2d6",
		gametable.NewRollResult("8-", "III", nil),
		gametable.NewRollResult("9..10", "II", nil),
		gametable.NewRollResult("11", "Ib", nil),
		gametable.NewRollResult("12+", "Ia", nil),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create giantsTable: %v", err)
	}
	specialTable, err := gametable.NewTable("Special", "2d6",
		gametable.NewRollResult("5-", "VI", nil),
		gametable.NewRollResult("6..8", "IV", nil),
		gametable.NewRollResult("9..10", "III", nil),
		gametable.NewRollResult("11+", "Giants", giantsTable),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create specialTable: %v", err)
	}
	hotTable, err := gametable.NewTable("Hot", "2d6",
		gametable.NewRollResult("9-", "A", nil),
		gametable.NewRollResult("10..11", "B", nil),
		gametable.NewRollResult("12+", "O", nil),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create hotTable: %v", err)
	}
	typeTable, err := gametable.NewTable("Type", "2d6",
		gametable.NewRollResult("2-", "Unusual", specialTable),
		gametable.NewRollResult("3..6", "M", nil),
		gametable.NewRollResult("7..8", "K", nil),
		gametable.NewRollResult("9..10", "G", nil),
		gametable.NewRollResult("11", "F", nil),
		gametable.NewRollResult("12+", "Hot", hotTable),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create typeTable: %v", err)
	}
	starType, err := typeTable.Roll(dp)
	if err != nil {
		return "", "", fmt.Errorf("table roll failed: %v", err)
	}
	done := false
	class := "V"
	for !done {
		switch starType {
		case "O", "B", "A", "F", "G", "K", "M":
			done = true
		case "Ia", "Ib", "II", "III":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", "", fmt.Errorf("table roll failed: %v", err)
			}
		case "IV":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", "", fmt.Errorf("table roll failed: %v", err)
			}
			switch starType {
			case "M":
				starType = "IV"
			case "O":
				starType = "B"
			}

		case "VI":
			class = starType
			starType, err = typeTable.WithMod(1).Roll(dp)
			if err != nil {
				return "", "", fmt.Errorf("table roll failed: %v", err)
			}
			switch starType {
			case "F":
				starType = "G"
			case "A":
				starType = "B"
			}
		}
	}
	return starType, class, nil
}

func StarSubTypeDetermination(dp *dice.Dicepool, st Star) (*int, error) {
	subtypeTable := &gametable.GameTable{}
	err := fmt.Errorf("table not created")
	switch st.Type {
	case "M":
		subtypeTable, err = gametable.NewTable("m-type", "2d6",
			gametable.NewRollResult("2", "8", nil),
			gametable.NewRollResult("3", "6", nil),
			gametable.NewRollResult("4", "5", nil),
			gametable.NewRollResult("5", "4", nil),
			gametable.NewRollResult("6", "0", nil),
			gametable.NewRollResult("7", "2", nil),
			gametable.NewRollResult("8", "1", nil),
			gametable.NewRollResult("9", "3", nil),
			gametable.NewRollResult("10", "5", nil),
			gametable.NewRollResult("11", "7", nil),
			gametable.NewRollResult("12", "9", nil),
		)
	case "O", "B", "A", "F", "G", "K":
		subtypeTable, err = gametable.NewTable("numeric", "2d6",
			gametable.NewRollResult("2", "0", nil),
			gametable.NewRollResult("3", "1", nil),
			gametable.NewRollResult("4", "3", nil),
			gametable.NewRollResult("5", "5", nil),
			gametable.NewRollResult("6", "7", nil),
			gametable.NewRollResult("7", "9", nil),
			gametable.NewRollResult("8", "8", nil),
			gametable.NewRollResult("9", "6", nil),
			gametable.NewRollResult("10", "4", nil),
			gametable.NewRollResult("11", "2", nil),
			gametable.NewRollResult("12", "0", nil),
		)
	default:
		return nil, fmt.Errorf("what shaall we do with %v?", st.Type)
	}
	r, err := subtypeTable.Roll(dp)
	if err != nil {
		return nil, fmt.Errorf("sybtype table roll: %v", err)
	}
	n, _ := strconv.Atoi(r)
	return &n, nil
}

// func (st Star) String() string {
// 	return fmt.Sprintf("%v%v %v", st.Type, *st.SubType, st.Class)
// }

type mtlData struct {
	mass        float64
	temperature float64
	luminocity  float64
}

func (st Star) index() int {
	index := 0
	switch st.Class {
	case "Ia":
		index += 100
	case "Ib":
		index += 200
	case "II":
		index += 300
	case "III":
		index += 400
	case "IV":
		index += 500
	case "V":
		index += 600
	case "VI":
		index += 700
	case "VII", "D":
		index += 800
	case "BD":
		index += 900
	}

	switch st.Type {
	case "O", "L":
		index += 10
	case "B", "T":
		index += 20
	case "A", "Y":
		index += 30
	case "F":
		index += 40
	case "G":
		index += 50
	case "K":
		index += 60
	case "M":
		index += 70
	}
	if st.SubType != nil {
		index += *st.SubType
	}
	return index
}

func fromIndex(i int) (string, string, int) {
	cl := ""
	switch i / 100 {
	case 1:
		cl = "Ia"
	case 2:
		cl = "Ib"
	case 3:
		cl = "II"
	case 4:
		cl = "III"
	case 5:
		cl = "IV"
	case 6:
		cl = "V"
	case 7:
		cl = "VI"
	case 8:
		cl = "D"
	case 9:
		cl = "BD"
	}
	noclass := i % 100
	scl := ""
	switch noclass / 10 {
	case 1:
		scl = "O"
		if cl == "BD" {
			scl = "L"
		}
	case 2:
		scl = "B"
		if cl == "BD" {
			scl = "T"
		}
	case 3:
		scl = "A"
		if cl == "BD" {
			scl = "Y"
		}
	case 4:
		scl = "F"
	case 5:
		scl = "G"
	case 6:
		scl = "K"
	case 7:
		scl = "M"
	}
	sub := noclass % 10
	return cl, scl, sub
}

func massByIndex(i int) float64 {
	massMap := make(map[int]float64)

	massMap[110] = 200
	massMap[115] = 80
	massMap[120] = 60
	massMap[125] = 30
	massMap[130] = 20
	massMap[135] = 15
	massMap[140] = 13
	massMap[145] = 12
	massMap[150] = 12
	massMap[155] = 13
	massMap[160] = 14
	massMap[165] = 18
	massMap[170] = 20
	massMap[175] = 25
	massMap[179] = 30

	massMap[210] = 150
	massMap[215] = 60
	massMap[220] = 40
	massMap[225] = 25
	massMap[230] = 15
	massMap[235] = 13
	massMap[240] = 12
	massMap[245] = 10
	massMap[250] = 10
	massMap[255] = 11
	massMap[260] = 12
	massMap[265] = 13
	massMap[270] = 15
	massMap[275] = 20
	massMap[279] = 25

	massMap[310] = 130
	massMap[315] = 40
	massMap[320] = 30
	massMap[325] = 20
	massMap[330] = 14
	massMap[335] = 11
	massMap[340] = 10
	massMap[345] = 8
	massMap[350] = 8
	massMap[355] = 10
	massMap[360] = 10
	massMap[365] = 12
	massMap[370] = 14
	massMap[375] = 16
	massMap[379] = 18

	massMap[410] = 110
	massMap[415] = 30
	massMap[420] = 20
	massMap[425] = 10
	massMap[430] = 8
	massMap[435] = 6
	massMap[440] = 4
	massMap[445] = 3
	massMap[450] = 2.5
	massMap[455] = 2.4
	massMap[460] = 1.1
	massMap[465] = 1.5
	massMap[470] = 1.8
	massMap[475] = 2.4
	massMap[479] = 8

	// massMap[510] = 110
	// massMap[515] = 30
	massMap[520] = 20
	massMap[525] = 10
	massMap[530] = 4
	massMap[535] = 2.3
	massMap[540] = 2
	massMap[545] = 1.5
	massMap[550] = 1.7
	massMap[555] = 1.2
	massMap[560] = 1.5
	// massMap[565] = 1.5
	// massMap[570] = 1.8
	// massMap[575] = 2.4
	// massMap[579] = 8

	massMap[610] = 90
	massMap[615] = 60
	massMap[620] = 18
	massMap[625] = 5
	massMap[630] = 2.2
	massMap[635] = 1.8
	massMap[640] = 1.5
	massMap[645] = 1.3
	massMap[650] = 1.1
	massMap[655] = 0.9
	massMap[660] = 0.8
	massMap[665] = 0.7
	massMap[670] = 0.5
	massMap[675] = 0.16
	massMap[679] = 0.08

	massMap[710] = 2
	massMap[715] = 1.5
	massMap[720] = 0.5
	massMap[725] = 0.4
	// massMap[730] = 2.2
	// massMap[735] = 1.8
	// massMap[740] = 1.5
	// massMap[745] = 1.3
	massMap[750] = 0.8
	massMap[755] = 0.7
	massMap[760] = 0.6
	massMap[765] = 0.5
	massMap[770] = 0.4
	massMap[775] = 0.12
	massMap[779] = 0.075

	massMap[910] = 0.08
	massMap[915] = 0.06
	massMap[920] = 0.05
	massMap[925] = 0.04
	massMap[930] = 0.025
	massMap[935] = 0.013

	return interpolate(massMap, i)
}

func tempByIndex(i int) float64 {
	tempMap := make(map[int]float64)
	if i < 900 {
		i = i % 100
	}
	tempMap[10] = 50000
	tempMap[15] = 40000
	tempMap[20] = 30000
	tempMap[25] = 15000
	tempMap[30] = 10000
	tempMap[35] = 8000
	tempMap[40] = 7500
	tempMap[45] = 6500
	tempMap[50] = 6000
	tempMap[55] = 5600
	tempMap[60] = 5200
	tempMap[65] = 4400
	tempMap[70] = 3700
	tempMap[75] = 3000
	tempMap[79] = 2400

	tempMap[910] = 2400
	tempMap[915] = 1850
	tempMap[920] = 1300
	tempMap[925] = 900
	tempMap[930] = 550
	tempMap[935] = 300
	return interpolate(tempMap, i)
}

func diamByIndex(i int) float64 {
	diamMap := make(map[int]float64)

	diamMap[110] = 25
	diamMap[115] = 22
	diamMap[120] = 20
	diamMap[125] = 60
	diamMap[130] = 120
	diamMap[135] = 180
	diamMap[140] = 210
	diamMap[145] = 280
	diamMap[150] = 330
	diamMap[155] = 360
	diamMap[160] = 420
	diamMap[165] = 600
	diamMap[170] = 900
	diamMap[175] = 1200
	diamMap[179] = 1800

	diamMap[210] = 24
	diamMap[215] = 20
	diamMap[220] = 14
	diamMap[225] = 25
	diamMap[230] = 50
	diamMap[235] = 75
	diamMap[240] = 85
	diamMap[245] = 115
	diamMap[250] = 135
	diamMap[255] = 150
	diamMap[260] = 180
	diamMap[265] = 260
	diamMap[270] = 380
	diamMap[275] = 600
	diamMap[279] = 800

	diamMap[310] = 22
	diamMap[315] = 18
	diamMap[320] = 12
	diamMap[325] = 14
	diamMap[330] = 30
	diamMap[335] = 45
	diamMap[340] = 50
	diamMap[345] = 66
	diamMap[350] = 77
	diamMap[355] = 90
	diamMap[360] = 110
	diamMap[365] = 160
	diamMap[370] = 230
	diamMap[375] = 350
	diamMap[379] = 500

	diamMap[410] = 21
	diamMap[415] = 15
	diamMap[420] = 10
	diamMap[425] = 6
	diamMap[430] = 5
	diamMap[435] = 5
	diamMap[440] = 5
	diamMap[445] = 5
	diamMap[450] = 10
	diamMap[455] = 15
	diamMap[460] = 20
	diamMap[465] = 40
	diamMap[470] = 60
	diamMap[475] = 100
	diamMap[479] = 200

	// massMap[510] = 110
	// massMap[515] = 30
	diamMap[520] = 8
	diamMap[525] = 5
	diamMap[530] = 4
	diamMap[535] = 3
	diamMap[540] = 3
	diamMap[545] = 2
	diamMap[550] = 3
	diamMap[555] = 4
	diamMap[560] = 6
	// massMap[565] = 1.5
	// massMap[570] = 1.8
	// massMap[575] = 2.4
	// massMap[579] = 8

	diamMap[610] = 20
	diamMap[615] = 12
	diamMap[620] = 7
	diamMap[625] = 3.5
	diamMap[630] = 2.2
	diamMap[635] = 2
	diamMap[640] = 1.7
	diamMap[645] = 1.5
	diamMap[650] = 1.1
	diamMap[655] = 0.95
	diamMap[660] = 0.9
	diamMap[665] = 0.8
	diamMap[670] = 0.7
	diamMap[675] = 0.2
	diamMap[679] = 0.1

	diamMap[710] = 0.18
	diamMap[715] = 0.18
	diamMap[720] = 0.2
	diamMap[725] = 0.5
	// massMap[730] = 2.2
	// massMap[735] = 1.8
	// massMap[740] = 1.5
	// massMap[745] = 1.3
	diamMap[750] = 0.8
	diamMap[755] = 0.7
	diamMap[760] = 0.6
	diamMap[765] = 0.5
	diamMap[770] = 0.4
	diamMap[775] = 0.1
	diamMap[779] = 0.08

	diamMap[910] = 0.1
	diamMap[915] = 0.08
	diamMap[920] = 0.9
	diamMap[925] = 0.11
	diamMap[930] = 0.1
	diamMap[935] = 0.1
	return interpolate(diamMap, i)
}

func lumaByIndex(i int) float64 {
	lumaMap := make(map[int]float64)

	lumaMap[110] = 3400000
	lumaMap[115] = 1100000
	lumaMap[120] = 290000
	lumaMap[125] = 160000
	lumaMap[130] = 130000
	lumaMap[135] = 120000
	lumaMap[140] = 120000
	lumaMap[145] = 120000
	lumaMap[150] = 120000
	lumaMap[155] = 110000
	lumaMap[160] = 110000
	lumaMap[165] = 120000
	lumaMap[170] = 130000
	lumaMap[175] = 100000
	lumaMap[179] = 90000

	lumaMap[210] = 3200000
	lumaMap[215] = 900000
	lumaMap[220] = 140000
	lumaMap[225] = 28000
	lumaMap[230] = 22000
	lumaMap[235] = 20000
	lumaMap[240] = 20000
	lumaMap[245] = 20000
	lumaMap[250] = 20000
	lumaMap[255] = 20000
	lumaMap[260] = 21000
	lumaMap[265] = 22000
	lumaMap[270] = 24000
	lumaMap[275] = 26000
	lumaMap[279] = 19000

	lumaMap[310] = 2700000
	lumaMap[315] = 730000
	lumaMap[320] = 100000
	lumaMap[325] = 8800
	lumaMap[330] = 8000
	lumaMap[335] = 7300
	lumaMap[340] = 7000
	lumaMap[345] = 6900
	lumaMap[350] = 6800
	lumaMap[355] = 7000
	lumaMap[360] = 7800
	lumaMap[365] = 8400
	lumaMap[370] = 8800
	lumaMap[375] = 8800
	lumaMap[379] = 7300

	lumaMap[410] = 2400000
	lumaMap[415] = 510000
	lumaMap[420] = 72000
	lumaMap[425] = 1600
	lumaMap[430] = 220
	lumaMap[435] = 90
	lumaMap[440] = 70
	lumaMap[445] = 39
	lumaMap[450] = 120
	lumaMap[455] = 200
	lumaMap[460] = 260
	lumaMap[465] = 530
	lumaMap[470] = 600
	lumaMap[475] = 720
	lumaMap[479] = 1200

	// massMap[510] = 110
	// massMap[515] = 30
	lumaMap[520] = 46000
	lumaMap[525] = 1100
	lumaMap[530] = 140
	lumaMap[535] = 33
	lumaMap[540] = 25
	lumaMap[545] = 6
	lumaMap[550] = 10
	lumaMap[555] = 14
	lumaMap[560] = 23
	// massMap[565] = 1.5
	// massMap[570] = 1.8
	// massMap[575] = 2.4
	// massMap[579] = 8

	lumaMap[610] = 2200000
	lumaMap[615] = 330000
	lumaMap[620] = 35000
	lumaMap[625] = 550
	lumaMap[630] = 43
	lumaMap[635] = 15
	lumaMap[640] = 8.1
	lumaMap[645] = 3.5
	lumaMap[650] = 1.4
	lumaMap[655] = 0.78
	lumaMap[660] = 0.52
	lumaMap[665] = 0.21
	lumaMap[670] = 0.82
	lumaMap[675] = 0.0029
	lumaMap[679] = 0.00029

	lumaMap[710] = 180
	lumaMap[715] = 73
	lumaMap[720] = 29
	lumaMap[725] = 11
	// massMap[730] = 2.2
	// massMap[735] = 1.8
	// massMap[740] = 1.5
	// massMap[745] = 1.3
	lumaMap[750] = 0.73
	lumaMap[755] = 0.43
	lumaMap[760] = 0.23
	lumaMap[765] = 0.083
	lumaMap[770] = 0.027
	lumaMap[775] = 0.00072
	lumaMap[779] = 0.00019

	return interpolate(lumaMap, i)
}

func interpolate(massMap map[int]float64, index int) float64 {
	// Если значение уже есть в карте, возвращаем его
	if mass, exists := massMap[index]; exists {
		return mass
	}

	// Находим ближайшие известные индексы для интерполяции
	lowerIndex, upperIndex := findClosestIndices(massMap, index)

	// Если не найдены подходящие индексы, возвращаем 0
	if lowerIndex == -1 || upperIndex == -1 {
		return 0
	}

	// Линейная интерполяция
	lowerMass := massMap[lowerIndex]
	upperMass := massMap[upperIndex]

	return lowerMass + (float64(index-lowerIndex))*(upperMass-lowerMass)/float64(upperIndex-lowerIndex)
}

func findClosestIndices(massMap map[int]float64, target int) (int, int) {
	var lower, upper int = -1, -1

	// Ищем ближайший меньший индекс
	for i := target - 1; i >= 10; i-- {
		if _, exists := massMap[i]; exists {
			lower = i
			break
		}
	}

	// Ищем ближайший больший индекс
	for i := target + 1; i <= 939; i++ {
		if _, exists := massMap[i]; exists {
			upper = i
			break
		}
	}

	return lower, upper
}
