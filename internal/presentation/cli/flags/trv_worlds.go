package flags

import "github.com/urfave/cli/v3"

const (
	DRY_RUN = "dry-run"
	RINGS   = "rings"
)

func TrvWorlds_Import() []cli.Flag {
	return []cli.Flag{
		dryRun,
		rings,
	}
}

func TrvWorlds_Survey() []cli.Flag {
	return []cli.Flag{
		dryRun,
		radius,
	}
}

var dryRun = &cli.BoolFlag{
	Name:     "dry-run",
	Category: "",
	Usage:    "do not download data or save data",
	Local:    true,
	Value:    false,
	Aliases:  []string{"d"},
}

var rings = &cli.IntFlag{
	Name:    "rings",
	Usage:   "calibration rings quantity [1-15]",
	Value:   13,
	Aliases: []string{"rn"},
}

var radius = &cli.IntFlag{
	Name:    "radius",
	Usage:   "set radius to survey worlds in",
	Value:   6,
	Aliases: []string{"r"},
}
