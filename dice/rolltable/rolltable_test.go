package rolltable

import (
	"fmt"
	"testing"

	"github.com/Galdoba/t5/pkg/dice"
)

func TestNewTable(t *testing.T) {
	starports := NewTable(
		WithDescription("Starports on Mainworld"),
		WithDiceCode("1d6"),
		WithEntries(
			NewEntry("A", WithKey("0...4"),
				WithStringValue("Quality", "Exelent"),
				WithStringValue("Repairs", "Overhaul"),
				WithStringValue("Downport", "Yes"),
			),
			NewEntry("B", WithKey("5,6"),
				WithStringValue("Quality", "Good"),
				WithStringValue("Repairs", "Overhaul"),
				WithStringValue("Downport", "Yes"),
			),
			NewEntry("C", WithKey("7,8"),
				WithStringValue("Quality", "Exelent"),
				WithStringValue("Repairs", "Overhaul"),
				WithStringValue("Downport", "Yes"),
			),
			NewEntry("D", WithKey("9"),
				WithStringValue("Quality", "Exelent"),
				WithStringValue("Repairs", "Overhaul"),
				WithStringValue("Downport", "Yes"),
			),
			NewEntry("E", WithKey("10,11"),
				WithStringValue("Quality", "Exelent"),
				WithStringValue("Repairs", "Overhaul"),
				WithStringValue("Downport", "Yes"),
			),
			NewEntry("X", WithKey("12+"),
				WithStringValue("Quality", "Exelent"),
				WithStringValue("Repairs", "Overhaul"),
				WithStringValue("Downport", "Yes"),
			),
		),
	)
	result := starports.Check(dice.NewDicepool())

	fmt.Println(result)
}
