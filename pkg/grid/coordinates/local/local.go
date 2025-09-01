package local

type SpaceSectorLocal struct {
	SectorX int
	SectorY int
	X       int
	Y       int
}

func NewLocal(sx, sy, lx, ly int) SpaceSectorLocal {
	return SpaceSectorLocal{SectorX: sx, SectorY: sy, X: lx, Y: ly}
}
