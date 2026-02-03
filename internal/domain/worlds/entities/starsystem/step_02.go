package starsystem

import "github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/stellar"

func (b *Builder) runStep2(ss *StarSystem) error {
	// - [+] 2. **Determine if system has multiple stars, if yes, then:**
	b.step2.starSchema = stellar.RollDesignations(b.rng)
	if b.imported.Allegiance != "" && b.imported.Stellar != "" {
		b.step2.starSchema = stellar.RollStellarDesignations(b.rng, stellar.Stellar(b.imported.Stellar))
	}
	if len(b.step2.starSchema) < 2 {
		b.step2.completed = true
		return nil
	}

	//   - [ ] a. Determine Orbit#s of secondary and companion stars

	//
	//   - [ ] b. Determine eccentricity of secondary stars and check for overlaps
	//   - [ ] c. Determine secondary and companion star types
	//   - [ ] d. Adjust system age to account for post-stellar objects (if any)
	//   - [ ] e. Determine star orbital periods
	return nil
}
