package coordinates

import (
	"fmt"
)

const (
	sectorWidth           = 32
	sectorHeight          = 40
	loc2Glo_Offset_Width  = -1
	loc2Glo_Offset_Height = -40
)

func (l Local) ToGlobal() Global {
	x := l.lx + loc2Glo_Offset_Width + (l.sx * sectorWidth)
	y := l.ly + loc2Glo_Offset_Height + (l.sy * sectorHeight)
	return NewGlobal(x, y)
}

func (l Local) ToCube() Cube {
	return l.ToGlobal().ToCube()
}

func (l Local) Hex() string {
	return fmt.Sprintf("%02d%02d", l.lx, l.ly)
}
