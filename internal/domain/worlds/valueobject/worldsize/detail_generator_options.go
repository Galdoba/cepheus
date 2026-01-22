package worldsize

import "github.com/Galdoba/cepheus/pkg/dice"

type SizeDetails func(*DetailGenerator)

func WithHZCO(hzco float64) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.hzco = hzco
	}
}

func WithRNG(dp *dice.Roller) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.rng = dp
	}
}

func WithDiameter(diam int) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.forcedDiameter = diam
	}
}

func WithComposition(composition string) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.forcedComposition = composition
	}
}

func WithDencity(dencity float64) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.forcedDencity = dencity
	}
}

func WithGravity(gravity float64) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.forcedGravity = gravity
	}
}

func WithMass(mass float64) SizeDetails {
	return func(dg *DetailGenerator) {
		dg.forcedMass = mass
	}
}
