package rtg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/cepheus/internal/infrastructure/filepaths"
	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/tttable"
)

const (
	starSystemSubfolder = "star_system"

	MOD_NonMainSequenceClass = "non main sequence class determined"

	TableVariant_Type_Realistic = "realistic types"
)

func NewStarTypeDeterminationGenerator(seed string, opts ...string) (RandomTableGenerator, error) {
	options := make(map[string]bool)
	for _, o := range opts {
		options[o] = true
	}
	tables := []*tttable.Table{}
	for name, path := range StarSystemTablesMap() {
		table, err := tttable.Load(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load table %v: %v", name, err)
		}
		switch name {
		case "Type":
			if options[TableVariant_Type_Realistic] {
			}
		}
		tables = append(tables, table)
	}
	tc, err := tttable.NewTableCollection(
		tttable.WithRoller(dice.New(seed)),
		tttable.WithTables(tables...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create random tables colection: %v", err)
	}

	generator := randomTableGenerator{}
	generator.tableCollection = tc
	return &generator, nil
}

func starTypeDeterminationTableNames() []string {
	return []string{
		"Type",
		"Hot",
		"Special",
		"Unusual",
		"Giants",
		"Peculiar",
	}
}

func StarSystemTablesMap() map[string]string {
	tableMap := make(map[string]string)
	for _, name := range starTypeDeterminationTableNames() {
		tableMap[name] = filepath.Join(filepaths.RandomTablesDirectory(), starSystemSubfolder, tableFileName(name))
	}
	return tableMap
}

func InitStarSystemDeterminationTables() error {
	for name, path := range StarSystemTablesMap() {
		err := AssertTable(path)
		if errors.Is(err, os.ErrNotExist) {
			switch name {
			case "Type":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"2-":   "Special/Unusual",
						"3-6":  "M",
						"7-8":  "K",
						"9-10": "G",
						"11":   "F",
						"12+":  "Hot",
					},
					map[string]int{MOD_NonMainSequenceClass: 1},
				)
			case "Hot":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"9-":    "A",
						"10-11": "B",
						"12+":   "O",
					},
					map[string]int{},
				)
			case "Special":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"5-":   "VI",
						"6-8":  "IV",
						"9-10": "III",
						"11+":  "Giants",
					},
					map[string]int{},
				)
			case "Unusual":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"2-":   "Peculiar",
						"3":    "VI",
						"4":    "IV",
						"5-7":  "BD",
						"8-10": "D",
						"11":   "III",
						"12+":  "Giants",
					},
					map[string]int{},
				)
			case "Giants":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"8-":   "III",
						"9-10": "II",
						"11":   "Ib",
						"12+":  "Ia",
					},
					map[string]int{},
				)
			case "Peculiar":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"2-":  "BH",
						"3":   "PSR",
						"4":   "NS",
						"5-6": "Nb",
						"7-9": "Protostar",
						"10":  "Star Cluster",
						"11+": "Anomaly",
					},
					map[string]int{},
				)
			}
			if err != nil {
				return fmt.Errorf("failed to initiate table: %v", err)
			}

		}
	}
	return nil
}

// 	HotTable, err := tttable.NewTable("Hot",
// 		tttable.WithDiceExpression("2d6"),
// 		tttable.WithRows(
// 			tttable.NewRow("9-", "A"),
// 			tttable.NewRow(tttable.MustIndex(10, 11), "B"),
// 			tttable.NewRow("12+", "O"),
// 		),
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	UnusualTable, err := tttable.NewTable("Unusual",
// 		tttable.WithDiceExpression("2d6"),
// 		tttable.WithRows(
// 			tttable.NewRow("2-", "Peculiar"),
// 			tttable.NewRow("3", "VI"),
// 			tttable.NewRow("4", "IV"),
// 			tttable.NewRow("5-7", "BD"),
// 			tttable.NewRow("8-10", "D"),
// 			tttable.NewRow("11", "III"),
// 			tttable.NewRow("12+", "Giants"),
// 		),
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	PeculiarTable, err := tttable.NewTable("Peculiar",
// 		tttable.WithDiceExpression("2d6"),
// 		tttable.WithRows(
// 			tttable.NewRow("2-", "BH"),
// 			tttable.NewRow("3", "PSR"),
// 			tttable.NewRow("4", "NS"),
// 			tttable.NewRow("5-6", "NB"),
// 			tttable.NewRow("7-9", "Protostar"),
// 			tttable.NewRow("10", "Star Cluster"),
// 			tttable.NewRow("11+", "Anomaly"),
// 		),
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	GiantsTable, _ := tttable.NewTable("Giants",
// 		tttable.WithDiceExpression("2d6"),
// 		tttable.WithRows(
// 			tttable.NewRow(tttable.MustIndex(8, tttable.AndLess), "III"),
// 			tttable.NewRow(tttable.MustIndex(9, 10), "II"),
// 			tttable.NewRow(tttable.MustIndex(11), "Ib"),
// 			tttable.NewRow(tttable.MustIndex(12, tttable.AndMore), "Ia"),
// 		),
// 	)
// }
