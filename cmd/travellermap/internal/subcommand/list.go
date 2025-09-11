package subcommand

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Galdoba/cepheus/cmd/travellermap/internal/database"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/Galdoba/cepheus/pkg/interaction"
	"github.com/urfave/cli/v3"
)

func List(actx *infra.Container) *cli.Command {
	return &cli.Command{
		Name:                            "list",
		Aliases:                         []string{},
		Usage:                           "List world data",
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
			key, err := interaction.GetInput("set key:")
			if err != nil {
				return fmt.Errorf("failed to get user input: %v", err)
			}
			fmt.Println("key:", key)
			db, err := database.Open(filepath.Join(actx.Config.Files.WorkSpaces, "M1105.json"))
			if err != nil {
				return err
			}
			worlds := database.Search(db, key)
			for _, world := range worlds {
				fmt.Println(world.Key())
			}
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
