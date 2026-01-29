package interpolate

import (
	"fmt"
	"strconv"
)

func Index(Type, SubType, Class string) int {
	index := 0
	switch Class {
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
	case "PSR":
		index += 2000
	case "NS":
		index += 3000
	case "BH":
		index += 4000
	}

	switch Type {
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
	n, _ := strconv.Atoi(SubType)
	index += n
	return index
}

func FromIndex(i int) (string, string, string) {
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
	Type := ""
	switch noclass / 10 {
	case 1:
		Type = "O"
		if cl == "BD" {
			Type = "L"
		}
	case 2:
		Type = "B"
		if cl == "BD" {
			Type = "T"
		}
	case 3:
		Type = "A"
		if cl == "BD" {
			Type = "Y"
		}
	case 4:
		Type = "F"
	case 5:
		Type = "G"
	case 6:
		Type = "K"
	case 7:
		Type = "M"
	}
	sub := fmt.Sprintf("%v", noclass%10)
	switch Type {
	case "O", "B", "A", "F", "G", "K", "M", "L", "T", "Y":
	default:
		sub = ""
	}
	return Type, sub, cl
}
