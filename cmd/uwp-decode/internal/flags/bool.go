package flags

import "github.com/urfave/cli/v3"

var Mapping = cli.BoolFlag{
	Name:    "mapping",
	Usage:   "print codes mapping",
	Aliases: []string{"m"},
}
