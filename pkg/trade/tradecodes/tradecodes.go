package tradecodes

import (
	"slices"
	"strings"

	"github.com/Galdoba/cepheus/pkg/uwp"
)

const (
	Ag = "Ag"
	As = "As"
	Ba = "Ba"
	De = "De"
	Fl = "Fl"
	Ga = "Ga"
	Hi = "Hi"
	Ht = "Ht"
	Ic = "Ic"
	In = "In"
	Lo = "Lo"
	Lt = "Lt"
	Na = "Na"
	Ni = "Ni"
	Po = "Po"
	Ri = "Ri"
	Va = "Va"
	Wa = "Wa"
)

type TradeCode struct {
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Rules       string `json:"rules,omitempty"`
}

func newTC(code string) TradeCode {
	return TradeCode{Code: code}
}

func GenerateFromUWP(u uwp.UWP) []string {
	codes := []string{}
	validations := tcUwpValidations()
codeLoop:
	for code, requirements := range validations {
	requirementLoop:
		for i, req := range requirements {
			if req == "" {
				continue requirementLoop
			}
			dataKey := ""
			switch i {
			case 0:
				dataKey = uwp.Size
			case 1:
				dataKey = uwp.Atmo
			case 2:
				dataKey = uwp.Hydr
			case 3:
				dataKey = uwp.Pops
			case 4:
				dataKey = uwp.Govr
			case 5:
				dataKey = uwp.Laws
			case 6:
				dataKey = uwp.TL
			}
			value := u.Data[dataKey].Code
			if !strings.Contains(req, value) {
				continue codeLoop
			}
		}
		codes = append(codes, code)
	}
	slices.Sort(codes)
	return codes
}

func tcUwpValidations() map[string][]string { //map[trade_code][]{valid_uwp_value}  /* = any
	validations := make(map[string][]string)
	validations[Ag] = []string{"", "456789", "45678", "567", "", "", ""}
	validations[As] = []string{"0", "0", "0", "", "", "", ""}
	validations[Ba] = []string{"", "", "", "0", "0", "0", ""}
	validations[De] = []string{"", "23456789", "0", "", "", "", ""}
	validations[Fl] = []string{"", "ABCDEFG", "123456789A", "", "", "", ""}
	validations[Ga] = []string{"678", "568", "567", "", "", "", ""}
	validations[Hi] = []string{"", "", "", "9ABCDE", "", "", ""}
	validations[Ht] = []string{"", "", "", "", "", "", "CDEFGHJ"}
	validations[Ic] = []string{"", "01", "123456789A", "", "", "", ""}
	validations[In] = []string{"", "012479ABC", "", "9ABCDE", "", "", ""}
	validations[Lo] = []string{"", "", "", "123", "", "", ""}
	validations[Lt] = []string{"", "", "", "123456789ABCDE", "", "", "012345"}
	validations[Na] = []string{"", "0123", "0123", "6789ABCDE", "", "", ""}
	validations[Ni] = []string{"", "", "", "456", "", "", ""}
	validations[Po] = []string{"", "2345", "0123", "", "", "", ""}
	validations[Ri] = []string{"", "68", "", "678", "456789", "", ""}
	validations[Va] = []string{"", "0", "", "", "", "", ""}
	validations[Wa] = []string{"", "3456789DEF", "A", "", "", "", ""}
	return validations
}
