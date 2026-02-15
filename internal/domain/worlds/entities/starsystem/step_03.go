package starsystem

import (
	"fmt"
	"strconv"

	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
)

func (b *Builder) runStep3(ss *starSystemPrecursor) error {
	// - [ ] 3. **Determine system's worlds** (*import if available*)
	//   - [+] a. Determine gas giant (GG) presence and quantity from table
	if err := b.determineGasGiantsPresenceAndNumber(ss); err != nil {
		return fmt.Errorf("failed to determine gas giants: %w", err)
	}
	//   - [ ] b. Determine planetoid belt (PB) presence and quantity from table
	if err := b.determineAsteroidBeltsPresenceAndNumber(ss); err != nil {
		return fmt.Errorf("failed to determine gas giants: %w", err)
	}
	//   - [ ] c. Determine terrestrial planet (TP) quantity
	b.determineTerrestrialWorldsNumber(ss)
	//   - [ ] d. Record total worlds (GG + PB + TP)
	ss.TotalWorlds = ss.GG + ss.Belts + ss.Planets
	if err := step3Validation(ss); err != nil {
		return err
	}
	return nil
}

func (b *Builder) determineGasGiantsPresenceAndNumber(ss *starSystemPrecursor) error {
	if b.rng.Roll("2d6") < 10 {
		mods := []string{}
		if len(ss.Stars) == 1 && ss.PrimaryStar.Class == "V" {
			mods = append(mods, rtg.MOD_SingleClassV)
		}
		if ss.PrimaryStar.Type == "BD" {
			mods = append(mods, rtg.MOD_PrimaryIsBD)
		}
		if ss.PrimaryStar.Dead {
			mods = append(mods, rtg.MOD_PrimaryIsPostStellar)
		}
		if len(ss.Stars) >= 4 {
			mods = append(mods, rtg.MOD_FourOrMoreStars)
		}
		for _, star := range ss.Stars {
			if star.Dead {
				mods = append(mods, rtg.MOD_PerEveryPostStellar)
			}
		}

		gg, err := b.step3.tablesPlanets.Roll(rtg.TableGGQuantity, mods...)
		if err != nil {
			return fmt.Errorf("failed to roll table: %v", err)
		}
		ss.GG, err = strconv.Atoi(gg)
		if err != nil {
			return fmt.Errorf("bad result rolled (%v): %w", gg, err)
		}

	} else {
		ss.GG = 0
	}
	return nil
}

func (b *Builder) determineAsteroidBeltsPresenceAndNumber(ss *starSystemPrecursor) error {
	if b.rng.Roll("2d6") >= 8 {
		mods := []string{}
		if ss.GG >= 1 {
			mods = append(mods, rtg.MOD_SystemHas1orMoreGG)
		}
		if ss.PrimaryStar.Protostar {
			mods = append(mods, rtg.MOD_PrimaryIsProtostar)
		}
		if ss.PrimaryStar.Age <= 0.1 {
			mods = append(mods, rtg.MOD_PrimaryIsPrimordial)
		}
		if ss.PrimaryStar.Dead {
			mods = append(mods, rtg.MOD_PrimaryIsPostStellar)
		}
		for _, star := range ss.Stars {
			if star.Dead {
				mods = append(mods, rtg.MOD_PerEveryPostStellar)
			}
		}
		if len(ss.Stars) > 1 {
			mods = append(mods, rtg.MOD_SystemHas2orMoreStars)
		}
		bl, err := b.step3.tablesPlanets.Roll(rtg.TableBeltsQuantity, mods...)
		if err != nil {
			return fmt.Errorf("failed to roll table: %v", err)
		}
		ss.Belts, err = strconv.Atoi(bl)
		if err != nil {
			return fmt.Errorf("bad result rolled (%v): %w", bl, err)
		}

	} else {
		ss.Belts = 0
	}
	return nil
}

func (b *Builder) determineTerrestrialWorldsNumber(ss *starSystemPrecursor) error {
	if ss.Planets > -1 {
		return nil
	}
	r := b.rng.Roll("2d6-2")
	dm := 0
	if r >= 3 {
		dm = b.rng.Roll("1d3-1")
	}
	if r < 3 {
		r = b.rng.Roll("1d3+2")
	}
	sit := newStarIterator(ss.Stars)
	for sit.next() {
		_, star, _ := sit.getValues()
		if star.Dead {
			dm--
		}
	}
	ss.Planets = r + dm

	return nil
}

func step3Validation(ss *starSystemPrecursor) error {
	if ss.GG < 0 {
		return fmt.Errorf("negative gas giants quantity")
	}
	if ss.Belts < 0 {
		return fmt.Errorf("negative belts quantity")
	}
	if ss.Planets < 0 {
		return fmt.Errorf("negative terrestrial planets quantity")
	}
	if ss.TotalWorlds != ss.GG+ss.Belts+ss.Planets {
		return fmt.Errorf("total worlds not match")
	}
	return nil
}
