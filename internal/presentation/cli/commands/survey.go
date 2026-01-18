package commands

import (
	"context"
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/aggregates/world"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
	"github.com/Galdoba/cepheus/internal/presentation/cli/flags"
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

		for _, key := range canonicalDataStorage.AllKeys() {
			wd, err := canonicalDataStorage.Read(key)
			if err != nil {
				return fmt.Errorf("failed to read canonical data (%v): %v", wd.Import_DB_Key(), err)
			}
			w, err := world.Import(wd)
			if err != nil {
				return fmt.Errorf("failed to import (%v): %v", wd.Import_DB_Key(), err)
			}
			dto := w.ToDTO()
			workingDataStorage.Update(dto.Key(), w.ToDTO())
			workingDataStorage.Create(dto.Key(), w.ToDTO())
		}
		canonicalDataStorage.Close()
		workingDataStorage.CommitAndClose()

		return nil
	}
}

var worldDB_path string
var importDB_path string
