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
		global.NextSpaceRound()
		fmt.Printf("%s %v          \r", global.Time().DateTime(), global.Time().DayWeek())
		time.Sleep(time.Millisecond)
	}

	currentTime := global.Time()

	fmt.Printf("Current time: %s\n", currentTime.Format())
	fmt.Printf("Year: %d\n", currentTime.Year())
	fmt.Printf("Month: %d\n", currentTime.Month())
	fmt.Printf("Day of year: %d\n", currentTime.Day())
	fmt.Printf("Day of month: %d\n", currentTime.DayOfMonth())
	fmt.Printf("Day of week: %d\n", currentTime.DayWeek())
	fmt.Printf("Hour: %d\n", currentTime.Hour())
	fmt.Printf("Minutes: %d\n", currentTime.Minutes())
	fmt.Printf("Is holiday: %t\n", currentTime.IsHoliday())
}
