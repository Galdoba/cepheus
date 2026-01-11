package check

import "fmt"

const (
	CharacteristicCheck            = "Characteristic Check"
	SkillCheck                     = "Skill Check"
	Raw                            = "Raw Check"
	Unresolved          Resolution = "Unresolved"
	ExceptionalFailure  Resolution = "Exceptional Failure"
	NormalFailure       Resolution = "Normal Failure"
	MarginalFailure     Resolution = "Marginal Failure"
	ExceptionalSuccess  Resolution = "Exceptional Success"
	NormalSuccess       Resolution = "Normal Success"
	MarginalSuccess     Resolution = "Marginal Success"
	Simple              Difficulty = 2
	Easy                Difficulty = 4
	Routine             Difficulty = 6
	Average             Difficulty = 8
	Difficult           Difficulty = 10
	VeryDifficult       Difficulty = 12
	Formidable          Difficulty = 14
	Impossible          Difficulty = 16
)

type Resolution string

type Difficulty int

type checkModifier struct {
	description string
	dm          int
}

type Resolver interface {
	Roll(string) int
}

type Check struct {
	checkType  string
	code       string
	dms        map[checkModifier]bool
	finalDM    int
	difficulty Difficulty
	effect     int
	result     int
	resolution Resolution
	useBounds  bool
	lowBound   int
	highBound  int
}

func New(opts ...CheckOption) *Check {
	ch := Check{
		checkType:  Raw,
		code:       "2d6",
		dms:        make(map[checkModifier]bool),
		finalDM:    0,
		difficulty: Average,
		resolution: Unresolved,
	}
	for _, modify := range opts {
		modify(&ch)
	}
	return &ch
}

type CheckOption func(*Check)

func WithCode(val string) CheckOption {
	return func(c *Check) {
		c.code = val
	}
}

func NewModifier(descr string, dm int) checkModifier {
	return checkModifier{description: descr, dm: dm}
}

func WithMods(mods ...checkModifier) CheckOption {
	return func(c *Check) {
		for _, mod := range mods {
			c.dms[mod] = false
		}
	}
}

func WithModsApplied(mods ...checkModifier) CheckOption {
	return func(c *Check) {
		for _, mod := range mods {
			c.dms[mod] = true
		}
	}
}

func WithDifficulty(val Difficulty) CheckOption {
	return func(c *Check) {
		c.difficulty = val
	}
}

func WithBounds(low, high int) CheckOption {
	return func(c *Check) {
		c.useBounds = true
		c.lowBound = low
		c.highBound = high
	}
}

func (c *Check) Resolve(r Resolver) error {
	if c.resolution != Unresolved {
		return fmt.Errorf("check was already resolved")
	}
	for mod, applied := range c.dms {
		if applied {
			c.finalDM += mod.dm
		}
	}
	roll := r.Roll(c.code)
	roll += c.finalDM
	if c.useBounds {
		roll = minmax(roll, c.lowBound, c.highBound)
	}
	c.result = roll
	c.effect = roll - int(c.difficulty)
	switch minmax(c.effect, -6, 6) {
	case -6:
		c.resolution = ExceptionalFailure
	case -5, -4, -3, -2:
		c.resolution = NormalFailure
	case -1:
		c.resolution = MarginalFailure
	case 0:
		c.resolution = MarginalSuccess
	case 1, 2, 3, 4, 5:
		c.resolution = NormalSuccess
	case 6:
		c.resolution = ExceptionalSuccess
	default:
		panic(fmt.Sprintf("unattended effect value: %v", c.effect))
	}
	return nil
}

func minmax(i, min, max int) int {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}

type outcome struct {
	Effect     int
	Roll       int
	Resolution Resolution
	Success    bool
	Err        error
}

func (c *Check) Outcome() outcome {
	if c.resolution == Unresolved {
		return outcome{
			Err: fmt.Errorf("check unresolved"),
		}
	}
	o := outcome{
		Effect:     c.effect,
		Roll:       c.result,
		Resolution: c.resolution,
	}
	switch c.resolution {
	case ExceptionalSuccess, NormalSuccess, MarginalSuccess:
		o.Success = true
	}
	return o
}
