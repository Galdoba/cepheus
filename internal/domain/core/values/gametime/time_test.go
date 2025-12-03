package gametime

import (
	"fmt"
	"testing"
	"time"
)

func TestGlobal_NextCombatRound(t *testing.T) {
	var global Global

	// Симулируем прохождение времени
	for range 100000000 {
		global.NextPortRound()
		fmt.Printf("%s %v-%v-%v %v  dow=%v     \r", global.Time().DateTime(), global.Time().Month(), global.Time().WeekOfMonth(), global.Time().DayOfWeek(), global.Time().Weekday(), global.Time().DayOfMonth())
		time.Sleep(time.Millisecond * 20)
	}

	currentTime := global.Time()

	fmt.Printf("Year: %d\n", currentTime.Year())
	fmt.Printf("Month: %d\n", currentTime.Month())
	fmt.Printf("Day of year: %d\n", currentTime.Day())
	fmt.Printf("Day of month: %d\n", currentTime.DayOfMonth())
	fmt.Printf("Day of week: %d\n", currentTime.DayOfWeek())
	fmt.Printf("Hour: %d\n", currentTime.Hour())
	fmt.Printf("Minutes: %d\n", currentTime.Minutes())
	fmt.Printf("Is holiday: %t\n", currentTime.IsHoliday())
}
