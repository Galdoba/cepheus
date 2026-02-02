package tttable

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

// MockRoller для тестирования
type MockRoller struct {
	predefinedIndexes []int
	currentIndex      int
	predefinedCodes   []string
	currentCode       int
}

func newMockRoller(results ...int) *MockRoller {
	return &MockRoller{
		predefinedIndexes: results,
		currentIndex:      0,
	}
}

func (m *MockRoller) RollSafe(code string) (int, error) {
	if m.currentIndex >= len(m.predefinedIndexes) {
		return 0, fmt.Errorf("no more predefined results")
	}
	result := m.predefinedIndexes[m.currentIndex]
	m.currentIndex++
	return result, nil
}

func (m *MockRoller) ConcatRollSafe(code string) (string, error) {
	if m.currentIndex >= len(m.predefinedCodes) {
		return "", fmt.Errorf("no more predefined results")
	}
	result := m.predefinedCodes[m.currentCode]
	m.currentIndex++
	return result, nil
}

// TestMustIndex тестирует функцию MustIndex
func TestMustIndex(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected string
		panics   bool
	}{
		// Одно число
		{"Single positive", []int{5}, "5", false},
		{"Single negative", []int{-3}, "-3", false},
		{"Single zero", []int{0}, "0", false},

		// Диапазон в пределах границ
		{"Range within bounds", []int{1, 5}, "1-5", false},
		{"Range multiple numbers", []int{3, 1, 5, 2}, "1-5", false},
		{"Range with negatives", []int{-3, -1, -5}, "-5--1", false},
		{"Mixed range", []int{-2, 0, 2}, "-2-2", false},

		// Ошибки
		{"No arguments", []int{}, "", true},
		{"Same numbers", []int{5, 5, 5}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.panics {
					t.Errorf("MustIndex panicked unexpectedly: %v", r)
				} else if r == nil && tt.panics {
					t.Error("MustIndex should have panicked but didn't")
				}
			}()

			if !tt.panics {
				result := MustIndex(tt.numbers...)
				if result != tt.expected {
					t.Errorf("MustIndex(%v) = %v, want %v", tt.numbers, result, tt.expected)
				}
			} else {
				_ = MustIndex(tt.numbers...)
			}
		})
	}
}

// TestIndexSafe тестирует функцию IndexSafe
func TestIndexSafe(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected string
		wantErr  bool
	}{
		{"Valid single", []int{10}, "10", false},
		{"Valid range", []int{1, 10}, "1-10", false},
		{"Empty slice", []int{}, "", true},
		{"Same numbers", []int{3, 3}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IndexSafe(tt.numbers...)

			if (err != nil) != tt.wantErr {
				t.Errorf("IndexSafe(%v) error = %v, wantErr = %v", tt.numbers, err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("IndexSafe(%v) = %v, want %v", tt.numbers, result, tt.expected)
			}
		})
	}
}

