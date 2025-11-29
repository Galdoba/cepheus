package value

import "github.com/Galdoba/cepheus/internal/domain/core/values/modifier"

const (
	defaultAbsenceValue = -3
	defaultSkillLimit   = 6
)

type SkillValue struct {
	AdjustableValue
}

func NewSkillValue() *SkillValue {
	s := SkillValue{}
	s.ebsenceValue = defaultAbsenceValue
	s.highLimit = defaultSkillLimit
	s.modifiers = make(map[string]modifier.Modifier)
	s.dmFunc = skillDM
	return &s
}

func skillDM(val int) int {
	return val
}

func (s *SkillValue) Value() int {
	switch s.exist {
	case false:
		return s.ebsenceValue
	default:
		return s.baseValue
	}
}

func (s *SkillValue) ValueModded(mods ...string) int {
	v := s.Value()
	return v + s.SumModifiers(mods...)
}

func (s *SkillValue) Increase() {
	s.exist = true
	s.baseValue = minmax(s.baseValue+1, 0, s.highLimit)
}

func (s *SkillValue) Ensure(val int) {
	s.exist = true
	val = minmax(val, 0, s.highLimit)
	s.baseValue = minmax(val, s.baseValue, s.highLimit)
}
