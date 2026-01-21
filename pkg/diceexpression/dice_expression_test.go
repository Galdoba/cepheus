package diceexpression

import (
	"fmt"
	"testing"
)

func Test_parseBaseSum(t *testing.T) {
	for _, s := range []string{
		"2d6",
		" 2d6 2d6",
		"2D6",
		"2d6+3",
		"5dd2",
		"d6",
	} {
		fmt.Printf(" input: %v\n", s)
		n, f, err := parseBaseSum(s)
		fmt.Printf("output: %v %v (%v)\n", n, f, err)
		fmt.Println("")

	}
	fmt.Println("SUM")
	for _, s := range []string{
		"2d6+1",
		"2d6-6",
		"2D6-2+3",
		"2d6+3-1",
		"2d6+3a5-1",
		"3d6-5+3-+2",
	} {
		fmt.Printf(" input: %v\n", s)
		n, err := parseAdditiveMod(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")

	}
	fmt.Println("MULT")
	for _, s := range []string{
		"2d6x1",
		"2d6x6",
		"2D6x3x3",
		"2d6x33x-1",
		"2d6xx44",
		"3d6x1000x5+6x5",
	} {
		fmt.Printf(" input: %v\n", s)
		n, err := parseMultiplicativeMod(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")

	}
	fmt.Println("DEL")
	for _, s := range []string{
		"2d6/2",
		"2d6/2/6",
	} {
		fmt.Printf(" input: %v\n", s)
		n, err := parseDeletiveMod(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")

	}
	fmt.Println("REPLACE")
	for _, s := range []string{
		"2d6r1:2",
		"2d6r1:2r4:4",
		"2d6r1:2r1:3",
	} {
		fmt.Printf(" input: %v\n", s)
		n, err := parseReplacements(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")

	}
	fmt.Println("DROP")
	for _, s := range []string{
		"4d6r1:2dh2",
		"3d6dl1+1",
	} {
		fmt.Printf(" input: %v\n", s)
		n, err := parseDropHigh(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")
		n, err = parseDropLow(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")

	}
	fmt.Println("INDIVIDUAL")
	for _, s := range []string{
		"4dd6i+1",
		"3d6i+1i-5",
	} {
		fmt.Printf(" input: %v\n", s)
		n, err := parseIndividualMod(s)
		fmt.Printf("output: %v (%v)\n", n, err)
		fmt.Println("")

	}
}