// TestParseKey тестирует парсинг ключей
func TestParseKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		wantErr  bool
		expected *RangeKey
	}{
		// Одиночные числа
		{"Single positive", "5", false, &RangeKey{Min: 5, Max: 5, MinInclusive: true, MaxInclusive: true}},
		{"Single negative", "-3", false, &RangeKey{Min: -3, Max: -3, MinInclusive: true, MaxInclusive: true}},
		{"Single zero", "0", false, &RangeKey{Min: 0, Max: 0, MinInclusive: true, MaxInclusive: true}},

		// Диапазоны
		{"Positive range", "1-5", false, &RangeKey{Min: 1, Max: 5, MinInclusive: true, MaxInclusive: true}},
		{"Negative range", "-5--1", false, &RangeKey{Min: -5, Max: -1, MinInclusive: true, MaxInclusive: true}},
		{"Mixed range", "-2-3", false, &RangeKey{Min: -2, Max: 3, MinInclusive: true, MaxInclusive: true}},

		// Нижние границы
		{"Lower bound positive", "5-", false, &RangeKey{Min: MinRollBound, Max: 5, MinInclusive: true, MaxInclusive: true}},
		{"Lower bound negative", "-5-", false, &RangeKey{Min: MinRollBound, Max: -5, MinInclusive: true, MaxInclusive: true}},

		// Верхние границы
		{"Upper bound positive", "5+", false, &RangeKey{Min: 5, Max: MaxRollBound, MinInclusive: true, MaxInclusive: true}},
		{"Upper bound negative", "-5+", false, &RangeKey{Min: -5, Max: MaxRollBound, MinInclusive: true, MaxInclusive: true}},

		// Неверные форматы
		{"Empty string", "", true, nil},
		{"Invalid format", "abc", true, nil},
		{"Double minus", "--5", true, nil},
		{"Double plus", "++5", true, nil},
		{"Range reversed", "10-5", true, nil},
		{"Trailing spaces", "5 -", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseKey(tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseKey(%q) error = %v, wantErr = %v", tt.key, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Min != tt.expected.Min || result.Max != tt.expected.Max {
					t.Errorf("ParseKey(%q) = %+v, want %+v", tt.key, result, tt.expected)
				}
			}
		})
	}
}

// TestNewTable тестирует создание таблицы
func TestNewTable(t *testing.T) {

	tests := []struct {
		name      string
		tableName string
		opts      []TableOption
		wantErr   bool
	}{
		{
			name:      "Valid table with rows",
			tableName: "TestTable",
			opts: []TableOption{
				WithIndexEntries(
					NewTableEntry("1-5", "Event A"),
					NewTableEntry("6-10", "Event B"),
				),
			},
			wantErr: false,
		},
		{
			name:      "Table with roller",
			tableName: "TestTable2",
			opts: []TableOption{
				WithIndexEntries(NewTableEntry("1-10", "Event")),
			},
			wantErr: false,
		},
		{
			name:      "Empty table should fail",
			tableName: "EmptyTable",
			opts:      []TableOption{},
			wantErr:   true,
		},
		{
			name:      "Table with invalid key",
			tableName: "InvalidTable",
			opts: []TableOption{
				WithIndexEntries(NewTableEntry("invalid", "Event")),
			},
			wantErr: true,
		},
		{
			name:      "Table with overlapping ranges",
			tableName: "OverlapTable",
			opts: []TableOption{
				WithIndexEntries(NewTableEntry("1-5", "Event A"), NewTableEntry("3-7", "Event B")), // Пересекается с 1-5
			},
			wantErr: true,
		},
		{
			name:      "Table with duplicate keys",
			tableName: "DuplicateTable",
			opts: []TableOption{
				WithIndexEntries(NewTableEntry("1-5", "Event A"), NewTableEntry("1-5", "Event B")), // Дубликат
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := NewTable(tt.tableName, tt.opts...)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTable() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && table == nil {
				t.Error("NewTable() returned nil table without error")
			}

			if !tt.wantErr && table.Name != tt.tableName {
				t.Errorf("Table name = %v, want %v", table.Name, tt.tableName)
			}
		})
	}
}

// // TestTableAddRow тестирует добавление строк в таблицу
// func TestTableAddRow(t *testing.T) {
// 	table, err := NewTable("TestTable")
// 	if err == nil || table != nil {
// 		t.Fatal("Empty table should have failed validation")
// 	}

// 	table = &Table{
// 		Name:   "TestTable",
// 		Rows:   make(map[string]TableEntry),
// 		parsed: make(map[string]*RangeKey),
// 	}

// 	// Добавление валидной строки
// 	err = table.AddRow("1-5", "Event A")
// 	if err != nil {
// 		t.Errorf("AddRow() error = %v, want nil", err)
// 	}

// 	// Добавление невалидной строки
// 	err = table.AddRow("invalid", "Event B")
// 	if err == nil {
// 		t.Error("AddRow() should have returned error for invalid key")
// 	}

