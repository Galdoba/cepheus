package subcommand

import (
	"context"
	"fmt"

	"github.com/Galdoba/cepheus/cmd/travellermap/internal/database"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/flags"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/urfave/cli/v3"
)

func CreateMap(actx *infra.Container) *cli.Command {
	return &cli.Command{
		Name:           "create",
		Aliases:        []string{},
		Usage:          "Create new map using canonical data",
		UsageText:      "",
		ArgsUsage:      "",
		Version:        "",
		Description:    "Canonical data will be localized",
		DefaultCommand: "",
		Category:       "",
		Commands:       []*cli.Command{},
		Flags: []cli.Flag{
			flags.Canonical,
		},
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
			source := c.String(flags.CANONICAL)
			if source == "" {
				return fmt.Errorf("empty key provided")
			}
			args := c.Args().Slice()
			for _, arg := range args {
				db, err := database.Create(actx, source, arg)
				switch err {
				default:
					fmt.Printf("%v: database creation failed: %v\n", arg, err)
				case nil:
					fmt.Printf("%v: database created: %v\n", arg, db.Path())
				}
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
		Metadata:                 map[string]interface{}{},
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
