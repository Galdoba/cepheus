package convert

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates/cube"
	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates/global"
	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates/local"
)

const (
	sectorWidth           = 32
	sectorHeight          = 40
	loc2Glo_Offset_Width  = -1
	loc2Glo_Offset_Height = -40
)

// type Combined struct {
// 	Q, R, S          int
// 	SectorX, SectorY int
// 	LocalX, LocalY   int
// 	GlobalX, GlobalY int
// }

func Local_to_global(sx, sy, lx, ly int) (int, int) {
	if lx > sectorWidth || lx < 1 {
		panic(fmt.Sprintf("lx (%v) not in range 1-32", lx))
	}
	if ly > sectorHeight || ly < 1 {
		panic(fmt.Sprintf("ly (%v) not in range 1-40", ly))
	}
	x := lx + loc2Glo_Offset_Width + (sx * sectorWidth)
	y := ly + loc2Glo_Offset_Height + (sy * sectorHeight)
	return x, y
}

func LocalToGlobal(local local.SpaceSectorLocal) global.SpaceGlobal {
	sx, sy, lx, ly := local.SectorX, local.SectorY, local.X, local.Y
	if local.X > sectorWidth || local.X < 1 {
		panic(fmt.Sprintf("lx (%v) not in range 1-32", lx))
	}
	if ly > sectorHeight || ly < 1 {
		panic(fmt.Sprintf("ly (%v) not in range 1-40", ly))
	}
	x := lx + loc2Glo_Offset_Width + (sx * sectorWidth)
	y := ly + loc2Glo_Offset_Height + (sy * sectorHeight)
	return global.NewSpaceGlobal(x, y)
}

func Global_to_local(x, y int) (int, int, int, int) {
	//reverse offset
	adj_x := x - loc2Glo_Offset_Width
	adj_y := y - loc2Glo_Offset_Height

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
	return sx, sy, lx, ly
}

func GlobalToLocal(global global.SpaceGlobal) local.SpaceSectorLocal {
	//reverse offset
	adj_x := global.X - loc2Glo_Offset_Width
	adj_y := global.Y - loc2Glo_Offset_Height

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
	return local.NewLocal(sx, sy, lx, ly)
}

func Local_to_cube(sx, sy, lx, ly int) (int, int, int) {
	x := lx + loc2Glo_Offset_Width + (sx * sectorWidth)
	y := ly + loc2Glo_Offset_Height + (sy * sectorHeight)

	// q-axis
	q := x

	// r-axis
	s := -(y + (x / 2))

	// s-axis
	r := -q - s
	return q, r, s
}

func LocalToCube(local local.SpaceSectorLocal) cube.Cube {
	x := local.X + loc2Glo_Offset_Width + (local.SectorX * sectorWidth)
	y := local.Y + loc2Glo_Offset_Height + (local.SectorY * sectorHeight)

	// q-axis
	q := x

	// r-axis
	s := -(y + (x / 2))

	// s-axis
	r := -q - s
	return cube.NewCube(q, r, s)
}

func Cube_to_local(q, r, s int) (int, int, int, int) {
	x := q
	y := -((q / 2) + s)

	adj_x := x - loc2Glo_Offset_Width
	adj_y := y - loc2Glo_Offset_Height
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

	return sx, sy, lx, ly
}

func CubeToLocal(cube cube.Cube) local.SpaceSectorLocal {
	x := cube.Q
	y := -((cube.Q / 2) + cube.S)

	adj_x := x - loc2Glo_Offset_Width
	adj_y := y - loc2Glo_Offset_Height
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

	return local.NewLocal(sx, sy, lx, ly)
}

func Cube_to_global(q, r, s int) (int, int) {
	x := q
	y := ((q / 2) + s) * -1
	return x, y
}

func CubeToGlobal(cube cube.Cube) global.SpaceGlobal {
	x := cube.Q
	y := ((cube.Q / 2) + cube.S) * -1
	return global.NewSpaceGlobal(x, y)
}

func Global_to_cube(x, y int) (int, int, int) {
	q := x
	s := -(y + (x / 2))
	r := -q - s
	return q, r, s
}

func GlobalToCube(global global.SpaceGlobal) cube.Cube {
	q := global.X
	s := -(global.Y + (global.X / 2))
	r := -q - s
	return cube.NewCube(q, r, s)
}

func RoundTrip(q, r, s int) error {
	gx, gy := Cube_to_global(q, r, s)
	sx, sy, lx, ly := Global_to_local(gx, gy)
	rq, rr, rs := Local_to_cube(sx, sy, lx, ly)
	if q != rq || r != rr || s != rs {
		return fmt.Errorf("round trip failed: cube(%v, %v, %v)->global(%v, %v)->local(%v, %v, %v, %v)->return(%v, %v, %v)", q, r, s, gx, gy, sx, sy, lx, ly, rq, rr, rs)
	}
	sx, sy, lx, ly = Cube_to_local(q, r, s)
	gx, gy = Local_to_global(sx, sy, lx, ly)
	rq, rr, rs = Global_to_cube(gx, gy)
	if q != rq || r != rr || s != rs {
		return fmt.Errorf("round trip failed: cube(%v, %v, %v)->local(%v, %v, %v, %v)->global(%v, %v)->return(%v, %v, %v)", q, r, s, sx, sy, lx, ly, gx, gy, rq, rr, rs)
	}
	return nil
}