// 	// Добавление дубликата
// 	err = table.AddRow("1-5", "Event C")
// 	if err == nil {
// 		t.Error("AddRow() should have returned error for duplicate key")
// 	}

// 	// Проверка количества строк
// 	if len(table.Rows) != 1 {
// 		t.Errorf("Expected 1 row, got %d", len(table.Rows))
// 	}
// }

// TestFindByRoll тестирует поиск по результату броска
func TestFindByRoll(t *testing.T) {
	table, err := NewTable("TestTable",
		WithIndexEntries(
			NewTableEntry("-10-", "Very Low"),
			NewTableEntry("-9--1", "Low"),
			NewTableEntry("0", "Zero"),
			NewTableEntry("1-5", "Low Medium"),
			NewTableEntry("6-10", "High Medium"),
			NewTableEntry("11+", "High"),
		),
	)

	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	tests := []struct {
		roll     int
		expected string
		wantErr  bool
	}{
		{-50, "Very Low", false},  // Минимальная граница
		{-15, "Very Low", false},  // Нижняя граница
		{-5, "Low", false},        // Среднее отрицательное
		{0, "Zero", false},        // Ноль
		{3, "Low Medium", false},  // Среднее положительное
		{8, "High Medium", false}, // Высокое положительное
		{50, "High", false},       // Максимальная граница
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Roll%d", tt.roll), func(t *testing.T) {
			result, err := table.FindByRoll(tt.roll)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByRoll(%d) error = %v, wantErr = %v", tt.roll, err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("FindByRoll(%d) = %v, want %v", tt.roll, result, tt.expected)
			}
		})
	}
}

