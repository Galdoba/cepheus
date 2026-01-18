package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/internal/declare"
	"github.com/Galdoba/cepheus/internal/domain/core/actions"
	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/presentation/cli/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	appname := declare.APP_TRV_WORLDS
	options, err := app.InitTrvWorlds()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start program: %v", err)
		os.Exit(1)
	}
	c := cli.Command{
		Name:           appname,
		Aliases:        []string{},
		Usage:          "",
		UsageText:      "",
		ArgsUsage:      "",
		Version:        "",
		Description:    "",
		DefaultCommand: "",
		Category:       "",
		Commands: []*cli.Command{
			commands.ImportWorldData(options),
			commands.InspectWorld(options),
			commands.Reset(options),
			commands.SurveyWorld(options),
		},
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
		Action:                          actions.TrvWorlds_Status(options),
		CommandNotFound:                 nil,
		OnUsageError:                    nil,
		InvalidFlagAccessHandler:        nil,
		Hidden:                          false,
		Authors:                         []any{"galdoba"},
		Copyright:                       "",
		Reader:                          nil,
		Writer:                          nil,
		ErrWriter:                       nil,
		ExitErrHandler:                  nil,
		Metadata:                        map[string]interface{}{},
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
	if err := c.Run(context.TODO(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "program error: %v\n", err)
		os.Exit(1)
	}

}

/*
logger
config
cli
tui

COMMANDS
import
port
starmap
generate

*/
