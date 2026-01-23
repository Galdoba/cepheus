package commands

import (
	"context"
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/aggregates/world"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
	"github.com/Galdoba/cepheus/internal/infrastructure/jsonstorage"
	"github.com/Galdoba/cepheus/internal/presentation/cli/flags"
	"github.com/Galdoba/consolio/prompt"
	"github.com/urfave/cli/v3"
)

const (
	INSPECT_WORLD_COMMAND = "inspect"
)

func InspectWorld(app *app.TrvWorldsInfrastructure) *cli.Command {
	add := cli.Command{
		Name:        INSPECT_WORLD_COMMAND,
		Aliases:     []string{},
		Usage:       "show detailed information on selected world",
		UsageText:   "trv_worlds [global options] inspect [options]",
		Description: "Print available information on selected world.",
		Action:      inspectAction(*app.Config),
		Flags:       flags.TrvWorlds_Survey(),
	}
	return &add
}

func inspectAction(cfg config.TrvWorldsCfg) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {

		derived_db_path := cfg.World.Derived_DB_File
		// //open
		fmt.Printf("open database... ")
		workingDataStorage, err := jsonstorage.OpenOrCreateStorage[world.WorldDTO](derived_db_path)
		if err != nil {
			return fmt.Errorf("failed open working data storage: %v", err)
		}

		// keyMap := make(map[string]string)
		keys := []*prompt.Item{}

		for _, entry := range workingDataStorage.AllEntries() {
			inspectKey := entry.Imported.SearchKey()
			keys = append(keys, prompt.NewItem(inspectKey, entry))
		}
		if len(keys) <= 1 {
			fmt.Println("working database contais", len(keys), "entries")
			fmt.Println("survey more worlds!!")
		}
		selected, err := prompt.SearchItem(
			prompt.FromItems(keys),
			prompt.WithTitle("search world:"),
		)
		if err != nil {
			return err
		}
		dto := selected.Payload().(world.WorldDTO)
		w, err := world.FromDTO(dto.DatabaseKey(), dto)
		if err != nil {
			return err
		}
		fmt.Println("returned object:")
		fmt.Println(w)
		return nil
	}
}
