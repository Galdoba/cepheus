package actions

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/cepheus/internal/domain/support/entities/paths"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
	"github.com/urfave/cli/v3"
)

func TrvWorlds_Status(app *app.TrvWorldsInfrastructure) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		fmt.Println("here be dragons...")
		fmt.Println("config:", app.Config)
		file, err := os.Stat(app.Config.Import.External_DB_File)
		if err != nil {
			switch errors.Is(err, os.ErrNotExist) {
			case true:
				fmt.Printf("no import file was found.\n")
				fmt.Printf("run `trv_worlds import` to download data from https://travellermap.com\n")
				return nil
			case false:
				return fmt.Errorf("failed to read import file: %v\n", err)
			}
		}
		fmt.Printf("database detected: %v bytes\n", file.Size())
		fmt.Println("test read database...")
		db, err := jsonstorage.OpenStorage[t5ss.WorldData](paths.DefaultExternalDB_File())
		if err != nil {
			return fmt.Errorf("failed to open database: %v\n", err)
		}
		fmt.Printf("database contains %v entries\n", db.Len())
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database")
		}
		fmt.Println("everything seems fine to me... ")
		time.Sleep(time.Second * 2)
		fmt.Println("or is it?")
		fmt.Println("")
		return nil
	}
}
