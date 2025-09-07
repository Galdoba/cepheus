package flags

import "github.com/urfave/cli/v3"

const (
	CANONICAL = "canonical"
)

var Canonical = &cli.StringFlag{
	Name:     CANONICAL,
	Usage:    "set canonical file key (Required)",
	Required: true,
	OnlyOnce: true,
}
