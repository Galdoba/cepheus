package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/domain/generic/entities/coordinates"
	"github.com/Galdoba/cepheus/internal/domain/generic/services/travellermap"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
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
		fmt.Printf("storage contains data on %v worlds\n", js.Len())

		wd, err := travellermap.GetWorldData(coordinates.NewSpaceCoordinates(0, 0), cfg.Import.CoordinatesRingSize)
		fmt.Println(err)
		fmt.Println(len(wd), "worlds:")
		for _, w := range wd {
			fmt.Println(w)
		}

		js.Close()
		return nil
	}
}
