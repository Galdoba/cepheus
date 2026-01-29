package main

import (
	"fmt"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/tttable"
)

func main() {
	count := 0
	TypeTable, err := tttable.NewTable("Type",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow("2-", "Unusual"),
			tttable.NewRow(tttable.MustIndex(3, 4, 5, 6), "M"),
			tttable.NewRow(tttable.MustIndex(7, 8), "K"),
			tttable.NewRow(tttable.MustIndex(9, 10), "G"),
			tttable.NewRow("11", "F"),
			tttable.NewRow("12+", "Hot"),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	HotTable, err := tttable.NewTable("Hot",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow("9-", "A"),
			tttable.NewRow(tttable.MustIndex(10, 11), "B"),
			tttable.NewRow("12+", "O"),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	UnusualTable, err := tttable.NewTable("Unusual",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow("2-", "Peculiar"),
			tttable.NewRow("3", "VI"),
			tttable.NewRow("4", "IV"),
			tttable.NewRow("5-7", "BD"),
			tttable.NewRow("8-10", "D"),
			tttable.NewRow("11", "III"),
			tttable.NewRow("12+", "Giants"),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	PeculiarTable, err := tttable.NewTable("Peculiar",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow("2-", "BH"),
			tttable.NewRow("3", "PSR"),
			tttable.NewRow("4", "NS"),
			tttable.NewRow("5-6", "NB"),
			tttable.NewRow("7-9", "Protostar"),
			tttable.NewRow("10", "Star Cluster"),
			tttable.NewRow("11+", "Anomaly"),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	GiantsTable, _ := tttable.NewTable("Giants",
		tttable.WithDiceExpression("2d6"),
		tttable.WithRows(
			tttable.NewRow(tttable.MustIndex(8, tttable.AndLess), "III"),
			tttable.NewRow(tttable.MustIndex(9, 10), "II"),
			tttable.NewRow(tttable.MustIndex(11), "Ib"),
			tttable.NewRow(tttable.MustIndex(12, tttable.AndMore), "Ia"),
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	gen := generator{}
	gen.rng = dice.New("")

	tables, err := tttable.NewTableCollection(
		tttable.WithCollectionRoller(gen.rng),
		tttable.WithTables(TypeTable, HotTable, UnusualTable, PeculiarTable, GiantsTable),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for !gen.star.completed {
		res, err := tables.Roll("Type")
		if err != nil {
			panic(err)
		}
		switch res {
		case "O", "B", "A", "F", "G", "K", "M":
			if gen.star.Type == "" {
				gen.star.Type = res
			}
			if gen.star.Class == "" {
				gen.star.Class = "V"
			}
		case "Ia", "Ib", "II", "III", "IV", "VI":
			if gen.star.Class == "" {
				gen.star.Class = res
			}
		case "Anomaly":
			gen.star.Type = res
			gen.star.completed = true
		case "BH", "NS", "D":
			gen.star.Type = res
			gen.star.EmptySystem = true
			gen.star.DeadSystem = true
			gen.star.completed = true
		case "BD":
			gen.star.Type = res
			gen.star.EmptySystem = true
			gen.star.completed = true
		case "PSR":
			gen.star.Type = res
			gen.star.DeadSystem = true
			gen.star.completed = true
		case "NB":
			gen.star.Nebula++
		case "Protostar":
			gen.star.PrimordialSystem = true
		case "Star Cluster":
			gen.star.Cluster = true
		default:
			panic(res)
		}
		count++
		if gen.star.Class != "" && gen.star.Type != "" {
			gen.star.completed = true
		}

	}
	fmt.Println(gen.star, count)
	fmt.Println("Anomaly", gen.star.Anomaly)
	fmt.Println("DeadSystem", gen.star.DeadSystem)
	fmt.Println("EmptySystem", gen.star.EmptySystem)
	fmt.Println("PrimordialSystem", gen.star.PrimordialSystem)
	fmt.Println("Cluster", gen.star.Cluster)
	fmt.Println("Nebula", gen.star.Nebula)
}

type generator struct {
	rng             *dice.Roller
	moreColdStars   bool //change Type table: index 11-
	dwarfsAsPrimary bool //change Type table: index 12
	simplePeculiar  bool //change Type table: index 12

	star Star
}

type Star struct {
	completed        bool
	Type             string
	SubType          int
	Class            string
	Nebula           int
	ProtoStar        bool
	Cluster          bool
	Anomaly          bool
	DeadSystem       bool
	EmptySystem      bool
	PrimordialSystem bool
}
