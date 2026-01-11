package money

import "fmt"

const (
	creditsInMegaCredit = 1000000
	megaCreditsInRU     = 100
)

type Credit int64

type MegaCredit float64

type RU int

func (c Credit) String() string {
	return fmt.Sprintf("%d cr", c)
}

func (mc MegaCredit) String() string {
	return fmt.Sprintf("%.6g MCr", mc)
}

func (c Credit) MegaCredit() MegaCredit {
	mc := float64(int64(c)) / creditsInMegaCredit
	return MegaCredit(mc)
}

func (mc MegaCredit) Credit() Credit {
	c := int64(float64(mc) * creditsInMegaCredit)
	return Credit(c)
}

func (mc MegaCredit) RU() RU {
	ru := int(mc) / megaCreditsInRU
	return RU(ru)
}

func (rc RU) MegaCredit() MegaCredit {
	mc := int(rc) * megaCreditsInRU
	return MegaCredit(mc)
}
