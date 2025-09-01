package cube

import "testing"

func TestNeighborsWithNewDirections(t *testing.T) {
	center := NewCube(0, 0, 0)
	neighbors := Neighbors(center)

	expected := []Cube{
		{Q: 0, R: -1, S: 1}, // север
		{Q: 1, R: -1, S: 0}, // северо-восток
		{Q: 1, R: 0, S: -1}, // юго-восток
		{Q: 0, R: 1, S: -1}, // юг
		{Q: -1, R: 1, S: 0}, // юго-запад
		{Q: -1, R: 0, S: 1}, // северо-запад
	}

	if len(neighbors) != 6 {
		t.Fatalf("Ожидалось 6 соседей, получил %d", len(neighbors))
	}

	for i, n := range neighbors {
		if n.Q != expected[i].Q || n.R != expected[i].R || n.S != expected[i].S {
			t.Errorf("Неверный сосед %d: ожидалось (%d,%d,%d), получил (%d,%d,%d)",
				i, expected[i].Q, expected[i].R, expected[i].S, n.Q, n.R, n.S)
		}
	}
}

func TestMovement(t *testing.T) {
	center := NewCube(0, 0, 0)

	tests := []struct {
		direction int
		expected  Cube
	}{
		{0, Cube{Q: 0, R: -1, S: 1}},
		{1, Cube{Q: 1, R: -1, S: 0}},
		{2, Cube{Q: 1, R: 0, S: -1}},
		{3, Cube{Q: 0, R: 1, S: -1}},
		{4, Cube{Q: -1, R: 1, S: 0}},
		{5, Cube{Q: -1, R: 0, S: 1}},
	}

	for _, test := range tests {
		result := neighbor(center, test.direction)
		if result != test.expected {
			t.Errorf("Направление %d: ожидалось (%d,%d,%d), получил (%d,%d,%d)",
				test.direction,
				test.expected.Q, test.expected.R, test.expected.S,
				result.Q, result.R, result.S)
		}
	}
}

func TestSpiralMaps(t *testing.T) {
	center := NewCube(0, 0, 0)
	indexToHex, hexToIndex := SpiralMaps(center, 2)

	// Проверка центрального гекса
	if indexToHex[0] != center {
		t.Error("Индекс 0 должен соответствовать центру")
	}
	if hexToIndex[center] != 0 {
		t.Error("Центр должен иметь индекс 0")
	}

	// Ожидаемый порядок для радиуса 1
	radius1Order := []Cube{
		neighbor(center, 0), // север
		neighbor(center, 1), // северо-восток
		neighbor(center, 2), // юго-восток
		neighbor(center, 3), // юг
		neighbor(center, 4), // юго-запад
		neighbor(center, 5), // северо-запад
	}

	// Проверка порядка для радиуса 1
	for i, expected := range radius1Order {
		idx := i + 1 // Индексы после центра
		if indexToHex[idx] != expected {
			t.Errorf("Позиция %d: ожидалось %v, получено %v",
				idx, expected, indexToHex[idx])
		}
		if hexToIndex[expected] != idx {
			t.Errorf("Обратная карта: ожидался индекс %d для %v",
				idx, expected)
		}
	}

	// Проверка количества элементов
	expectedCount := 1 + 6 + 12 // Центр + радиус1 + радиус2
	if len(indexToHex) != expectedCount {
		t.Errorf("Ожидалось %d элементов, получено %d",
			expectedCount, len(indexToHex))
	}
	if len(hexToIndex) != expectedCount {
		t.Errorf("Ожидалось %d элементов, получено %d",
			expectedCount, len(hexToIndex))
	}

	// Проверка уникальности индексов
	for i := 0; i < expectedCount; i++ {
		if _, exists := indexToHex[i]; !exists {
			t.Errorf("Индекс %d отсутствует в прямой карте", i)
		}
	}

	// Проверка расстояний
	for hex, index := range hexToIndex {
		dist := Distance(center, hex)
		expectedRadius := 0
		switch {
		case index == 0:
			expectedRadius = 0
		case index < 7:
			expectedRadius = 1
		default:
			expectedRadius = 2
		}

		if dist != expectedRadius {
			t.Errorf("Для индекса %d: расстояние %d, ожидалось %d",
				index, dist, expectedRadius)
		}
	}
}

func TestSpiralMapsEdgeCases(t *testing.T) {
	// Тест с радиусом 0
	center := NewCube(0, 0, 0)
	indexToHex, hexToIndex := SpiralMaps(center, 0)

	if len(indexToHex) != 1 || len(hexToIndex) != 1 {
		t.Error("Для радиуса 0 должна быть только центральная точка")
	}

	// Тест с отрицательным радиусом
	indexToHexNeg, hexToIndexNeg := SpiralMaps(center, -1)
	if len(indexToHexNeg) != 0 || len(hexToIndexNeg) != 0 {
		t.Error("Для отрицательного радиуса должны быть пустые карты")
	}

	// Тест с большим радиусом
	indexToHexLarge, _ := SpiralMaps(center, 3)
	expectedCount := 1 + 6 + 12 + 18 // Центр + 3 кольца
	if len(indexToHexLarge) != expectedCount {
		t.Errorf("Ожидалось %d элементов, получено %d",
			expectedCount, len(indexToHexLarge))
	}
}

