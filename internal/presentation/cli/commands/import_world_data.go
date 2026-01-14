package commands

import (
	"context"
	"encoding/json"
	"errors"
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
		Name:      IMPORT_WORLD_DATA_COMMAND,
		Aliases:   []string{},
		Usage:     "download world data in t5ss format",
		UsageText: "trv_worlds [global options] import [options]",
		Action:    importAction(*app.Config),
		Flags:     flags.TrvWorlds_Import(),
	}
	return &add
}

func importAction(cfg config.TrvWorldsCfg) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		dbPath := cfg.Import.ImportDataPath

		//Open Storage
		js, err := jsonstorage.OpenStorage[t5ss.WorldData](dbPath)
		if err != nil {
			switch errors.Is(err, os.ErrNotExist) {
			case true:
				fmt.Printf("import storage does not exits!\ncreate new: ")
				js, err = jsonstorage.NewStorage[t5ss.WorldData](dbPath)
				if err != nil {
					fmt.Println("failed!")
					fmt.Println("aborting program...")
					return err
				}
			case false:
				return fmt.Errorf("failed to create new storage: %v", err)
			}
		}

		worldDatamap, errormap := api.GetData(api.ImportUrlList(14)...)
		fmt.Println("")
		for url, err := range errormap {
			fmt.Printf("failed to get data from [%v]:\nerror: %v\n", url, err)
		}
		updated := 0
		created := 0
		for url, data := range worldDatamap {
			batch := t5ss.WorldBatch{}
			// fmt.Println(data)
			// fmt.Println(string(data))
			if err := json.Unmarshal(data, &batch); err != nil {
				fmt.Printf("failed to unmarshal data from [%v]:\nerror: %v\n", url, err)
			}
			// fmt.Println(batch)
			// fmt.Println(len(batch.List))
			for _, world := range batch.List {
				crd := world.Coordinates()
				if err := js.Update(crd.DatabaseKey(), world); err == nil {
					updated++
				}
				if err := js.Create(crd.DatabaseKey(), world); err == nil {
					created++
				}
			}

		}
		if err := js.CommitAndClose(); err != nil {
			return fmt.Errorf("failed to commit&&close database: %v", err)
		}
		if updated > 0 {
			fmt.Println("database entries updated:", updated)
		}
		if created > 0 {
			fmt.Println("database entries created:", created)
		}
		return nil
	}
}
