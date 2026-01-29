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
	MOD_ProtostarSystem      = "protostar system determined class determined"

	TableVariant_Type_Realistic = "realistic types"
)

func NewStarTypeDeterminationGenerator(roller *dice.Roller, opts ...string) (RandomTableGenerator, error) {
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
		tables = append(tables, table)
	}
	tc, err := tttable.NewTableCollection(
		tttable.WithTables(tables...),
	)
	tc.SetRoller(roller)
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

func starSybtypeDeterminationTableNames() []string {
	return []string{
		"Numeric",
		"M Type Primary",
	}
}

func StarSystemTablesMap() map[string]string {
	tableMap := make(map[string]string)
	for _, name := range starTypeDeterminationTableNames() {
		tableMap[name] = filepath.Join(filepaths.RandomTablesDirectory(), starSystemSubfolder, tableFileName(name))
	}
	for _, name := range starSybtypeDeterminationTableNames() {
		tableMap[name] = filepath.Join(filepaths.RandomTablesDirectory(), starSystemSubfolder, tableFileName(name))
	}
	return tableMap
}

func InitStarSystemDeterminationTables() error {
	for name, path := range StarSystemTablesMap() {
		dir := filepath.Dir(path)
		fmt.Fprintf(os.Stderr, "init table: %v                \r", name)
		err := AssertTable(path)
		if err == nil {
			continue
		}
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0775); err != nil {
				return fmt.Errorf("failed to create random tables directory: %v", err)
			}
			path = dir
			switch name {
			case "Type":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"2-":   "Unusual",
						"3-6":  "M",
						"7-8":  "K",
						"9-10": "G",
						"11":   "F",
						"12+":  "Hot",
					},
					map[string]int{
						MOD_NonMainSequenceClass: 1,
						MOD_ProtostarSystem:      1,
					},
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
			case "Numeric":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"2":  "0",
						"3":  "1",
						"4":  "3",
						"5":  "5",
						"6":  "7",
						"7":  "9",
						"8":  "8",
						"9":  "6",
						"10": "4",
						"11": "2",
						"12": "0",
					},
					map[string]int{},
				)
			case "M Type Primary":
				err = CreateRandomTable(path, name, "2d6",
					map[string]string{
						"2":  "8",
						"3":  "6",
						"4":  "5",
						"5":  "4",
						"6":  "0",
						"7":  "2",
						"8":  "1",
						"9":  "3",
						"10": "5",
						"11": "7",
						"12": "9",
					},
					map[string]int{},
				)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "\n")
				return fmt.Errorf("failed to initiate table: %v", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "\n")
			return fmt.Errorf("failed to initiate table %v: %v", name, err)
		}
	}
	return nil
}
