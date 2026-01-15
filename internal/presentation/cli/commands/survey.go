package commands

import (
	"context"
	"fmt"

	"github.com/Galdoba/cepheus/internal/infrastructure/app"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
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
	}
	return &add
}

func surveyAction(cfg config.TrvWorldsCfg) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		fmt.Println(cfg)
		fmt.Println("me survey command!")
		return nil
	}
}