// TestTableValidation тестирует валидацию таблицы
func TestTableValidation(t *testing.T) {
	t1, _ := NewTable("t1", WithDiceExpression("2d6"), WithIndexEntries(NewTableEntry("1-5", "A"), NewTableEntry("6-10", "B"), NewTableEntry("11-15", "C")))
	// t2, _ := NewTable("t2", WithDiceExpression("2d6"), WithIndexEntries(NewTableEntry("1-5", "A"), NewTableEntry("6-10", "B"), NewTableEntry("11-15", "C")))
	t3, _ := NewTable("t3", WithDiceExpression("2d6"), WithIndexEntries(NewTableEntry("1-10", "A"), NewTableEntry("11-12", "B")))
	// t4, err := NewTable("t4", WithDiceExpression("2d6"), WithIndexEntries(NewTableEntry("1-", "A"), NewTableEntry("-6-10", "B")))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	tests := []struct {
		name    string
		table   Table
		wantErr bool
	}{
		{
			name:    "Valid non-overlapping ranges",
			table:   *t1,
			wantErr: false,
		},
		// {
		// 	name:    "Overlapping ranges",
		// 	table:   *t2,
		// 	wantErr: true,
		// },
		{
			name:    "Touching ranges are allowed",
			table:   *t3,
			wantErr: false,
		},
		// {
		// 	name:    "Mixed bounds",
		// 	table:   *t4,
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.table.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// TestNewTableCollection тестирует создание коллекции таблиц
func TestNewTableCollection(t *testing.T) {
	roller := dice.New("42")

	table1, _ := NewTable("Table1", WithIndexEntries(NewTableEntry("1-10", "Result1")))
	table2, _ := NewTable("Table2", WithIndexEntries(NewTableEntry("1-10", "Result2")))

	tests := []struct {
		name    string
		opts    []CollectionOption
		wantErr bool
	}{
		{
			name: "Valid collection with roller",
			opts: []CollectionOption{
				WithRoller(roller),
				WithTables(table1),
			},
			wantErr: false,
		},
		{
			name: "Collection with multiple tables",
			opts: []CollectionOption{
				WithRoller(roller),
				WithTables(table1, table2),
			},
			wantErr: false,
		},
		{
			name:    "Empty collection",
			opts:    []CollectionOption{},
			wantErr: false,
		},
		{
			name: "Duplicate table names",
			opts: []CollectionOption{
				WithTables(table1),
				WithTables(table1), // Дубликат
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection, err := NewTableCollection(tt.opts...)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewTableCollection() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && collection == nil {
				t.Error("NewTableCollection() returned nil without error")
			}
		})
	}
}

// TestCollectionRoll тестирует каскадные броски в коллекции
func TestCollectionRoll(t *testing.T) {
	roller := newMockRoller(1, 1, 11, 1, 1, 11)
	table1, _ := NewTable("Table1",
		WithDiceExpression("2d6"),
		// WithRoller(NewMockRoller(1, 1, 12)),
		WithIndexEntries(
			NewTableEntry("1-5", "Table2"),
			NewTableEntry("6-20", "DirectResult"),
		),
	)

	table2, _ := NewTable("Table2",
		WithDiceExpression("2d6"),
		WithIndexEntries(
			NewTableEntry("1-5", "Table3"),
			NewTableEntry("6-10", "Table3"),
			NewTableEntry("11-20", "DifferentResult"),
		),
	)

	table3, _ := NewTable("Table3",
		WithDiceExpression("2d6"),
		WithIndexEntries(
			NewTableEntry("1-10", "AnotherResult"),
			NewTableEntry("11+", "FinalResult"),
		),
	)

	collection, err := NewTableCollection(
		WithRoller(roller),
		WithTables(table1, table2, table3),
	)

	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}
	// Тест простого каскадного броска
	result, err := collection.Roll("Table1")
	if err != nil {
		t.Errorf("Roll() error = %v", err)
	}

	if result != "FinalResult" {
		t.Errorf("Roll() = %v, want FinalResult", result)
	}

	// Тест каскадного броска с полной последовательностью

	_, results, err := collection.RollCascade("Table1")
	if err != nil {
		t.Errorf("RollCascade() error = %v", err)
	}

	expectedSequence := []string{"Table2", "Table3", "FinalResult"}
	if len(results) != len(expectedSequence) {
		t.Errorf("RollCascade() returned %d results, want %d", len(results), len(expectedSequence))
	}

	for i, expected := range expectedSequence {
		if i < len(results) && results[i] != expected {
			t.Errorf("RollCascade()[%d] = %v, want %v", i, results[i], expected)
		}
	}

	// Тест с отсутствующей таблицей
	_, err = collection.Roll("NonExistentTable")
	if err == nil {
		t.Error("Roll() should have returned error for non-existent table")
	}

	// Тест без роллера
	collectionNoRoller, _ := NewTableCollection(WithTables(table1))
	_, err = collectionNoRoller.Roll("Table1")
	if err == nil {
		t.Error("Roll() should have returned error when roller is not set")
	}
}

// TestCollectionCycleDetection тестирует обнаружение циклов
func TestCollectionCycleDetection(t *testing.T) {
	roller := newMockRoller(1, 1, 1, 1, 1, 1, 1, 1, 1)

	tableA, _ := NewTable("TableA",
		WithDiceExpression("2d6"),
		WithIndexEntries(
			NewTableEntry("1-10", "TableB"),
			NewTableEntry("11-20", "ResultA"),
		),
	)

	tableB, _ := NewTable("TableB",
		WithDiceExpression("2d6"),
		WithIndexEntries(
			NewTableEntry("1-10", "TableA"), // Циклическая ссылка
			NewTableEntry("11-20", "ResultB"),
		),
	)

	collection, err := NewTableCollection(
		WithRoller(roller),
		WithTables(tableA, tableB),
	)

	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Должно обнаружить цикл
	_, err = collection.Roll("TableA")
	if err == nil {
		t.Error("Roll() should have detected cycle")
	}

	// Проверяем конкретный тип ошибки
	if err != nil && !contains(err.Error(), "cycle") {
		t.Errorf("Expected cycle error, got: %v", err)
	}
}

