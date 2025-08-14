package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/cmd/uwp-decode/internal/action"
	"github.com/Galdoba/cepheus/cmd/uwp-decode/internal/flags"
	"github.com/Galdoba/cepheus/internal/declare"
	"github.com/urfave/cli/v3"
)

func main() {
	appName := declare.APP_UWP_DECODE
	cmd := cli.Command{
		Name:           appName,
		Aliases:        []string{},
		Usage:          "decode Cepheus engine UWP string to human readable text",
		UsageText:      fmt.Sprintf("%s [global options] [arguments]", appName),
		ArgsUsage:      "arg usage text",
		Version:        "0.1.0",
		Description:    "For tabletop games like Traveller, Cepheus Deluxe, Hostile, etc.",
		DefaultCommand: "",
		Category:       "",
		Commands:       []*cli.Command{},
		Flags: []cli.Flag{
			&flags.Language,
			&flags.Mapping,
		},
		HideHelp:                   false,
		HideHelpCommand:            false,
		HideVersion:                false,
		EnableShellCompletion:      false,
		ShellCompletionCommandName: "",

		Before: action.SetupDescription,
		Action: action.Decode,

		Hidden:    false,
		Authors:   []any{"galdoba"},
		Copyright: "",
		Reader:    nil,
		Writer:    nil,
		ErrWriter: nil,
		ExitErrHandler: func(context.Context, *cli.Command, error) {
		},
		Metadata:                      map[string]interface{}{},
		CustomRootCommandHelpTemplate: "",
		SliceFlagSeparator:            "",
		DisableSliceFlagSeparator:     false,
		UseShortOptionHandling:        false,
		Suggest:                       false,
		AllowExtFlags:                 false,
		SkipFlagParsing:               false,
		CustomHelpTemplate:            "",
		PrefixMatchCommands:           false,
		MutuallyExclusiveFlags:        []cli.MutuallyExclusiveFlags{},
		Arguments:                     []cli.Argument{},
		ReadArgsFromStdin:             false,
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v shutdown with error: %v", cmd.Name, err)
		os.Exit(1)
	}
	os.Exit(0)
}
