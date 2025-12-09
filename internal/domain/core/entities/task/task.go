package task

import (
	"github.com/Galdoba/cepheus/internal/domain/core/entities/check"
	"github.com/Galdoba/cepheus/internal/domain/core/entities/dice"
	"github.com/Galdoba/cepheus/internal/domain/core/values/characteristic"
	"github.com/Galdoba/cepheus/internal/domain/core/values/gametime"
	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
)

const (
	Undefined      TimeFrame = iota
	Instant                  //no time
	Seconds                  //1d6 seconds
	CombatRounds             //1d6 * 6 seconds
	Minutes                  //1d6 minutes
	About10Minutes           //10+flux minutes (5-15)
	AnHour                   //60+flux*10 minutes (10 - 110)
	AllDay                   // 10 + flux*hours (5 - 15)
	AWeek                    //6 + flux*days (1 - 11 days)
	AMonth                   //6 + flux * week (1 - 11 weeks)
)

type TimeFrame int

type Task struct {
	Description           string
	RelatedCharacteristic characteristic.CharacteristicName
	RelatedSkill          skill.Skill
	TimeFrame             int
	StartAt               gametime.GameTime
	FinishAt              gametime.GameTime
	Check                 *check.Check
}

func newTask() Task {
	ts := Task{
		Description:           "",
		RelatedCharacteristic: "",
		RelatedSkill:          "",
		TimeFrame:             0,
		StartAt:               0,
		FinishAt:              0,
		Check:                 check.New(),
	}
	return ts
}

func calculateTimeframe(tf TimeFrame, start gametime.GameTime) gametime.GameTime {
	switch tf {
	case Instant:
		return start
	case Seconds:
		return start.AddTicks(1)
	case CombatRounds:
		ticks := gametime.GameTime(6 * dice.Roll(dice.Code_1D))
		return start.AddTicks(ticks)
	case Minutes:
		ticks := gametime.GameTime(gametime.Minute * gametime.GameTime(dice.Roll(dice.Code_1D)))
		return start.AddTicks(ticks)
	case About10Minutes:
		ticks := gametime.GameTime((gametime.Minute * 10) + gametime.GameTime(dice.Flux()))
		return start.AddTicks(ticks)
	case AnHour:
		ticks := gametime.GameTime((gametime.Minute * 60) + gametime.GameTime(gametime.Minute*10*gametime.GameTime(dice.Flux())))
		return start.AddTicks(ticks)
	case AllDay:
		ticks := gametime.GameTime(gametime.Hour*10 + (gametime.Hour * gametime.GameTime(dice.Flux())))
		return start.AddTicks(ticks)
	case AWeek:
		ticks := gametime.GameTime(gametime.Day*6 + (gametime.Day * gametime.GameTime(dice.Flux())))
		return start.AddTicks(ticks)
	case AMonth:
		ticks := gametime.GameTime(gametime.Day*6 + (gametime.Day * gametime.GameTime(dice.Flux())))
		for range 3 {
			ticks = ticks.AddTicks(gametime.GameTime(gametime.Day*6 + (gametime.Day * gametime.GameTime(dice.Flux()))))
		}
		return start.AddTicks(ticks)
	}
	return 0
}
