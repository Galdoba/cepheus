package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/cepheus/cmd/travellermap/internal/files"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/Galdoba/cepheus/internal/declare"
	"github.com/urfave/cli/v3"
)

func main() {
	appName := declare.APP_TRAVELLERMAP
	actx, err := infra.Initiate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "init error: %v", err)
		os.Exit(1)
	}
	cmd := cli.Command{
		Name:                            appName,
		Aliases:                         []string{},
		Usage:                           "cli tool to manipulate traveller game maps",
		UsageText:                       "",
		ArgsUsage:                       "",
		Version:                         "",
		Description:                     "",
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
		Before:                          startupCheck(actx),
		After:                           nil,
		Action:                          nil,
		CommandNotFound:                 nil,
		OnUsageError:                    nil,
		InvalidFlagAccessHandler:        nil,
		Hidden:                          false,
		Authors:                         []any{},
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
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v error: %v", appName, err)
		os.Exit(1)
	}
	fmt.Println(actx.Config)
}

func startupCheck(actx *infra.Container) cli.BeforeFunc {
	return func(ctx context.Context, c *cli.Command) (context.Context, error) {
		os.MkdirAll(actx.Config.Files.DataDirectory, 0755)
		os.MkdirAll(actx.Config.Files.WorkSpaces, 0755)
		_, err := files.CanonicalData(actx)
		if err != nil {
			return ctx, err
		}
		return ctx, nil
	}
}
