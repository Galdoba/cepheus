package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/worlds/aggregates/world"
	"github.com/Galdoba/cepheus/internal/domain/worlds/entities/trade"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/tradegoods"
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
		fmt.Println("  ...sort of")
		time.Sleep(time.Second)
		crdList := coordinates.Spiral(wd.Coordinates().ToCube(), 4)
		// localUWP := uwp.UWP(wd.UWP)
		// ours := tradegoods.Available(classifications.Classify(localUWP)...)
		for _, crd := range crdList {
			if ok, err := trade.Exists(wd.Coordinates(), crd.ToGlobal()); err == nil {
				switch ok {
				case true:
					partner, _ := canonicalDataStorage.Read(crd.ToGlobal().DatabaseKey())
					fmt.Println("trades with", partner.SearchKey())
					importing, err := trade.CalculateImport(wd.Coordinates(), partner.Coordinates())
					if err != nil {
						panic(err)
					}
					if len(importing) > 0 {
						fmt.Println("importing:")
						for _, good := range tradegoods.Types(importing...) {
							fmt.Println(good)
						}
					}

				case false:

				}
			}

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
