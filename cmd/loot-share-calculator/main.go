package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Galdoba/appcontext"
	"github.com/Galdoba/cepheus/internal/actions"
	lscconfig "github.com/Galdoba/cepheus/internal/config/loot-share-calculator"
	"github.com/Galdoba/cepheus/internal/declare"
	cli "github.com/urfave/cli/v3"
)

func main() {
	appName := declare.APP_LOOT_CALCULATOR
	cfg := lscconfig.Default()
	actx := appcontext.New(appName, appcontext.WithConfig(&cfg))
	// appName := "loot-share-calculator"
	// cfg := config.Config{
	// ShipCrew:       11,
	// PlayersPresent: 5,
	// OfficerRatio:   10,
	// MasteryMod:     0,
	// }
	// actx := appcontex.New(appName, appcontex.WithConfig(&cfg))

	cmd := cli.Command{
		Name:                       "loot-share-calculator",
		Aliases:                    []string{},
		Usage:                      "fast calculate loot shares for ship crew in Pirates of Drinax Campaign",
		UsageText:                  "",
		ArgsUsage:                  "arg usage text",
		Version:                    "0.0.1",
		Description:                "For tabletop games like Traveller, Cepheus Deluxe, Hostile, etc.",
		DefaultCommand:             "",
		Category:                   "",
		Commands:                   []*cli.Command{},
		Flags:                      []cli.Flag{},
		HideHelp:                   false,
		HideHelpCommand:            false,
		HideVersion:                false,
		EnableShellCompletion:      false,
		ShellCompletionCommandName: "",

		Action: actions.CalculateLootShares(actx),

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
