package star

import (
	"fmt"
	"strconv"
	"strings"
)

type starClassData struct {
	str string
	cl  string
	ty  string
	sub string
}

func (scd starClassData) StarData() (string, string, string, string) {
	return scd.ty, scd.sub, scd.cl, scd.str
}

func FromStellar(stellar string) []starClassData {
	if stellar == "" {
		return []starClassData{starClassData{}}
	}
	classes := []string{"BD", " Ia", " Ib", " III", " II", " IV", " VII", " VI", " V", " D", "D", "BH", "PSR", "NS", ""}
	types := []string{"", "O", "B", "A", "F", "G", "K", "M", "L", "T", "Y"}
	subtypes := []string{"", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	stars := []starClassData{}
	stellar = strings.TrimSpace(stellar)
mainLoop:
	for len(stars) < 10 {
		for _, c := range classes {
			for _, t := range types {
				for _, s := range subtypes {
					if t+s+c == "" {
						fmt.Println(stellar, stars, "<===")
						panic(0)
					}
					star := fmt.Sprintf("%v%v%v", t, s, c)
					if strings.HasPrefix(stellar, star) {
						stars = append(stars, starClassData{str: t + s + c, cl: c, ty: t, sub: s})
						stellar = strings.TrimPrefix(stellar, star)
						stellar = strings.TrimPrefix(stellar, " ")
						if stellar == "" {
							break mainLoop
						}
						continue mainLoop
					}
				}
			}
		}
	}
	return stars
}

func Disassemle(starStr string) (starClassData, error) {
	sd := FromStellar(starStr)
	switch len(sd) {
	case 0:
		return starClassData{}, fmt.Errorf("bad data: %v", starStr)
	case 1:
		return sd[0], nil
	default:
		return starClassData{}, fmt.Errorf("multiple stars provided: %v", starStr)
	}
}

func KnownStellar(stellarData string) KnownStarData {
	return func(s *Star) {
		if stellarData == "" {
			return
		}
		sd, err := Disassemle(stellarData)
		if err != nil {
			panic(err)
		}
		s.Class = strings.TrimSpace(sd.cl)
		s.Type = sd.ty
		n, err := strconv.Atoi(sd.sub)
		switch err {
		case nil:
			s.SubType = &n
		default:
		}

	}
}
