package value

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	valueTypeCharacteristic = "characteristic"
	valueTypeSkill          = "skill"
	minAdj                  = -6
	maxAdj                  = 6
)

type AdjustableValue struct {
	exist        bool
	ebsenceValue int
	baseValue    int
	modifier     int
	adjustment   int
	highLimit    int
	valueType    string
	dmFunc       func(int) int
}

func New(opts ...ValueOption) *AdjustableValue {
	nv := AdjustableValue{}
	for _, set := range opts {
		set(&nv)
	}
	return &nv
}

type ValueOption func(*AdjustableValue)

func Base(i int) ValueOption {
	return func(av *AdjustableValue) {
		av.exist = true
		av.baseValue = i
	}
}

func ValueFor(valueType string) ValueOption {
	return func(av *AdjustableValue) {
		switch valueType {
		case valueTypeCharacteristic:
			av.valueType = valueType
			av.dmFunc = characteristicDM
			av.highLimit = 15
			av.ebsenceValue = -3
		case valueTypeSkill:
			av.valueType = valueType
			av.dmFunc = skillDM
			av.highLimit = 6
			av.ebsenceValue = -3
		}
	}
}

func defaultDMFunc(valueType string) func(int) int {
	switch valueType {
	case valueTypeCharacteristic:

	}
	return nil
}

func characteristicDM(value int) int {
	return max(-2, (value/3)-2)
}

func skillDM(value int) int {
	return value
}

func (av *AdjustableValue) BaseValue() int {
	if !av.exist {
		return av.ebsenceValue
	}
	return av.baseValue
}

func (av *AdjustableValue) Value() int {
	if !av.exist {
		return av.ebsenceValue
	}
	return minmax(av.baseValue+av.modifier, 0, av.highLimit)
}

func (av *AdjustableValue) AdjustedValue() int {
	if !av.exist {
		return av.ebsenceValue
	}
	adj := minmax(av.adjustment, minAdj, maxAdj)
	return minmax(av.Value()+adj, 0, av.highLimit)
}

func (av *AdjustableValue) DM() int {
	if !av.exist || av.dmFunc == nil {
		return av.ebsenceValue
	}
	return av.dmFunc(av.Value())
}

func minmax(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (av *AdjustableValue) String() string {
	s := ""
	val := av.BaseValue()
	mVal := av.Value()
	switch mVal == val {
	case true:
		s += fmt.Sprintf("%d", val)
	case false:
		s += fmt.Sprintf("%d/%d", mVal, val)
	}
	dm := av.DM()
	switch dm >= 0 {
	case true:
		s += fmt.Sprintf(" (+%v)", dm)
	case false:
		s += fmt.Sprintf(" (%v)", dm)
	}
	return s
}

func (av *AdjustableValue) Unstring(s string) error {
	parts := strings.Split(s, " ")
	valParts := strings.Split(parts[0], "/")
	for len(valParts) < 2 {
		valParts = append(valParts, valParts[0])
	}
	if bv, err := strconv.Atoi(valParts[1]); err != nil {
		return fmt.Errorf("imposible to parse base value '%s'", s)
	} else {
		av.baseValue = bv
	}
	if mv, err := strconv.Atoi(valParts[1]); err != nil {
		return fmt.Errorf("imposible to parse modified value '%s'", s)
	} else {
		av.modifier = mv - av.baseValue
	}
	return nil
}
