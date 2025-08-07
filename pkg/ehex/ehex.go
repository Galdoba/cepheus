package ehex

import (
	"fmt"
	"strings"
)

const (
	undesignated = "undesignated"
)

var stdRuneArray = []rune{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 65, 66, 67, 68, 69, 70, 71, 72, 74, 75, 76, 77, 78, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
var runes_x64 = []rune{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 65, 66, 67, 68, 69, 70, 71, 72, 74, 75, 76, 77, 78, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 63, 94, 42, 38}

type container struct {
	code           rune
	valueType      string
	runeToValueMap map[rune]int
	valueToRuneMap map[int]rune
}

type Ehex interface {
	Code() string
	Value() int
	Type() string
}

func (c *container) Code() string {
	if _, ok := c.runeToValueMap[c.code]; ok {
		return string(c.code)
	} else {
		return "*"
	}
}

func (c *container) Value() int {
	if r, ok := c.runeToValueMap[c.code]; ok {
		return r
	} else {
		return -1
	}
}

func (c *container) Type() string {
	return c.valueType
}

//Create funcs

func FromString(s string, valueType ...string) *container {
	ct := newContainer()
	ct.setValueType(valueType...)
	ct.setString(s)
	return ct
}

func FromInt(i int, valueType ...string) *container {
	ct := newContainer()
	ct.setValueType(valueType...)
	ct.setInt(i)
	return ct
}

//inner funcs

func newContainer() *container {
	ct := container{}
	ct.runeToValueMap = make(map[rune]int)
	ct.valueToRuneMap = make(map[int]rune)
	return &ct
}

func (ct *container) setValueType(valueType ...string) {
	ct.valueType = undesignated
	for _, vt := range valueType {
		ct.valueType = vt
	}

	switch ct.valueType {
	case undesignated:
		for v, r := range runes_x64 {
			ct.runeToValueMap[r] = v
			ct.valueToRuneMap[v] = r
		}
	default:
		panic(fmt.Sprintf("rune array for ehex.container type '%v' is not implemented", ct.valueType))
	}
}

func (ct *container) setString(s string) {
	if len(s) == 0 {
		panic("string argument required for ehex.container to set")
	}
	if len(s) != 1 {
		panic("string argument MUST contain 1 character for ehex.container to set")
	}
	r := []rune(s)
	ct.code = r[0]
}

func (ct *container) setInt(i int) {
	if i < 0 {
		panic("int argiment MUST NOT BE NEGATIVE")
	}
	if r, ok := ct.valueToRuneMap[i]; ok {
		ct.code = r
	} else {
		r := []rune("*")
		ct.code = r[0]
	}
}

//Independent funcs

// func Code(i int, valueType ...string) string {
// 	return ""
// }

// func Value(code string) int {
// 	return 0
// }

func Add(eh1, eh2 Ehex) *container {
	if eh1.Type() != eh2.Type() {
		panic("can't add different types of ehex.containers")
	}
	newEh := newContainer()
	newEh.setValueType(eh1.Type())
	newEh.setInt(eh1.Value() + eh2.Value())
	return newEh
}

func PrintRunes(str string) {
	s := "{"
	for _, rn := range str {
		s += fmt.Sprintf("%v, ", rn)
	}
	s = strings.TrimSuffix(s, ", ")
	s += "}"
	fmt.Print(s)
}

/*
0   0 1 2 3 4 5 6 7
1   8 9 A B C D E F
2   G H J K L M N P
3   Q R S T U V W X
4   Y Z a b c d e f
5   g h i j k l m n
6   o p q r s t u v
7   w x y z ? ^ * &

8
9
A
B
C
D
E
*/
