package global

type SpaceGlobal struct {
	X int
	Y int
}

func NewSpaceGlobal(x, y int) SpaceGlobal {
	return SpaceGlobal{X: x, Y: y}
}
