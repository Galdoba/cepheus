package actions

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/urfave/cli/v3"
)

func TrvWorlds_Status(app *app.TrvWorldsInfrastructure) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		fmt.Println("here be dragons...")
		fmt.Println("config:", app.Config)
		file, err := os.Stat(app.Config.Import.ImportDataPath)
		if err != nil {
			switch errors.Is(err, os.ErrNotExist) {
			case true:
				fmt.Printf("no import file was found.\n")
				fmt.Printf("run `trv_worlds import` to download data from https://travellermap.com\n")
				return nil
			case false:
				return fmt.Errorf("failed to read import file: %v", err)
			}
		}
		fmt.Printf("import file detected: %v bytes", file.Size())
		return nil
	}
}
