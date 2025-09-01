package cube

import (
	"math"
)

const (
	north     = 0
	northEast = 1
	southEast = 2
	south     = 3
	southWest = 4
	northWest = 5
)

// Distance возвращает расстояние между двумя гексами в шагах
func Distance(a, b Cube) int {
	return (abs(a.Q-b.Q) + abs(a.R-b.R) + abs(a.S-b.S)) / 2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var directions = []Cube{
	{Q: 0, R: -1, S: 1}, // 0: north
	{Q: 1, R: -1, S: 0}, // 1: northEast
	{Q: 1, R: 0, S: -1}, // 2: southEast
	{Q: 0, R: 1, S: -1}, // 3: south
	{Q: -1, R: 1, S: 0}, // 4: southWest
	{Q: -1, R: 0, S: 1}, // 5: northWest
}

// Neighbors возвращает 6 соседних гексов в порядке направлений
func Neighbors(h Cube) []Cube {
	neighbors := make([]Cube, 0, 6)
	for _, d := range directions {
		neighbors = append(neighbors, Cube{
			Q: h.Q + d.Q,
			R: h.R + d.R,
			S: h.S + d.S,
		})
	}
	return neighbors
}

// LineDrawing возвращает гексы на прямой линии между двумя точками
func LineDrawing(a, b Cube) []Cube {
	N := Distance(a, b)
	results := make([]Cube, 0, N+1)

	// Добавляем начальную точку
	results = append(results, a)

	// Линейная интерполяция для N сегментов
	for i := 1; i <= N; i++ {
		t := float64(i) / float64(N)
		q := lerp(float64(a.Q), float64(b.Q), t)
		r := lerp(float64(a.R), float64(b.R), t)
		s := lerp(float64(a.S), float64(b.S), t)

		// Округляем до ближайшего гекса
		hex := roundCube(q, r, s)
		results = append(results, hex)
	}

	return results
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func roundCube(q, r, s float64) Cube {
	rq := math.Round(q)
	rr := math.Round(r)
	rs := math.Round(s)

	// Коррекция округления для соблюдения q + r + s = 0
	dq := math.Abs(rq - q)
	dr := math.Abs(rr - r)
	ds := math.Abs(rs - s)

	if dq > dr && dq > ds {
		rq = -rr - rs
	} else if dr > ds {
		rr = -rq - rs
	} else {
		rs = -rq - rr
	}

	return Cube{
		Q: int(rq),
		R: int(rr),
		S: int(rs),
	}
}

// Ring возвращает гексы на заданном расстоянии от центра
func Ring(center Cube, radius int) []Cube {
	if radius < 0 {
		return []Cube{}
	}
	if radius == 0 {
		return []Cube{center}
	}

	// Начинаем с направления north и двигаемся на radius шагов
	// после чего выбираем следующее направление
	ring := make([]Cube, 0, 6*radius)
	current := move(center, 0, radius)
	directionOptimization := 2
	for direction := range 6 {
		for range radius {
			ring = append(ring, current)
			current = neighbor(current, (direction+directionOptimization)%6)
		}
	}
	return ring
}

// Вспомогательная функция для перемещения на несколько шагов
func move(h Cube, direction, steps int) Cube {
	current := h
	for range steps {
		current = neighbor(current, direction)
	}
	return current
}

// Получение соседа в определенном направлении
func neighbor(h Cube, direction int) Cube {
	d := directions[direction]
	return Cube{
		Q: h.Q + d.Q,
		R: h.R + d.R,
		S: h.S + d.S,
	}
}

// Spiral возвращает все гексы в пределах заданного радиуса
func Spiral(center Cube, radius int) []Cube {
	results := []Cube{center}

	for k := 1; k <= radius; k++ {
		ring := Ring(center, k)
		results = append(results, ring...)
	}

	return results
}

// SpiralMaps возвращает две карты для спирального обхода
func SpiralMaps(center Cube, radius int) (map[int]Cube, map[Cube]int) {
	indexToCube := make(map[int]Cube)
	hexToIndex := make(map[Cube]int)
	if radius < 0 {
		return indexToCube, hexToIndex
	}

	// Центральный гекс
	counter := 0
	indexToCube[counter] = center
	hexToIndex[center] = counter
	counter++

	// Обходим кольца от 1 до заданного радиуса
	for r := 1; r <= radius; r++ {
		ringOfHexes := Ring(center, r)
		for _, hex := range ringOfHexes {
			indexToCube[counter] = hex
			hexToIndex[hex] = counter
			counter++
		}
	}

	return indexToCube, hexToIndex
}

// Vector выполняет покоординатное сложение двух гексов
func Vector(start, change Cube) Cube {
	return Cube{
		Q: start.Q + change.Q,
		R: start.R + change.R,
		S: start.S + change.S,
	}
}

func Rotate(start, center Cube, steps int) Cube {
	// Если start совпадает с центром или нулевые шаги
	if start == center || steps == 0 {
		return start
	}

	// Вычисляем радиус кольца
	r := Distance(center, start)

	// Получаем кольцо с нужным радиусом
	ring := Ring(center, r)

	// Находим позицию start в кольце
	pos := -1
	for i, hex := range ring {
		if hex == start {
			pos = i
			break
		}
	}

	// Если точка не найдена в кольце
	if pos == -1 {
		return start
	}

	// Вычисляем новую позицию с учетом направления
	n := len(ring)
	if steps > 0 {
		// Вращение по часовой стрелке
		return ring[(pos+steps)%n]
	} else {
		// Вращение против часовой стрелки
		newPos := (pos + steps) % n
		if newPos < 0 {
			newPos += n
		}
		return ring[newPos]
	}
}

/*
X12345678|
0  ____  |
1 /    \ |+--------+
2/      \||        |
3\      /||        +--------+
          |        |        |
4 \____/ |+--------+        +
          |        |        |
		  +        +--------+



*/
