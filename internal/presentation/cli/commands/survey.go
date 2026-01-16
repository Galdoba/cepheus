package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/cepheus/internal/domain/worlds/aggregates/world"
	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/t5ss"
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
		Usage:       "show detailed information on selected world",
		UsageText:   "trv_worlds [global options] inspect [options]",
		Description: "Print available information on selected world.\nAdditional data will be calculated and stored based on MOARN principle.",
		Action:      surveyAction(*app.Config),
		Flags:       flags.TrvWorlds_Survey(),
	}
	return &add
}

func surveyAction(cfg config.TrvWorldsCfg) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		worldDB_path = cfg.World.WorldsDataPath
		importDB_path = cfg.Import.ImportDataPath
		//open
		fmt.Printf("open database... ")
		world_db, import_db, err := openDatabases()
		if err != nil {
			fmt.Println("failed!!")
			return err
		}
		fmt.Println("success!")
		//search
		searchKey := c.String("s")
		if searchKey == "" {
			input, err := prompt.Input(
				prompt.WithTitle("input search key:"),
				prompt.WithDescription("coordinates, name or other data..."),
			)
			if err != nil {
				return fmt.Errorf("failed to set search key: %v", err)
			}
			searchKey = input
		}
		worldsMatch := []string{}
		importMatch := []string{}
		for _, key := range world_db.AllKeys() {
			if searchKey == "" || strings.Contains(key, searchKey) {
				worldsMatch = append(worldsMatch, key)
			}
		}
		if len(worldsMatch) == 0 {
			for _, key := range import_db.AllKeys() {
				data, err := import_db.Read(key)
				if err != nil {
					return nil
				}
				if searchKey == "" || strings.Contains(data.Details_DB_Key(), searchKey) {
					importMatch = append(importMatch, key)

				}
			}
		}
		fmt.Println(len(worldsMatch), len(importMatch))
		switch len(worldsMatch) {
		default:
			fmt.Println("data found:")
			items, err := collectWorlds(world_db, worldsMatch)
			if err != nil {
				return err
			}
			selected, err := prompt.SearchItem(
				prompt.FromItems(items),
				prompt.WithTitle("select world:"),
			)
			if err != nil {
				return err
			}
			dto := selected.Payload().(world.WorldDTO)
			fmt.Println(dto)
		case 0:
			fmt.Println("no data found...")
			confirmed, err := prompt.Confirm(
				prompt.WithTitle("search in imported?"),
			)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}

			items, err := collectWorldsImported(import_db, importMatch)
			if err != nil {
				return err
			}
			selected, err := prompt.SearchItem(
				prompt.FromItems(items),
				prompt.WithTitle("select world:"),
			)
			if err != nil {
				return err
			}
			importedData := selected.Payload().(t5ss.WorldData)
			fmt.Println(importedData)
		}
		world_db.CommitAndClose()
		import_db.CommitAndClose()

		return nil
	}
}

var worldDB_path string
var importDB_path string

type worldDB interface {
	AllKeys() []string
	Create(string, world.WorldDTO) error
	Read(string) (world.WorldDTO, error)
	Update(string, world.WorldDTO) error
	CommitAndClose() error
}

type importDB interface {
	AllKeys() []string
	Read(string) (t5ss.WorldData, error)
	CommitAndClose() error
}

func openDatabases() (worldDB, importDB, error) {
	js, err := jsonstorage.OpenStorage[world.WorldDTO](worldDB_path)
	if err != nil {
		switch errors.Is(err, os.ErrNotExist) {
		case true:
			js, err = jsonstorage.NewStorage[world.WorldDTO](worldDB_path)
			if err != nil {
				return nil, nil, err
			}
		case false:
			return nil, nil, fmt.Errorf("failed to create new storage: %v", err)
		}
	}
	js2, err := jsonstorage.OpenStorage[t5ss.WorldData](importDB_path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open import_DB: %v", err)
	}
	return js, js2, nil
}

func collectWorlds(db worldDB, keys []string) ([]*prompt.Item, error) {
	items := []*prompt.Item{}
	for _, key := range keys {
		data, err := db.Read(key)
		if err != nil {
			return nil, err
		}
		itemKey := fmt.Sprintf("%v %v %v", data.Coordinates, data.Name, data.MainworldUWP)
		itemKey = strings.TrimSpace(itemKey)
		items = append(items, prompt.NewItem(itemKey, data))
	}
	return items, nil
}

func collectWorldsImported(db importDB, keys []string) ([]*prompt.Item, error) {
	items := []*prompt.Item{}
	for _, key := range keys {
		data, err := db.Read(key)
		if err != nil {
			return nil, err
		}
		itemKey := fmt.Sprintf("%v", data.Details_DB_Key())
		itemKey = strings.TrimSpace(itemKey)
		items = append(items, prompt.NewItem(itemKey, data))
	}
	return items, nil
}
