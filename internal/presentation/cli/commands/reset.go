package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/urfave/cli/v3"
)

const (
	RESET_COMMAND = "reset"
)

func Reset(app *app.TrvWorldsInfrastructure) *cli.Command {
	add := cli.Command{
		Name:        RESET_COMMAND,
		Aliases:     []string{},
		Usage:       "reset program files",
		UsageText:   "trv_worlds [global options] inspect [options]",
		Description: "Choose files to reset: config, import_database, worlds_database",
		Action:      resetAction(app),
	}
	return &add
}

func resetAction(app *app.TrvWorldsInfrastructure) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		paths := []string{
			app.Config.World.WorldsDataPath,
			app.Config.Import.ImportDataPath,
			app.CfgPath,
		}
		for _, path := range paths {
			switch err := os.Remove(path); err {
			case nil:
				fmt.Println("deleted: ", path)
			default:
				fmt.Printf("failed to delete %v:\n  %v\n", path, err)
			}
		}
		return nil
	}
}
