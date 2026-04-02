package starsystem

import (
	"fmt"
	"strconv"

	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
)

// runStep3 executes Step 3 of star system generation: determining system worlds.
// This includes:
// 1. Determining gas giant presence and quantity
// 2. Determining asteroid belt presence and quantity
// 3. Determining terrestrial planet quantity
// 4. Recording total worlds
//
// Note: This is partially implemented - steps 2b, 2c are marked as not implemented in World_Creation.md
func (builder *Builder) runStep03(systemPrecursor *starSystemPrecursor) error {
	if err := builder.determineGasGiantsPresenceAndNumber(systemPrecursor); err != nil {
		return fmt.Errorf("failed to determine gas giants: %w", err)
	}
	// TODO: Implement step 3b - Determine planetoid belt presence and quantity from table
	if err := builder.determineAsteroidBeltsPresenceAndNumber(systemPrecursor); err != nil {
		return fmt.Errorf("failed to determine asteroid belts: %w", err)
	}
	// TODO: Implement step 3c - Determine terrestrial planet quantity (currently uses simple logic)
	builder.determineTerrestrialWorldsNumber(systemPrecursor)
	systemPrecursor.TotalWorlds = systemPrecursor.GG + systemPrecursor.Belts + systemPrecursor.Planets
	if err := validateStep03Results(systemPrecursor); err != nil {
		return err
	}
	return nil
}

// determineGasGiantsPresenceAndNumber determines if the system has gas giants and how many.
// A 2d6 roll less than 10 means gas giants are possible, then modifiers are applied
// based on star type, number of stars, and other factors.
func (builder *Builder) determineGasGiantsPresenceAndNumber(systemPrecursor *starSystemPrecursor) error {
	if builder.rng.Roll("2d6") < 10 {
		mods := []string{}
		if len(systemPrecursor.Stars) == 1 && systemPrecursor.PrimaryStar.Class == "V" {
			mods = append(mods, rtg.MOD_SingleClassV)
		}
		if systemPrecursor.PrimaryStar.Type == "BD" {
			mods = append(mods, rtg.MOD_PrimaryIsBD)
		}
		if systemPrecursor.PrimaryStar.Dead {
			mods = append(mods, rtg.MOD_PrimaryIsPostStellar)
		}
		if len(systemPrecursor.Stars) >= 4 {
			mods = append(mods, rtg.MOD_FourOrMoreStars)
		}
		for _, star := range systemPrecursor.Stars {
			if star.Dead {
				mods = append(mods, rtg.MOD_PerEveryPostStellar)
			}
		}

		gg, err := builder.step3.tablesPlanets.Roll(rtg.TableGGQuantity, mods...)
		if err != nil {
			return fmt.Errorf("failed to roll table: %v", err)
		}
		systemPrecursor.GG, err = strconv.Atoi(gg)
		if err != nil {
			return fmt.Errorf("bad result rolled (%v): %w", gg, err)
		}

	} else {
		systemPrecursor.GG = 0
	}
	return nil
}

// determineAsteroidBeltsPresenceAndNumber determines if the system has asteroid belts and how many.
// A 2d6 roll of 8 or higher means asteroid belts are possible, with modifiers applied
// based on gas giant presence, star type, age, and other factors.
func (builder *Builder) determineAsteroidBeltsPresenceAndNumber(systemPrecursor *starSystemPrecursor) error {
	if builder.rng.Roll("2d6") >= 8 {
		mods := []string{}
		if systemPrecursor.GG >= 1 {
			mods = append(mods, rtg.MOD_SystemHas1orMoreGG)
		}
		if systemPrecursor.PrimaryStar.Protostar {
			mods = append(mods, rtg.MOD_PrimaryIsProtostar)
		}
		if systemPrecursor.PrimaryStar.Age <= 0.1 {
			mods = append(mods, rtg.MOD_PrimaryIsPrimordial)
		}
		if systemPrecursor.PrimaryStar.Dead {
			mods = append(mods, rtg.MOD_PrimaryIsPostStellar)
		}
		for _, star := range systemPrecursor.Stars {
			if star.Dead {
				mods = append(mods, rtg.MOD_PerEveryPostStellar)
			}
		}
		if len(systemPrecursor.Stars) > 1 {
			mods = append(mods, rtg.MOD_SystemHas2orMoreStars)
		}
		bl, err := builder.step3.tablesPlanets.Roll(rtg.TableBeltsQuantity, mods...)
		if err != nil {
			return fmt.Errorf("failed to roll table: %v", err)
		}
		systemPrecursor.Belts, err = strconv.Atoi(bl)
		if err != nil {
			return fmt.Errorf("bad result rolled (%v): %w", bl, err)
		}

	} else {
		systemPrecursor.Belts = 0
	}
	return nil
}

// determineTerrestrialWorldsNumber determines the number of terrestrial (rocky) planets in the system.
// Uses 2d6-2 as base with modifiers for dead stars in the system.
// TODO: Implement proper logic as specified in World_Creation.md step 3c
// TODO: Return error instead of ignoring it
func (builder *Builder) determineTerrestrialWorldsNumber(systemPrecursor *starSystemPrecursor) error {
	if systemPrecursor.Planets > -1 {
		return nil
	}
	baseRoll := builder.rng.Roll("2d6-2")
	diceModifier := 0
	if baseRoll >= 3 {
		diceModifier = builder.rng.Roll("1d3-1")
	}
	if baseRoll < 3 {
		baseRoll = builder.rng.Roll("1d3+2")
	}
	starIterator := newStarIterator(systemPrecursor.Stars)
	for starIterator.next() {
		_, star, _ := starIterator.getValues()
		if star.Dead {
			diceModifier--
		}
	}
	systemPrecursor.Planets = baseRoll + diceModifier

	return nil
}

// validateStep3Results validates the results of step 3 generation.
// Checks that all world counts are non-negative and that total matches the sum.
func validateStep03Results(systemPrecursor *starSystemPrecursor) error {
	if systemPrecursor.GG < 0 {
		return fmt.Errorf("negative gas giants quantity")
	}
	if systemPrecursor.Belts < 0 {
		return fmt.Errorf("negative belts quantity")
	}
	if systemPrecursor.Planets < 0 {
		return fmt.Errorf("negative terrestrial planets quantity")
	}
	if systemPrecursor.TotalWorlds != systemPrecursor.GG+systemPrecursor.Belts+systemPrecursor.Planets {
		return fmt.Errorf("total worlds not match")
	}
	return nil
}
