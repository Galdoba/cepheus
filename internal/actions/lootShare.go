package actions

import (
	"context"
	"fmt"

	"github.com/Galdoba/appcontext"
	lscconfig "github.com/Galdoba/cepheus/internal/config/loot-share-calculator"
	"github.com/urfave/cli/v3"
)

func CalculateLootShares(actx *appcontext.AppContext) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		cfg := lscconfig.Default()
		if err := actx.Config.LoadConfig("", &cfg); err != nil {
			return err
		}
		ships := []ShipData{
			ShipData{Tag: "Harrier now", Crew: 5, OfficerRatio: 10, Players: 5, MasteryMod: 0},
			ShipData{Tag: "Harrier full crew (green military)", Crew: 11, OfficerRatio: 10, Players: 5, MasteryMod: 0},
			ShipData{Tag: "Harrier full crew (green military no player)", Crew: 11, OfficerRatio: 10, Players: 0, MasteryMod: 0},
			ShipData{Tag: "Harrier full crew (green trader)", Crew: 11, OfficerRatio: 20, Players: 5, MasteryMod: 0},
			ShipData{Tag: "Harrier full crew (legendary)", Crew: 11, OfficerRatio: 10, Players: 5, MasteryMod: 3},
		}
		for _, ship := range ships {
			calcShipShares(ship)
		}

		return nil
	}
}

type ShipData struct {
	Tag          string
	Crew         int
	OfficerRatio int
	Players      int
	MasteryMod   int
}

func calcShipShares(sd ShipData) {
	loot := 1000000
	totalCrew := sd.Crew
	captains := 1
	shareCaptaib := captains * 5
	offciers := totalCrew / sd.OfficerRatio
	shareOfficers := offciers * 2
	ordinary := totalCrew - captains - offciers
	shareOrdinary := ordinary
	playerShare := sd.Players * 2
	shareMastery := sd.MasteryMod * totalCrew
	npcShares := shareCaptaib + shareOfficers + shareOrdinary + shareMastery

	sharesTotal := npcShares + playerShare
	oneShare := loot / sharesTotal
	fmt.Printf("example: %v\n", sd.Tag)
	fmt.Printf("crew: %v; 1 officer per %v crew; players on ship=%v; mastery mod = %v\n", sd.Crew, sd.OfficerRatio, sd.Players, sd.MasteryMod)
	fmt.Println("for loot of", loot, "share cost", oneShare)
	fmt.Println("captain player share:", (2+sd.MasteryMod+5)*oneShare)
	fmt.Println("captain npc share:", (sd.MasteryMod+5)*oneShare)
	fmt.Println("officer player share:", (sd.MasteryMod+2)*oneShare)
	fmt.Println("officer npc share:", (2+sd.MasteryMod+2)*oneShare)
	fmt.Println("ordinary player share:", (2+sd.MasteryMod+1)*oneShare)
	fmt.Println("ordinary share:", (sd.MasteryMod+1)*oneShare)
	fmt.Println("")
}
