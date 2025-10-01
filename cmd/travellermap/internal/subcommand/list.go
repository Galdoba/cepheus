package subcommand

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/cepheus/cmd/travellermap/internal/database"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/Galdoba/consolio/prompt"
	"github.com/urfave/cli/v3"
)

func List(actx *infra.Container) *cli.Command {
	return &cli.Command{
		Name:                            "search",
		Aliases:                         []string{},
		Usage:                           "Search world data",
		UsageText:                       "",
		ArgsUsage:                       "",
		Version:                         "",
		Description:                     "Read map file, search worlds by key provided and print data",
		DefaultCommand:                  "",
		Category:                        "",
		Commands:                        []*cli.Command{},
		Flags:                           []cli.Flag{},
		HideHelp:                        false,
		HideHelpCommand:                 false,
		HideVersion:                     false,
		EnableShellCompletion:           false,
		ShellCompletionCommandName:      "",
		ShellComplete:                   nil,
		ConfigureShellCompletionCommand: nil,
		Before:                          nil,
		After:                           nil,
		Action: func(ctx context.Context, c *cli.Command) error {
			db, err := database.Open(filepath.Join(actx.Config.Files.WorkSpaces, "M1105.json"))
			if err != nil {
				return err
			}
			worldMap, err := db.GetAll()
			if err != nil {
				return fmt.Errorf("failed to read worlds database")
			}
			items := []*prompt.Item{}
			for key := range worldMap {
				items = append(items, prompt.CreateItem(key))
			}
			found, err := prompt.SearchItem(
				prompt.WithTitle("search world"),
				prompt.WithDescription("enter world name or UWP"),
				prompt.FromItems(items...),
			)
			if err != nil {
				return err
			}
			world, err := db.Get(found.GetKey())
			if err != nil {
				return fmt.Errorf("failed to get world: %v", err)
			}
			fmt.Fprintf(os.Stderr, "%v\n", world)

			return nil
		},
		CommandNotFound:          nil,
		OnUsageError:             nil,
		InvalidFlagAccessHandler: nil,
		Hidden:                   false,
		Authors:                  []any{},
		Copyright:                "",
		Reader:                   nil,
		Writer:                   nil,
		ErrWriter:                nil,
		ExitErrHandler:           nil,
		Metadata:                 map[string]any{},
		ExtraInfo: func() map[string]string {
			panic("TODO")
		},
		CustomRootCommandHelpTemplate: "",
		SliceFlagSeparator:            "",
		DisableSliceFlagSeparator:     false,
		UseShortOptionHandling:        false,
		Suggest:                       false,
		AllowExtFlags:                 false,
		SkipFlagParsing:               false,
		CustomHelpTemplate:            "",
		PrefixMatchCommands:           false,
		SuggestCommandFunc:            nil,
		MutuallyExclusiveFlags:        []cli.MutuallyExclusiveFlags{},
		Arguments:                     []cli.Argument{},
		ReadArgsFromStdin:             false,
	}
}
