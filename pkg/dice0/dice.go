package dice

type dice struct {
	edges  int
	result int
}

type Mod interface {
	Mod() int
	ModKey() string
}

const (
	CheckEasy          = -4
	CheckRoutine       = -6
	CheckAverage       = -8
	CheckDifficult     = -10
	CheckVeryDifficult = -12
	CheckFormidable    = -14
	CheckImposible     = -16
)
