package gametime

import (
	"fmt"
	"sync"
)

const (
	ImperialCalendar     Calendar = "IC"
	secondsPerTick                = 6
	minutesPerHour                = 60
	minutesPerSpaceRound          = 6
	hoursPerDay                   = 24
	daysPerWeek                   = 7
	daysPerMonth                  = 28
	daysPerYear                   = 365
	ticsPerMinute                 = 10
	ticsPerSpaceRound             = ticsPerMinute * minutesPerSpaceRound
	ticsPerHour                   = minutesPerHour * ticsPerMinute
	ticsPer4Hours                 = 4 * ticsPerHour
	ticsPer10Hours                = 10 * ticsPerHour
	ticsPerDay                    = hoursPerDay * ticsPerHour
	ticsPerWeek                   = daysPerWeek * ticsPerDay
	ticsPerMonth                  = daysPerMonth * ticsPerDay
	ticsPerYear                   = daysPerYear * ticsPerDay
	Holiday                       = "Holiday"
	Wonday                        = "Wonday"
	Tuday                         = "Tuday"
	Thirday                       = "Thirday"
	Forday                        = "Forday"
	Fiday                         = "Fiday"
	Sixday                        = "Sixday"
	Senday                        = "Senday"
)

type Calendar string

var calendar = ImperialCalendar

type GameTime int64

type Global struct {
	time GameTime
	mu   sync.RWMutex
}

func (g *Global) NextCombatRound() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.time++
}

func (g *Global) NextSpaceRound() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.time += ticsPerMinute * 6
}

func (g *Global) NextPortRound() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.time += ticsPerHour
}

func (g *Global) Time() GameTime {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.time
}

func (gt GameTime) Seconds() int {
	return int(gt%10) * secondsPerTick
}

func (gt GameTime) Minutes() int {
	return int((gt / 10) % minutesPerHour)
}

func (gt GameTime) Hour() int {
	ticsInDay := int64(gt) % ticsPerDay
	return int(ticsInDay / ticsPerHour)
}

func (gt GameTime) Day() int {
	totalDays := int64(gt) / ticsPerDay
	return int(totalDays%daysPerYear) + 1
}

func (gt GameTime) DayOfMonth() int {
	dayOfYear := gt.Day()
	if dayOfYear == 1 {
		return 0
	}

	dayInMonth := (dayOfYear - 2) % daysPerMonth
	return dayInMonth + 1
}

func (gt GameTime) Month() int {
	dayOfYear := gt.Day()

	if dayOfYear == 1 {
		return 0
	}
	month := (dayOfYear - 2) / daysPerMonth
	return month + 1
}

func (gt GameTime) DayOfWeek() int {
	dayOfYear := gt.Day()

	if dayOfYear == 1 {
		return 0
	}
	return ((dayOfYear - 2) % daysPerWeek) + 1
}

func (gt GameTime) Weekday() string {
	switch gt.DayOfWeek() {
	case 0:
		return Holiday
	case 1:
		return Wonday
	case 2:
		return Tuday
	case 3:
		return Thirday
	case 4:
		return Forday
	case 5:
		return Fiday
	case 6:
		return Sixday
	case 7:
		return Senday
	}
	return "Undefined"
}

func (gt GameTime) WeekOfMonth() int {
	d := gt.DayOfMonth()
	switch d {
	case -1, 0:
		return 0
	default:
		return ((d - 1) / daysPerWeek) + 1
	}
}

func (gt GameTime) Year() int {
	totalYears := int64(gt) / ticsPerYear
	return int(totalYears) + 1 // Года начинаются с 1
}

func (gt GameTime) IsHoliday() bool {
	return gt.Day() == 1
}

func (gt GameTime) AddTicks(ticks GameTime) GameTime {
	return gt + ticks
}

func (gt GameTime) AddMinutes(minutes int) GameTime {
	return gt + GameTime(minutes*ticsPerMinute)
}

func (gt GameTime) AddHours(hours int) GameTime {
	return gt + GameTime(hours*ticsPerHour)
}

func (gt GameTime) AddDays(days int) GameTime {
	return gt + GameTime(days*ticsPerDay)
}

func (gt GameTime) DateTime() string {
	return fmt.Sprintf("%03d-%04d %02d:%02d:%02d", gt.Day(), gt.Year(), gt.Hour(), gt.Minutes(), gt.Seconds())
}

func (gt GameTime) DaysUntilWeekStart() int {
	days := 0
	for gt.Weekday() != Wonday {
		gt = gt.AddDays(1)
		days++
	}
	return days
}

func (gt GameTime) DaysUntilMonthStart() int {
	days := 0
	for gt.DayOfMonth() != 1 {
		gt = gt.AddDays(1)
		days++
	}
	return days
}

func (gt GameTime) DaysUntilYearStart() int {
	switch d := gt.Day(); d {
	case 1:
		return 0
	default:
		return daysPerYear - d + 1
	}
}
