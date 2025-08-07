package flags

import (
	"github.com/urfave/cli/v3"
)

var Language = cli.StringFlag{
	Name:        "language",
	DefaultText: "en",
	Usage:       "set output language",
	Value:       "en",
	Aliases:     []string{"l"},
}
