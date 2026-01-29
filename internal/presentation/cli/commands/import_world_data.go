package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
	"github.com/Galdoba/cepheus/internal/presentation/api"
	"github.com/Galdoba/cepheus/internal/presentation/cli/flags"
	"github.com/urfave/cli/v3"
)

const (
	IMPORT_WORLD_DATA_COMMAND = "import"
)

func ImportWorldData(app *app.TrvWorldsInfrastructure) *cli.Command {
	add := cli.Command{
		Name:        IMPORT_WORLD_DATA_COMMAND,
		Aliases:     []string{},
		Usage:       "download world data in t5ss format",
		UsageText:   "trv_worlds [global options] import [options]",
		Description: "Download jumpmaps with t5ss data from https://travellermap.com",
		Action:      importAction(app.Config),
		Flags:       flags.TrvWorlds_Import(),
	}
	return &add
}

func importAction(cfg config.TrvWorldsCfg) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		dbPath := cfg.Import.External_DB_File
		rings := cfg.Import.CoordinatesRingSize
		if c.Int(flags.RINGS) != 0 {
			rings = c.Int(flags.RINGS)
		}

		//open db
		db, err := jsonstorage.OpenOrCreateStorage[t5ss.WorldData](dbPath)
		if err != nil {
			return err
		}

		//setup request links
		ursList := []string{}
		switch c.Bool(flags.DRY_RUN) {
		case false:
			ursList = api.ImportUrlList(rings)
		case true:
			fmt.Fprintf(os.Stderr, "warning: dry-run mode activated\n")
		}

		//download data
		worldDatamap, errormap := api.GetData(ursList...)
		for url, err := range errormap {
			fmt.Printf("failed to get data from [%v]:\nerror: %v\n", url, err)
		}

		//write to db and close
		updated := 0
		created := 0
		for url, data := range worldDatamap {
			batch := t5ss.WorldBatch{}
			if err := json.Unmarshal(data, &batch); err != nil {
				fmt.Printf("failed to unmarshal data from [%v]:\nerror: %v\n", url, err)
			}
			for _, surveyImported := range batch.List {
				crd := surveyImported.Coordinates()
				if err := db.Update(crd.DatabaseKey(), surveyImported); err == nil {
					updated++
				}
				if err := db.Create(crd.DatabaseKey(), surveyImported); err == nil {
					created++
				}
			}

		}
		if err := db.CommitAndClose(); err != nil {
			return fmt.Errorf("failed to commit&&close import database: %v", err)
		}

		//exit stats
		if updated > 0 {
			fmt.Println("database entries updated:", updated)
		}
		if created > 0 {
			fmt.Println("database entries created:", created)
		}
		return nil
	}
}
