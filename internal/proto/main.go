package main

import (
	"container/heap"
	"fmt"
	"math"

	. "github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
)

// Cube - кубические координаты гексагональной сетки на плоскости

// CubeMoveCost - стоимость перехода между двумя гексами
type CubeMoveCost struct {
	From, To Cube
	Cost     float64
}

// Node - узел для алгоритма A*
type Node struct {
	Cube
	G, H, F float64 // стоимости
	Parent  *Node
	Index   int // для кучи
}

// Направления в гексагональной сетке (шесть соседей)

// HexPathfinder - основной структура поиска пути в гексагональной сетке
type HexPathfinder struct {
	// Данные о стоимости
	EntryCosts map[Cube]float64         // цена входа в гекс
	ExitCosts  map[CubeMoveCost]float64 // цена выхода из гекса

	// Параметры
	MaxJump   int                     // максимальная длина прыжка (6)
	Heuristic func(a, b Cube) float64 // эвристическая функция

	// Вспомогательные структуры
	openSet   *PriorityQueue
	closedSet map[Cube]bool
	gScore    map[Cube]float64
	cameFrom  map[Cube]Cube

	// Границы поля (опционально)
	MinQ, MaxQ, MinR, MaxR, MinS, MaxS int
}

// PriorityQueue - минимальная куча для openSet
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	node := x.(*Node)
	node.Index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}

// func CubeNeighbor(c Cube, direction int) Cube {
// 	return Neighbors(c)[direction]
// }