func TestVector(t *testing.T) {
	tests := []struct {
		name     string
		start    Cube
		change   Cube
		expected Cube
	}{
		{
			name:     "Базовое сложение",
			start:    NewCube(1, 2, -3),
			change:   NewCube(2, -1, -1),
			expected: Cube{Q: 3, R: 1, S: -4},
		},
		{
			name:     "Нулевое изменение",
			start:    NewCube(5, -3, -2),
			change:   NewCube(0, 0, 0),
			expected: Cube{Q: 5, R: -3, S: -2},
		},
		{
			name:     "Отрицательные координаты",
			start:    NewCube(-3, 2, 1),
			change:   NewCube(4, -5, 1),
			expected: Cube{Q: 1, R: -3, S: 2},
		},
		{
			name:     "Крайние значения",
			start:    NewCube(1000000, -500000, -500000),
			change:   NewCube(-1000000, 1000000, 0),
			expected: Cube{Q: 0, R: 500000, S: -500000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Vector(tt.start, tt.change)

			if result.Q != tt.expected.Q ||
				result.R != tt.expected.R ||
				result.S != tt.expected.S {
				t.Errorf("Ожидалось (%d, %d, %d), получил (%d, %d, %d)",
					tt.expected.Q, tt.expected.R, tt.expected.S,
					result.Q, result.R, result.S)
			}

			// Проверка инварианта q+r+s=0
			sum := result.Q + result.R + result.S
			if sum != 0 {
				t.Errorf("Нарушен инвариант: q+r+s=%d, ожидалось 0", sum)
			}
		})
	}
}

func TestVectorWithDirections(t *testing.T) {
	center := NewCube(0, 0, 0)
	north := NewCube(0, -1, 1) // Направление 0: север

	// Проверка перемещения
	moved := Vector(center, north)
	expected := NewCube(0, -1, 1)

	if moved.Q != expected.Q || moved.R != expected.R || moved.S != expected.S {
		t.Errorf("Ожидалось (%d,%d,%d), получил (%d,%d,%d)",
			expected.Q, expected.R, expected.S,
			moved.Q, moved.R, moved.S)
	}

	// Двойное перемещение
	movedTwice := Vector(moved, north)
	expectedTwice := NewCube(0, -2, 2)

	if movedTwice.Q != expectedTwice.Q ||
		movedTwice.R != expectedTwice.R ||
		movedTwice.S != expectedTwice.S {
		t.Errorf("Ожидалось (%d,%d,%d), получил (%d,%d,%d)",
			expectedTwice.Q, expectedTwice.R, expectedTwice.S,
			movedTwice.Q, movedTwice.R, movedTwice.S)
	}
}

func TestRing(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		center Cube
		radius int
		want   []Cube
	}{
		{
			name:   "zero radius",
			center: Cube{0, 0, 0},
			radius: 0,
			want:   []Cube{{0, 0, 0}},
		}, // TODO: Add test cases.
		{
			name:   "negative radius",
			center: Cube{0, 0, 0},
			radius: -4,
			want:   []Cube{},
		}, // TODO: Add test cases.
		{
			name:   "radius 1",
			center: Cube{0, 0, 0},
			radius: 1,
			want:   []Cube{{0, -1, 1}, {1, -1, 0}, {1, 0, -1}, {0, 1, -1}, {-1, 1, 0}, {-1, 0, 1}},
		}, // TODO: Add test cases.
		{
			name:   "radius 2",
			center: Cube{0, 0, 0},
			radius: 2,
			want:   []Cube{{0, -2, 2}, {1, -2, 1}, {2, -2, 0}, {2, -1, -1}, {2, 0, -2}, {1, 1, -2}, {0, 2, -2}, {-1, 2, -1}, {-2, 2, 0}, {-2, 1, 1}, {-2, 0, 2}, {-1, -1, 2}},
		}, // TODO: Add test cases.
		{
			name:   "radius 3",
			center: Cube{0, 0, 0},
			radius: 3,
			want:   []Cube{{0, -3, 3}, {1, -3, 2}, {2, -3, 1}, {3, -3, 0}, {3, -2, -1}, {3, -1, -2}, {3, 0, -3}, {2, 1, -3}, {1, 2, -3}, {0, 3, -3}, {-1, 3, -2}, {-2, 3, -1}, {-3, 3, 0}, {-3, 2, 1}, {-3, 1, 2}, {-3, 0, 3}, {-2, -1, 3}, {-1, -2, 3}},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ring(tt.center, tt.radius)
			// TODO: update the condition below to compare got with tt.want.
			for i, hex := range tt.want {
				if hex != got[i] {
					t.Errorf("Ring() = %v, want %v", got, tt.want)
				}
			}
			if len(got) != len(tt.want) {
				t.Errorf("Ring size = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}
