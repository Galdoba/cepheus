package sector

type Sector struct {
	x, y           int
	name           string
	subsectorNames []string
}

func (s *Sector) Subsector(x, y int) int {
	if x < 1 {
		return -1
	}

	ss := 0
	return ss
}
