package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Galdoba/cepheus/cmd/mgt2-boarding-action/internal/action"
	"github.com/Galdoba/cepheus/internal/declare"
	"github.com/urfave/cli/v3"
)

func main() {
	appName := declare.APP_BOARDING_ACTION

	cmd := cli.Command{
		Name:        declare.APP_BOARDING_ACTION,
		Aliases:     []string{},
		Usage:       "calculate Boarding Action according to rules from MgT2 CRB p.175",
		UsageText:   fmt.Sprintf("%s [global options] [arguments]", appName),
		Version:     "0.0.1",
		Description: "For tabletop game Traveller",

		Action: action.BoardingAction,

		Authors: []any{"galdoba"},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v shutdown with error: %v", cmd.Name, err)
		os.Exit(1)
	}
	os.Exit(0)
}