// TestRandomRolls тестирует случайные броски
func TestRandomRolls(t *testing.T) {
	// Используем настоящий случайный генератор
	realRoller := dice.New("42")

	table, err := NewTable("RandomTable",
		WithDiceExpression("2d6"),
		WithIndexEntries(
			NewTableEntry("1-5", "Low"),
			NewTableEntry("6-15", "Medium"),
			NewTableEntry("16-20", "High"),
		),
	)

	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Выполняем несколько бросков, проверяем что они не падают
	for i := 0; i < 100; i++ {
		_, result, err := table.roll(realRoller)
		if err != nil {
			t.Errorf("Roll() failed on iteration %d: %v", i, err)
			break
		}

		// Проверяем что результат один из ожидаемых
		validResults := []string{"Low", "Medium", "High"}
		valid := false
		for _, validResult := range validResults {
			if result == validResult {
				valid = true
				break
			}
		}

		if !valid {
			t.Errorf("Invalid result on iteration %d: %v", i, result)
		}
	}
}

// Вспомогательная структура для реальных случайных бросков
type randRoller struct {
	rng *rand.Rand
}

func (r *randRoller) Roll(code string) (int, error) {
	// Простая реализация для 1d20
	if code == "1d20" {
		return r.rng.Intn(20) + 1, nil
	}
	return 0, fmt.Errorf("unsupported dice code: %s", code)
}

// TestEdgeCases тестирует граничные случаи
func TestEdgeCases(t *testing.T) {
	// Таблица с единственной строкой
	table, err := NewTable("SingleRowTable",
		WithIndexEntries(NewTableEntry("1-", "Always")),
	)

	if err != nil {
		t.Fatalf("Failed to create single row table: %v", err)
	}

	// Проверяем что таблица валидна
	if err := table.Validate(); err != nil {
		t.Errorf("Single row table should be valid: %v", err)
	}

	// Таблица с полным покрытием
	table2, err := NewTable("FullCoverageTable",
		WithIndexEntries(
			NewTableEntry("-10-", "Low"),
			NewTableEntry("-9-10", "Middle"),
			NewTableEntry("11+", "High"),
		),
	)

	if err != nil {
		t.Fatalf("Failed to create full coverage table: %v", err)
	}

	// Проверяем несколько значений
	tests := []struct {
		roll     int
		expected string
	}{
		{-50, "Low"},   // Минимум
		{-10, "Low"},   // Граница
		{-5, "Middle"}, // Середина
		{0, "Middle"},  // Ноль
		{5, "Middle"},  // Середина
		{10, "Middle"}, // Граница
		{11, "High"},   // Граница
		{50, "High"},   // Максимум
	}

	for _, tt := range tests {
		result, err := table2.FindByRoll(tt.roll)
		if err != nil {
			t.Errorf("FindByRoll(%d) failed: %v", tt.roll, err)
			continue
		}

		if result != tt.expected {
			t.Errorf("FindByRoll(%d) = %v, want %v", tt.roll, result, tt.expected)
		}
	}
}

// TestGetMethods тестирует getter-методы
func TestGetMethods(t *testing.T) {
	table, err := NewTable("TestTable",
		WithIndexEntries(
			NewTableEntry("1-5", "Event A"),
			NewTableEntry("6-10", "Event B"),
		),
	)

	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// GetEvents
	events := table.GetAll()
	if len(events) != 2 {
		t.Errorf("GetEvents() should return 2 events, got %d", len(events))
	}

	// GetKeys
	keys := table.GetKeys()
	if len(keys) != 2 {
		t.Errorf("GetKeys() should return 2 keys, got %d", len(keys))
	}

	// Проверяем что ключи присутствуют
	expectedKeys := map[string]bool{"1-5": true, "6-10": true}
	for _, key := range keys {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key: %v", key)
		}
	}

	// Test collection GetTableNames
	collection, err := NewTableCollection(WithTables(table))
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	names := collection.GetTableNames()
	if len(names) != 1 || names[0] != "TestTable" {
		t.Errorf("GetTableNames() = %v, want [TestTable]", names)
	}
}

// Вспомогательная функция для проверки строк
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || contains(s[1:], substr)))
}
