package starsystem

import (
	"fmt"
	"strconv"

	"github.com/Galdoba/cepheus/internal/infrastructure/rtg"
)

func (b *Builder) runStep3(ss *StarSystem) error {
	// - [ ] 3. **Determine system's worlds** (*import if available*)
	//   - [+] a. Determine gas giant (GG) presence and quantity from table
	if err := b.determineGasGiantsPresenceAndNumber(ss); err != nil {
		return fmt.Errorf("failed to determine gas giants: %w", err)
	}
	//   - [ ] b. Determine planetoid belt (PB) presence and quantity from table
	//   - [ ] c. Determine terrestrial planet (TP) quantity
	//   - [ ] d. Record total worlds (GG + PB + TP)
	return nil
}

func (b *Builder) determineGasGiantsPresenceAndNumber(ss *StarSystem) error {
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
		fmt.Println(mods)

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
	fmt.Println(ss.GG)
	return nil
}
