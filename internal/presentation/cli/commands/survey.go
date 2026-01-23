package commands

import (
	"context"
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/aggregates/world"
	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/astrogation"
	"github.com/Galdoba/cepheus/internal/domain/worlds/services/traderoute"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
	"github.com/Galdoba/cepheus/internal/presentation/cli/flags"
	"github.com/Galdoba/consolio/prompt"
	"github.com/urfave/cli/v3"
)

const (
	SURVEY_WORLD_COMMAND = "survey"
)

func SurveyWorld(app *app.TrvWorldsInfrastructure) *cli.Command {
	add := cli.Command{
		Name:        SURVEY_WORLD_COMMAND,
		Aliases:     []string{},
		Usage:       "fil derived_db with external_db data",
		UsageText:   "trv_worlds [global options] inspect [options]",
		Description: "Transfer world data to working database\nAdditional data will be calculated and stored.",
		Action:      surveyAction(*app.Config),
		Flags:       flags.TrvWorlds_Survey(),
	}
	return &add
}

func surveyAction(cfg config.TrvWorldsCfg) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		external_db_path := cfg.Import.External_DB_File
		derived_db_path := cfg.World.Derived_DB_File
		//open
		fmt.Printf("open database... ")
		workingDataStorage, err := jsonstorage.OpenOrCreateStorage[world.WorldDTO](derived_db_path)
		if err != nil {
			return fmt.Errorf("failed open working data storage: %v", err)
		}
		canonicalDataStorage, err := jsonstorage.OpenStorage[t5ss.WorldData](external_db_path)
		if err != nil {
			return fmt.Errorf("failed to open canonical data storage")
		}
		fmt.Println("success!!")

		searcKeys := []string{}
		items := []*prompt.Item{}

		for _, entry := range canonicalDataStorage.AllEntries() {
			sk := entry.SearchKey()
			searcKeys = append(searcKeys, sk)
			items = append(items, prompt.NewItem(sk, entry))
		}

		selected, err := prompt.SearchItem(
			prompt.FromItems(items),
			prompt.WithTitle("Select world to survey:"),
		)

		if err != nil {
			return err
		}

		wd := selected.Payload().(t5ss.WorldData)

		fmt.Printf("selected entry: %v\n", wd)

		w, err := world.Import(wd)
		if err != nil {
			return err
		}
		fmt.Println("simulation hydration:")
		fmt.Println("Trade Routes:")
		as, err := astrogation.New()
		if err != nil {
			return err
		}
		crdList := coordinates.Spiral(wd.Coordinates().ToCube(), 4)
		for _, crd := range crdList {
			if !as.TradePathExist(w.Coordinates(), crd) {
				continue
			}
			partner, _ := canonicalDataStorage.Read(crd.ToGlobal().DatabaseKey())
			if !traderoute.HasLink(wd, partner) {
				continue
			}
			imp, exp := traderoute.Calculate(wd, partner)
			if len(imp)+len(exp) == 0 {
				continue
			}
			fmt.Println("\ntrades with", partner.SearchKey(), ":")
			for i, goods := range imp {
				if i == 0 {
					fmt.Println("importing:")
				}
				fmt.Println("  ", goods.TradeGoodType)

			}
			for i, goods := range exp {
				if i == 0 {
					fmt.Println("exporting:")
				}
				fmt.Println("  ", goods.TradeGoodType)
			}
			w.CreateTradeConnection(crd.ToGlobal(), imp, exp)
		}

		workingDataStorage.Update(w.DatabaseKey(), w.ToDTO())
		workingDataStorage.Create(w.DatabaseKey(), w.ToDTO())
		canonicalDataStorage.Close()
		workingDataStorage.CommitAndClose()

		return nil
	}
}

var worldDB_path string
var importDB_path string
