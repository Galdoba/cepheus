package dice

/*
Roll - is a prime method of rolling dice: roll dice pool, modify and return sum of values as integer
ConcatRoll - is a support method of rolling dice: roll dice pool, modify and return concatenaded values as a string (code). this roll will return values consisting only with 0-9

*/

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"time"
)

type Roller struct {
	seed            string
	rng             *rand.Rand
	lastResult      []int
	lastModifiedSum int
	lastConcatRoll  string
}

func New(seed string) *Roller {
	return newRoller(seed)
}

func newRoller(seed string) *Roller {
	seedInt := randomSeed()
	if seed != "" {
		seedInt = stringToInt64(seed)
	}
	r := Roller{}
	r.rng = rand.New(rand.NewSource(int64(seedInt)))
	return &r
}

func stringToInt64(seed string) int64 {
	if seed == "" {
		return 0
	}

	var h1, h2 int64 = 0, 0

	for i := 0; i < len(seed); i++ {
		c := int64(seed[i])

		// Первый полином с большим простым множителем
		h1 = h1*131 + c + int64(i)*31

		// Второй полином с другим множителем для лучшего перемешивания
		h2 = h2*257 + c - int64(i)*17

		// Перемешиваем хэши между собой
		h1 ^= h2
		h2 ^= h1
	}

	// Объединяем оба хэша
	return h1 ^ (h2 << 32) ^ int64(len(seed))
}
func randomSeed() int64 {
	return time.Now().UnixNano()
}

func (r *Roller) baseRoll(number, faces int) []int {
	r.lastResult = []int{}
	for i := 0; i < number; i++ {
		v := r.rng.Intn(faces) + 1
		r.lastResult = append(r.lastResult, v)
	}
	return r.lastResult
}

// rollConcat бросает кубики для конкатенации
func (r *Roller) rollConcat(faces int) int {
	if faces == 0 {
		return 0
	}
	return r.rng.Intn(faces) + 1
}

func (r *Roller) rollSafe(sd SumDirectives) (int, error) {
	// 1. Бросок кубиков
	res := r.baseRoll(sd.Num, sd.Faces)

	// 2. Применение модификаторов к каждому броску
	modifiedRolls := make([]int, len(res))
	for i, oldVal := range res {
		currentVal := oldVal

		// 2a. Reroll значений
		if len(sd.ReRoll) > 0 {
			reroledVal, err := reroll(r, sd.Faces, currentVal, sd.ReRoll)
			if err != nil {
				return 0, err
			}
			currentVal = reroledVal
		}

		// 2b. Замена значений
		if newVal, ok := sd.Replace[currentVal]; ok {
			currentVal = newVal
		}

		// 2c. Индивидуальные модификаторы
		if ind, ok := sd.SumMods[Individual]; ok {
			currentVal = currentVal + ind
		}

		modifiedRolls[i] = currentVal
	}

	// 3. Drop low и drop high
	// Создаем копию для сортировки
	sortedRolls := make([]int, len(modifiedRolls))
	copy(sortedRolls, modifiedRolls)

	if dl, ok := sd.SumMods[DropLow]; ok && dl > 0 {
		// Сортируем по возрастанию
		slices.Sort(sortedRolls)
		if dl >= len(sortedRolls) {
			return 0, fmt.Errorf("all dices dropped") // Отбрасываем все кубики
		}
		// Отбрасываем dl низких значений
		sortedRolls = sortedRolls[dl:]
	}

	if dh, ok := sd.SumMods[DropHigh]; ok && dh > 0 {
		if dh >= len(sortedRolls) {
			return 0, fmt.Errorf("add dices dropped") // Отбрасываем все кубики
		}
		// Отбрасываем dh высоких значений
		sortedRolls = sortedRolls[:len(sortedRolls)-dh]
	}

	// 4. Суммируем оставшиеся значения
	sum := 0
	r.lastResult = sortedRolls
	for _, v := range sortedRolls {
		sum += v
	}

	// 5. Применение модификаторов к сумме
	if mod, ok := sd.SumMods[Additive]; ok {
		sum += mod
	}
	if mod, ok := sd.SumMods[Multiplicative]; ok {
		sum = sum * mod
	}
	if mod, ok := sd.SumMods[Deletive]; ok {
		if mod != 0 {
			sum = sum / mod
		}
	}
	if mod, ok := sd.SumMods[SumMininum]; ok {
		sum = max(mod, sum)
	}
	if mod, ok := sd.SumMods[SumMaximum]; ok {
		sum = min(mod, sum)
	}

	r.lastModifiedSum = sum
	return sum, nil
}

func (r *Roller) concat(cd ConcatDirectives) string {
	result := ""
	for i, faces := range cd.Faces {
		// Бросок кубика
		rollValue := r.rollConcat(faces)

		// Применение модификатора
		if i < len(cd.Mods) {
			rollValue += cd.Mods[i]
		}

		// Ограничение значения от 0 до 9
		rollValue = max(0, rollValue)
		rollValue = min(9, rollValue)

		// Добавление к результату
		result += strconv.Itoa(rollValue)
	}
	r.lastConcatRoll = result

	return result
}

func reroll(r *Roller, faces, val int, excluded map[int]bool) (int, error) {
	count := 0
	done := false
	current := val
	for !done {
		// Бросаем кубик
		current = r.baseRoll(1, faces)[0]

		// Проверяем, нужно ли перебрасывать
		needReroll := false
		for ex := range excluded {
			if current == ex {
				needReroll = true
				count++
				if count > 1000 {
					return 0, fmt.Errorf("impossible roll 1d%v to exclude %v", faces, excluded)
				}
				break
			}
		}

		if !needReroll {
			done = true
		}
	}

	return current, nil
}

func RollSafe(expression string, seed ...string) (int, error) {
	var seedStr string
	if len(seed) > 0 {
		seedStr = seed[0]
	}

	rd, err := DiceExpression(expression).ParseRoll()
	if err != nil {
		return 0, err
	}

	r := newRoller(seedStr)
	return r.rollSafe(rd)
}

func ConcatRollSafe(expression string, seed ...string) (string, error) {
	var seedStr string
	if len(seed) > 0 {
		seedStr = seed[0]
	}

	rd, err := DiceExpression(expression).ParseConcatRoll()
	if err != nil {
		return "", err
	}

	r := newRoller(seedStr)
	return r.concat(rd), nil
}

func Roll(expression string, seed ...string) int {
	result, err := RollSafe(expression, seed...)
	if err != nil {
		panic(err)
	}
	return result
}

func ConcatRoll(expression string, seed ...string) string {
	result, err := ConcatRollSafe(expression, seed...)
	if err != nil {
		panic(err)
	}
	return result
}

func (r *Roller) Roll(expression string) int {
	rd, err := DiceExpression(expression).ParseRoll()
	if err != nil {
		panic(err)
	}
	s, err := r.rollSafe(rd)
	if err != nil {
		panic(err)
	}
	return s
}

func (r *Roller) ConcatRoll(expression string) string {
	rd, err := DiceExpression(expression).ParseConcatRoll()
	if err != nil {
		panic(err)
	}
	return r.concat(rd)
}

func (r *Roller) LastRoll() int {
	return r.lastModifiedSum
}

func (r *Roller) LastConcatRoll() string {
	return r.lastConcatRoll
}

func (r *Roller) Result() []int {
	return r.lastResult
}
