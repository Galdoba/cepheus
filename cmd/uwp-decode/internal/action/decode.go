package action

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/appcontext"
	"github.com/Galdoba/cepheus/pkg/uwp"
	"github.com/Galdoba/cepheus/pkg/uwp/descriptor"
	"github.com/urfave/cli/v3"
)

func Decode(ctx context.Context, c *cli.Command) error {
	//Init
	// context.
	path, err := descriptionFilePath(c.Name)
	if err != nil {
		return fmt.Errorf("failed to get description file path: %v", err)
	}
	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	return fmt.Errorf("failed to read description file: %v", err)
	// }
	ds := descriptor.New(path)
	if err := ds.Load(); err != nil {
		return fmt.Errorf("failed to load descriptor: %v", err)
	}

	args := c.Args().Slice()
	//Process
	all := len(args)
	if all == 0 {
		fmt.Fprintf(os.Stderr, "no arguments provided\ntype '%v -h' for program description\n", c.Name)
		os.Exit(0)
	}
	dataTypes := dataTypes()
	lang := c.String("language")
	for i, uwpcode := range args {
		description := make(map[string]string)
		printHeader(i, all)
		codes := strings.Split(uwpcode, "")
		if len(codes) != len(dataTypes) {
			fmt.Fprintf(os.Stderr, "error: bad argument '%v': expect string with lenght of 9 characters\n", uwpcode)
			continue
		}
		for i, category := range dataTypes {
			if category == "" {
				continue
			}
			description[category] = ds.Get(category, strings.ToUpper(codes[i]), lang)
		}
		fmt.Fprintf(os.Stdout, "UWP code '%v' description:\n", uwpcode)
		printMapping(c.Bool("mapping"))
		for i, category := range dataTypes {
			if i == 7 {
				continue
			}
			fmt.Fprintf(os.Stdout, "\n%v:\n", strings.ToUpper(category))
			fmt.Fprintf(os.Stdout, " -%v\n", description[category])
		}
		printFooter(i, all)
	}

	return nil
}

func printHeader(i, all int) {
	if all < 2 {
		return
	}
	fmt.Fprintf(os.Stdout, "===process %v of %v===========\n", i+1, all)
}

func printMapping(print bool) {
	if !print {
		return
	}
	//UWP code 'c555555-5'
	//          |
	fmt.Fprintf(os.Stdout, "          ||||||| |\n")
	fmt.Fprintf(os.Stdout, "          ||||||| +----- Technological Level\n")
	fmt.Fprintf(os.Stdout, "          ||||||+------- Law Level\n")
	fmt.Fprintf(os.Stdout, "          |||||+-------- Government Type\n")
	fmt.Fprintf(os.Stdout, "          ||||+--------- Population\n")
	fmt.Fprintf(os.Stdout, "          |||+---------- Hydrosphere\n")
	fmt.Fprintf(os.Stdout, "          ||+----------- Atmosphere\n")
	fmt.Fprintf(os.Stdout, "          |+------------ Size\n")
	fmt.Fprintf(os.Stdout, "          +------------- Spaceport\n")
}

func printFooter(i, all int) {
	if all < 2 {
		return
	}
	fmt.Fprintf(os.Stdout, "\n \n")
}

func dataTypes() []string {
	return []string{
		uwp.Port,
		uwp.Size,
		uwp.Atmo,
		uwp.Hydr,
		uwp.Pops,
		uwp.Govr,
		uwp.Laws,
		"",
		uwp.TL,
	}
}

// type Action func(context.Context, *Command) error
