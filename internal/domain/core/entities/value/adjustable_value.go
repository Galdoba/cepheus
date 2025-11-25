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

// AdjustableValue represents a numeric value that can be modified and adjusted within defined limits
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

// New creates a new AdjustableValue with the provided options
func New(opts ...ValueOption) *AdjustableValue {
	nv := AdjustableValue{}
	for _, set := range opts {
		set(&nv)
	}
	return &nv
}

// ValueOption defines a function type for configuring AdjustableValue
type ValueOption func(*AdjustableValue)

// Base sets the initial base value for AdjustableValue
func Base(i int) ValueOption {
	return func(av *AdjustableValue) {
		av.exist = true
		av.baseValue = i
	}
}

// ValueFor configures the value type and associated properties
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

// characteristicDM calculates the dice modifier for characteristic values
func characteristicDM(value int) int {
	return max(-2, (value/3)-2)
}

// skillDM calculates the dice modifier for skill values
func skillDM(value int) int {
	return value
}

// BaseValue returns the original unmodified value
func (av *AdjustableValue) BaseValue() int {
	if !av.exist {
		return av.ebsenceValue
	}
	return av.baseValue
}

// Value returns the current modified value within defined limits
func (av *AdjustableValue) Value() int {
	if !av.exist {
		return av.ebsenceValue
	}
	return minmax(av.baseValue+av.modifier, 0, av.highLimit)
}

// AdjustedValue returns the value with adjustments applied
func (av *AdjustableValue) AdjustedValue() int {
	if !av.exist {
		return av.ebsenceValue
	}
	adj := minmax(av.adjustment, minAdj, maxAdj)
	return minmax(av.Value()+adj, 0, av.highLimit)
}

// DM returns the dice modifier based on current value
func (av *AdjustableValue) DM() int {
	if !av.exist || av.dmFunc == nil {
		return av.ebsenceValue
	}
	return av.dmFunc(av.Value())
}

// minmax constrains a value between minimum and maximum bounds
func minmax(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// String returns a formatted string representation of the value
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

// Unstring parses a string to populate the AdjustableValue fields
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

func (av *AdjustableValue) ChangeBase(change int) int {
	initial := av.baseValue
	av.baseValue = minmax(av.baseValue+change, 0, av.highLimit)
	return av.baseValue - initial
}

func (av *AdjustableValue) ChangeModifier(change int) int {
	initial := av.modifier
	av.modifier = av.modifier + change
	return av.modifier - initial
}

func (av *AdjustableValue) Adjust(change int) int {
	initial := av.adjustment
	av.adjustment = minmax(av.adjustment+change, minAdj, maxAdj)
	return av.adjustment - initial
}
