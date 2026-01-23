package coordinates

import (
	"fmt"
	"strconv"
	"strings"
)

func (g Global) ToCube() Cube {
	q := g.x
	s := -(g.y + (g.x / 2))
	r := -q - s
	return MustCube(q, r, s)
}

func (g Global) ToLocal() Local {
	//reverse offset
	adj_x := g.x - loc2Glo_Offset_Width
	adj_y := g.y - loc2Glo_Offset_Height

	//x-axis
	temp := adj_x - 1
	sx := temp / sectorWidth
	rx := temp % sectorWidth //negative remainder
	//adjust negative remainder
	if rx < 0 {
		rx += sectorWidth
		sx--
	}
	lx := rx + 1

	//y-axis
	temp = adj_y - 1
	sy := temp / sectorHeight
	ry := temp % sectorHeight
	if ry < 0 {
		ry += sectorHeight
		sy--
	}
	ly := ry + 1
	return MustLocal(sx, sy, lx, ly)
}

func (g Global) DatabaseKey() string {
	return fmt.Sprintf("{%v,%v}", g.x, g.y)
}

func GlobalFromDatabaseKey(key string) (Global, error) {
	keyCopy := key
	keyParts := strings.Split(keyCopy, ",")
	if len(keyParts) != 2 {
		return Global{}, fmt.Errorf("bad key provided: '%v'", key)
	}
	x, err := strconv.Atoi(strings.TrimPrefix(keyParts[0], "{"))
	if err != nil {
		return Global{}, fmt.Errorf("failed to get X coordinate from '%v': %v", key, err)
	}
	y, err := strconv.Atoi(strings.TrimSuffix(keyParts[1], "}"))
	if err != nil {
		return Global{}, fmt.Errorf("failed to get Y coordinate from '%v': %v", key, err)
	}
	return NewGlobal(x, y), nil

}

func (g Global) OutOfReach() bool {
	if g.x > 10000 || g.x < -10000 {
		return true
	}
	if g.y > 10000 || g.y < -10000 {
		return true
	}
	return false
}
