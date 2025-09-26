package star

func (st Star) Index() int {
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
	case "PSR":
		index += 2000
	case "NS":
		index += 3000
	case "BH":
		index += 4000
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

func FromIndex(i int) (string, string, int) {
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
