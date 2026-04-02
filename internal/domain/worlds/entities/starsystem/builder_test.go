package starsystem

import (
	"fmt"
	"testing"
)

// TestBuilder tests the star system builder by generating systems with seeds from 160000 to 999999.
// It prints unique system profiles to help identify generation patterns and potential issues.
// TODO: Convert to proper test with assertions instead of print statements
// TODO: Add test for specific edge cases (single star, multiple stars, special types)
func TestBuilder(t *testing.T) {
	lastProfile := ""
	for i := 160000; i < 1000000; i++ {
		builder, err := NewBuilder(fmt.Sprintf("%v", i))
		if err != nil {
			fmt.Println(err)
			return
		}
		systemPrecursor, err := builder.Build()
		if err != nil {
			fmt.Println(err)
			return
		}
		if systemPrecursor.Profile() != lastProfile {
			fmt.Println(systemPrecursor.Profile())
			lastProfile = systemPrecursor.Profile()

		}
		si := newStarIterator(systemPrecursor.Stars)
		for si.next() {
			_, star, _ := si.getValues()
			fmt.Println(star.Designation)
		}
		fmt.Println(systemPrecursor.Stars)
		// if len(systemPrecursor.Stars) > 3 {
		// 	return
		// }
	}
}
