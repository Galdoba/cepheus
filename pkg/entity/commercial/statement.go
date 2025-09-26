package commercial

import "github.com/Galdoba/cepheus/pkg/dice"

const (
	Aggressive  MissionStatementType = "Aggressive"
	Altruistic  MissionStatementType = "Altruistic"
	Defensive   MissionStatementType = "Defensive"
	Greedy      MissionStatementType = "Greedy"
	Mysterious  MissionStatementType = "Mysterious"
	Outrageus   MissionStatementType = "Outrageus"
	Simple      MissionStatementType = "Simple"
	Threatening MissionStatementType = "Threatening"
)

type MissionStatementType string

type MissionStatement struct {
	Type          MissionStatementType
	Control       int
	Dependability int
	Guide         int
	Management    int
}

func randomStatement() MissionStatement {
	return []MissionStatement{
		AggressiveStatemet,
		AltruisticStatemet,
		DefensiveStatemet,
		GreedyStatemet,
		MysteriousStatemet,
		OutrageusStatemet,
		SimpleStatemet,
		ThreateningStatemet,
	}[dice.NewDicepool().Sum("1d8")-1]
}

var noStatement = MissionStatement{}

var AggressiveStatemet = MissionStatement{
	Type:          Aggressive,
	Control:       9,
	Dependability: 7,
	Guide:         8,
	Management:    8,
}
var AltruisticStatemet = MissionStatement{
	Type:          Altruistic,
	Control:       7,
	Dependability: 9,
	Guide:         6,
	Management:    7,
}
var DefensiveStatemet = MissionStatement{
	Type:          Defensive,
	Control:       7,
	Dependability: 9,
	Guide:         8,
	Management:    8,
}
var GreedyStatemet = MissionStatement{
	Type:          Greedy,
	Control:       10,
	Dependability: 5,
	Guide:         9,
	Management:    8,
}
var MysteriousStatemet = MissionStatement{
	Type:          Mysterious,
	Control:       6,
	Dependability: 6,
	Guide:         8,
	Management:    9,
}
var OutrageusStatemet = MissionStatement{
	Type:          Outrageus,
	Control:       8,
	Dependability: 9,
	Guide:         8,
	Management:    6,
}
var SimpleStatemet = MissionStatement{
	Type:          Simple,
	Control:       7,
	Dependability: 9,
	Guide:         5,
	Management:    10,
}
var ThreateningStatemet = MissionStatement{
	Type:          Threatening,
	Control:       9,
	Dependability: 7,
	Guide:         9,
	Management:    7,
}