func CubeDistance(a, b Cube) int {
	return Distance(a, b)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CubeHeuristic - эвристика для гексагональной сетки (расстояние между гексами)
func CubeHeuristic(a, b Cube) float64 {
	return float64(CubeDistance(a, b))
}

// // GetRing - возвращает все гексы на заданном расстоянии от центра
// func GetRing(center Cube, radius int) []Cube {
// 	if radius == 0 {
// 		return []Cube{center}
// 	}

// 	var results []Cube
// 	// Начинаем с гекса на расстоянии radius по направлению 4
// 	cube := CubeAdd(center, CubeScale(CubeDirections[4], radius))

// 	for i := 0; i < 6; i++ {
// 		for j := 0; j < radius; j++ {
// 			results = append(results, cube)
// 			cube = CubeNeighbor(cube, i)
// 		}
// 	}

// 	return results
// }

// GetAllInRadius - возвращает все гексы в пределах заданного радиуса
// func GetAllInRadius(center Cube, radius int) []Cube {
// 	var results []Cube
// 	for q := -radius; q <= radius; q++ {
// 		for r := max(-radius, -q-radius); r <= min(radius, -q+radius); r++ {
// 			s := -q - r
// 			if abs(s) <= radius {
// 				results = append(results, CubeAdd(center, Cube{Q: q, R: r, S: s}))
// 			}
// 		}
// 	}
// 	return results
// }

// getNeighbors - возвращает всех соседей в радиусе прыжка (исключая сам гекс)
func (pf *HexPathfinder) getNeighbors(c Cube) []Cube {
	var neighbors []Cube

	// Генерируем все гексы в радиусе до MaxJump
	// Исключаем сам гекс (радиус 0)
	for radius := 1; radius <= pf.MaxJump; radius++ {
		// Получаем все гексы на заданном расстоянии
		ring := Ring(c, radius)

		for _, neighbor := range ring {
			// Проверяем границы поля, если заданы
			if pf.MinQ != pf.MaxQ {
				if neighbor.Q() < pf.MinQ || neighbor.Q() > pf.MaxQ ||
					neighbor.R() < pf.MinR || neighbor.R() > pf.MaxR ||
					neighbor.S() < pf.MinS || neighbor.S() > pf.MaxS {
					continue
				}
			}

			// Проверяем, существует ли переход (стоимость выхода задана)
			if _, exists := pf.ExitCosts[CubeMoveCost{From: c, To: neighbor}]; !exists {
				// Если стоимость перехода не задана, пропускаем
				// Можно также установить стоимость по умолчанию
				pf.ExitCosts[CubeMoveCost{
					From: c,
					To:   neighbor,
					Cost: 1,
				}] = float64(7 - Distance(c, neighbor))
				continue
			}

			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

// IsInBounds - проверяет, находится ли гекс в пределах поля
func (pf *HexPathfinder) IsInBounds(c Cube) bool {
	if pf.MinQ == pf.MaxQ {
		return true // Границы не заданы
	}

	return c.Q() >= pf.MinQ && c.Q() <= pf.MaxQ &&
		c.R() >= pf.MinR && c.R() <= pf.MaxR &&
		c.S() >= pf.MinS && c.S() <= pf.MaxS
}

// SetBounds - устанавливает границы поля
func (pf *HexPathfinder) SetBounds(minQ, maxQ, minR, maxR int) {
	pf.MinQ = minQ
	pf.MaxQ = maxQ
	pf.MinR = minR
	pf.MaxR = maxR
	// S вычисляем из условия Q+R+S=0
	pf.MinS = -maxQ - maxR
	pf.MaxS = -minQ - minR
}

// FindPath - основной алгоритм поиска пути A*
func (pf *HexPathfinder) FindPath(start, goal Cube) ([]Cube, float64, bool) {
	// Проверяем корректность координат
	// if !start.IsValid() || !goal.IsValid() {
	// 	return nil, 0, false
	// }

	// Инициализация
	pf.openSet = &PriorityQueue{}
	pf.closedSet = make(map[Cube]bool)
	pf.gScore = make(map[Cube]float64)
	pf.cameFrom = make(map[Cube]Cube)

	// Начальный узел
	startNode := &Node{
		Cube:   start,
		G:      0,
		H:      pf.Heuristic(start, goal),
		F:      pf.Heuristic(start, goal),
		Parent: nil,
	}
	pf.gScore[start] = 0

	heap.Push(pf.openSet, startNode)

	for pf.openSet.Len() > 0 {
		// Извлекаем узел с минимальной f-стоимостью
		current := heap.Pop(pf.openSet).(*Node)

		// Если достигли цели
		if current.Cube == goal {
			return pf.reconstructPath(current), current.G, true
		}

		pf.closedSet[current.Cube] = true

		// Обрабатываем соседей (все гексы в радиусе прыжка)
		neighbors := pf.getNeighbors(current.Cube)

		for _, neighbor := range neighbors {
			// Пропускаем уже обработанные
			if pf.closedSet[neighbor] {
				continue
			}

			// Получаем стоимость перехода
			moveCost, moveExists := pf.ExitCosts[CubeMoveCost{From: current.Cube, To: neighbor}]
			if !moveExists {
				// Если стоимость перехода не задана, переход невозможен
				continue
			}

			// Получаем стоимость входа в соседний гекс
			entryCost, entryExists := pf.EntryCosts[neighbor]
			if !entryExists {
				// Если стоимость входа не задана, считаем бесконечно большой
				entryCost = math.Inf(1)
			}

			// Вычисляем новую стоимость пути
			tentativeG := current.G + moveCost + entryCost

			// Получаем текущую стоимость соседа
			currentG, exists := pf.gScore[neighbor]
			if !exists {
				currentG = math.Inf(1)
			}

			// Если нашли лучший путь
			if tentativeG < currentG {
				pf.cameFrom[neighbor] = current.Cube
				pf.gScore[neighbor] = tentativeG

				h := pf.Heuristic(neighbor, goal)
				f := tentativeG + h

				neighborNode := &Node{
					Cube:   neighbor,
					G:      tentativeG,
					H:      h,
					F:      f,
					Parent: current,
				}

				// Добавляем или обновляем в openSet
				heap.Push(pf.openSet, neighborNode)
			}
		}
	}

	return nil, 0, false // Путь не найден
}

// reconstructPath - восстанавливает путь от конечной точки до старта
func (pf *HexPathfinder) reconstructPath(end *Node) []Cube {
	var path []Cube
	current := end

	for current != nil {
		path = append([]Cube{current.Cube}, path...)
		current = current.Parent
	}

	return path
}

// NewHexPathfinder - создает новый экземпляр поисковика пути
func NewHexPathfinder(maxJump int, heuristic func(a, b Cube) float64) *HexPathfinder {
	if heuristic == nil {
		heuristic = CubeHeuristic
	}

	return &HexPathfinder{
		EntryCosts: make(map[Cube]float64),
		ExitCosts:  make(map[CubeMoveCost]float64),
		MaxJump:    maxJump,
		Heuristic:  heuristic,
	}
}

// Пример использования
func main() {
	// Создаем поисковик с максимальным прыжком 6
	pf := NewHexPathfinder(6, CubeHeuristic)

	// Заполнение стоимостей входа (пример)
	for q := -21; q <= 21; q++ {
		for r := -21; r <= 21; r++ {
			s := -q - r
			if abs(s) <= 21 {
				cube := MustCube(q, r, s)
				// Базовая стоимость входа
				pf.EntryCosts[cube] = 1.0
			}
		}
	}

	// Устанавливаем стоимости выходов (пример)
	// Здесь можно задать сложную логику с нетранзитивными переходами
	for q := -21; q <= 21; q++ {
		for r := -21; r <= 21; r++ {
			s := -q - r
			if abs(s) <= 21 {
				from := MustCube(q, r, s)
				// Для каждого гекса задаем стоимости прыжков в соседние гексы
				for radius := 1; radius <= 6; radius++ {
					ring := Ring(from, radius)
					for _, to := range ring {
						if pf.IsInBounds(to) {
							// Пример: стоимость зависит от расстояния, но можно сделать сложнее
							distance := CubeDistance(from, to)
							cost := float64(distance) * 1.5
							pf.ExitCosts[CubeMoveCost{From: from, To: to}] = cost
						}
					}
				}
			}
		}
	}

	// // Пример неcccранзитивного перехода: A->C != A->B + B->C
	// // Можно задать явно:
	// pf.ExitCosts[CubeMoveCost{
	// 	From: Cube{Q: 0, R: 0, S: 0},
	// 	To:   Cube{Q: 3, R: -1, S: -2},
	// }] = 2.0 // Дешевле, чем через промежуточные гексы

	// Поиск пути
	start := MustCube(0, 0, 0)
	goal := MustCube(21, -14, -7)

	path, cost, found := pf.FindPath(start, goal)

	if found {
		fmt.Printf("Путь найден! Стоимость: %.2f\n", cost)
		fmt.Printf("Длина пути: %d гексов\n", len(path)-1)
		for i, cube := range path {
			fmt.Printf("Шаг %d: (%d, %d, %d)\n", i, cube.Q(), cube.R(), cube.S())
		}
	} else {
		fmt.Println("Путь не найден")
	}
}
