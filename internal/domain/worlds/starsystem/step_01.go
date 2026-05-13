package starsystem

const (
	Empty        SubsectorType = "Empty"
	Scattered    SubsectorType = "Scattered"
	Dispesed     SubsectorType = "Dispesed"
	Average      SubsectorType = "Average"
	Crowded      SubsectorType = "Crowded"
	Dense        SubsectorType = "Dense"
	Realistic    SystemType    = "Realistic"
	SemiRelistic SystemType    = "Semi-Realistic"
	Fantastic    SystemType    = "Fantastic"
)

type SubsectorType string
type SystemType string

func (b *Builder) Step_01(ss *StarSystem) error {
	r := b.dice.MustRoll("1d100")
	starPresent := false
	switch b.cfg.RegionType {
	case Empty:
		starPresent = r <= 5
	case Scattered:
		starPresent = r <= 20
	case Dispesed:
		starPresent = r <= 35
	case "", Average:
		starPresent = r <= 50
	case Crowded:
		starPresent = r <= 60
	case Dense:
		starPresent = r <= 75
	}
	b.print("object presense: %v\n", starPresent)
	switch starPresent {
	case true:
		b.nextStep = step_02
	case false:
		b.nextStep = done
	}
	return nil
}
