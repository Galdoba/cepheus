package orbit

import (
	"math"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/float"
)

type Orbit struct {
	Parent            string
	FromParent        float64
	Eccentricity      float64
	AU                float64
	Min               float64
	Max               float64
	Period            float64
	starsDM           int
	systemAge         float64
	starMass1         float64
	starMass2         float64
	isAsteroid        bool
	isStar            bool
	isProtostarSystem bool
}

func NewStar(dp *dice.Dicepool, code string, opts ...OrbitOption) *Orbit {
	or := Orbit{}
	or.Parent = StarParent(code)
	switch code {
	case "Ab", "Bb", "Cb", "Db":
		or.FromParent = float64(dp.Sum("1d61")-1)/100.0 + 0.05
	case "Ba":
		or.FromParent = float64(dp.Sum1D()-1) + float64(dp.Sum("1d100"))/100.0 - 0.5
	case "Ca":
		or.FromParent = float64(dp.Sum1D()+5) + float64(dp.Sum("1d100"))/100.0 - 0.5
		or.starsDM = 1
	case "Da":
		or.FromParent = float64(dp.Sum1D()+11) + float64(dp.Sum("1d100"))/100.0 - 0.5
		or.starsDM = 2
	case "Aa":
		for _, modify := range opts {
			modify(&or)
		}
		return &or
	}
	or.FromParent = float.Round(or.FromParent)
	if or.FromParent > 20.0 {
		panic(or.FromParent)
	}
	for _, modify := range opts {
		modify(&or)
	}
	dm := or.starsDM
	if or.systemAge > 1 && or.FromParent < 1 {
		dm--
	}
	if or.isAsteroid {
		dm++
	}
	if or.isStar {
		dm = dm + 2
	}
	if or.isProtostarSystem {
		dm = dm + 2
	}
	or.Eccentricity = eccentricity(dp, dm)
	if or.FromParent == 0 {
		or.Eccentricity = 0
	}
	or.AU = OrbitN_To_AU(or.FromParent)
	or.Min = OrbitN_To_AU(or.FromParent * (1.0 - or.Eccentricity))
	or.Min = OrbitN_To_AU(or.FromParent * (1.0 + or.Eccentricity))
	or.Period = Period(or.starMass1, or.starMass2, or.FromParent)
	return &or
}

type OrbitOption func(*Orbit)

func StarsDM(dm int) OrbitOption {
	return func(o *Orbit) {
		o.starsDM = dm
	}
}

func IsStar(isStar bool) OrbitOption {
	return func(o *Orbit) {
		o.isStar = isStar
	}
}

func SystemAge(age float64) OrbitOption {
	return func(o *Orbit) {
		o.systemAge = age
	}
}

func IsAsteroid(ast bool) OrbitOption {
	return func(o *Orbit) {
		o.isAsteroid = ast
	}
}

func StarMass(i int, mass float64) OrbitOption {
	return func(o *Orbit) {
		switch i {
		case 1:
			o.starMass1 = mass
		case 2:
			o.starMass2 = mass
		}
	}
}

func eccentricity(dp *dice.Dicepool, dm int) float64 {
	fr := dp.Sum("2d6") + dm
	if fr <= 5 {
		return -0.001 + float64(dp.Sum1D())/1000.0
	}
	if fr >= 12 {
		return 0.3 + float64(dp.Sum("2d6"))/20.0
	}
	switch fr {
	case 6, 7:
		return 0.00 + float64(dp.Sum("1d6"))/200.0
	case 8, 9:
		return 0.03 + float64(dp.Sum("1d6"))/100.0
	case 10:
		return 0.05 + float64(dp.Sum("1d6"))/20.0
	case 11:
		return 0.05 + float64(dp.Sum("2d6"))/20.0
	}
	return 0
}

func StarParent(code string) string {
	switch code {
	case "Ab":
		return "Aa"
	case "Ba":
		return "Aa"
	case "Bb":
		return "Ba"
	case "Ca":
		return "Aa"
	case "Cb":
		return "Ca"
	case "Da":
		return "Aa"
	case "Db":
		return "Da"
	}
	return ""
}

func OrbitN_To_AU(orbitN float64) float64 {
	if orbitN < 0 {
		return 0
	}
	if orbitN > 20 {
		return 78700.0
	}
	fixed := make(map[float64]float64)
	fixed[0] = 0.0
	fixed[1] = 0.4
	fixed[2] = 0.7
	fixed[3] = 1.0
	fixed[4] = 1.6
	fixed[5] = 2.8
	fixed[6] = 5.2
	fixed[7] = 10.0
	fixed[8] = 20.0
	fixed[9] = 40.0
	fixed[10] = 77.0
	fixed[11] = 154.0
	fixed[12] = 308.0
	fixed[13] = 615.0
	fixed[14] = 1230.0
	fixed[15] = 2500.0
	fixed[16] = 4900.0
	fixed[17] = 9800.0
	fixed[18] = 19500.0
	fixed[19] = 39500.0
	fixed[20] = 78700.0
	orInt := int(orbitN)
	low := float64(orInt)
	high := float64(orInt + 1)
	fraction := orbitN - low
	fracMod := fixed[high] - fixed[low]
	return float.Round((fracMod * fraction) + fixed[low])
}

func OrbitN_To_Mkm(orbitN float64) float64 {
	au := OrbitN_To_AU(orbitN)
	return AU_To_Mkm(au)
}

func AU_To_Mkm(au float64) float64 {
	return float.Round(au * 149.598)
}

// years
func Period(m1, m2, distAU float64) float64 {
	if m1+m2 == 0 {
		return 0
	}
	years := math.Pow((math.Pow(distAU, 3.0))/(m1+m2), 0.5)
	return float.Round(years)
}

type OrbitAllowanceSegment struct {
	Start float64
	End   float64
}

type OrbitAllowanceSequence struct {
	segments []OrbitAllowanceSegment
}

var Full = OrbitAllowanceSegment{0.005, 20.0}

var InitialPrimarySequance = OrbitAllowanceSequence{
	segments: []OrbitAllowanceSegment{Full},
}

func InitialSequance(mao float64, end float64) OrbitAllowanceSequence {
	return OrbitAllowanceSequence{
		segments: []OrbitAllowanceSegment{
			{
				Start: mao,
				End:   end,
			},
		},
	}
}

func SubtractSubSequence(sequence OrbitAllowanceSequence, center, width float64) OrbitAllowanceSequence {
	remove := OrbitAllowanceSegment{
		Start: center - width,
		End:   center + width,
	}

	var resultSegments []OrbitAllowanceSegment

	for _, seg := range sequence.segments {
		// Проверяем, есть ли перекрытие между сегментом и удаляемым диапазоном
		if seg.End <= remove.Start || seg.Start >= remove.End {
			// Нет перекрытия - добавляем весь сегмент
			resultSegments = append(resultSegments, seg)
			continue
		}

		// Есть перекрытие - разделяем сегмент на части
		if seg.Start < remove.Start {
			// Добавляем левую часть
			resultSegments = append(resultSegments, OrbitAllowanceSegment{
				Start: seg.Start,
				End:   remove.Start,
			})
		}

		if seg.End > remove.End {
			// Добавляем правую часть
			resultSegments = append(resultSegments, OrbitAllowanceSegment{
				Start: remove.End,
				End:   seg.End,
			})
		}
	}

	return OrbitAllowanceSequence{segments: resultSegments}
}

func CenterWidth(segment OrbitAllowanceSegment) (float64, float64) {
	center := (segment.Start + segment.End) / 2
	width := (segment.End - segment.Start) / 2
	return center, width
}
