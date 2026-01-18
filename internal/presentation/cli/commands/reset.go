package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/consolio/prompt"
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
		objectsToDelete := make(map[string]string)
		keys := []string{"Config File", "Derived Database", "External Database"}
		objectsToDelete[keys[0]] = app.CfgPath
		objectsToDelete[keys[1]] = app.Config.World.Derived_DB_File
		objectsToDelete[keys[2]] = app.Config.Import.External_DB_File

		paths := []*prompt.Item{}
		for _, key := range keys {
			paths = append(paths, prompt.NewItem(key, objectsToDelete[key]))
		}
		toBeDeleted, err := prompt.SelectMultiple(
			prompt.WithTitle("select objects to delete:"),
			prompt.FromItems(paths),
		)
		if err != nil {
			return fmt.Errorf("selection failed: %v", err)
		}

		for _, item := range toBeDeleted {
			path := item.Payload().(string)
			switch err := os.Remove(path); err {
			case nil:
				fmt.Fprintf(os.Stderr, "deleted: %v (%v)\n", item.Key(), path)
			default:
				fmt.Fprintf(os.Stderr, "failed to delete %v: %v\n", item.Key(), err)
			}
		}
		return nil
	}
}
