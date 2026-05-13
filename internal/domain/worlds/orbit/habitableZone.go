package orbit

import "fmt"

type HabitableZone [2]float64

func (hz HabitableZone) Validate() error {
	if hz[0] < 0 {
		return fmt.Errorf("habiatable zone minimum is less than 0")
	}
	if hz[1] < 0 {
		return fmt.Errorf("habiatable zone maximum is less than 0")
	}
	if hz[1] < hz[0] {
		return fmt.Errorf("habiatable zone maximum is less minimum")
	}
	return nil
}

func (hz HabitableZone) Center() float64 {
	return (hz[1] + hz[0]) / 2
}

func (hz HabitableZone) Width() float64 {
	return hz[1] - hz[0]
}
