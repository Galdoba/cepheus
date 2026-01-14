package flags

import "github.com/urfave/cli/v3"

func TrvWorlds_Import() []cli.Flag {
	return []cli.Flag{
		dryRun,
	}
}

var dryRun = &cli.BoolFlag{
	Name:     "dry-run",
	Category: "",
	Usage:    "do not download data",
	Local:    true,
	Value:    false,
	Aliases:  []string{"d"},
}
