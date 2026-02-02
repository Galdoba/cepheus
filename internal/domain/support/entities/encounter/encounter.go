package encounter

import (
	"fmt"
	"time"

	"github.com/Galdoba/cepheus/pkg/dice"
	"github.com/Galdoba/cepheus/pkg/tttable"
)

type Encounter struct {
	ID           string
	CreatedAt    time.Time
	LastEditedAt time.Time
	Rolls        map[string]TableRollResult
}

type TableRollResult struct {
	TableRoll       int
	TableResult     string
	CanNotBeIgnored bool
}

func New(etype string) (*Encounter, error) {
	now := time.Now()
	id := now.Format("060102150405999")
	e := Encounter{
		ID:           id,
		CreatedAt:    now,
		LastEditedAt: now,
		Rolls:        make(map[string]TableRollResult),
	}
	tc, err := tableCollectionFor(etype)
	if err != nil {
		return nil, fmt.Errorf("failed to setup %v encounter tables: %v", etype, err)
	}
	for _, start := range startPoints(etype) {
		result, err := tc.Roll(start)
		if err != nil {
			return nil, err
		}
		e.Rolls[start] = TableRollResult{
			TableResult: result,
		}
	}
	return &e, nil
}

func encounterTables() []string {
	return []string{
		"Encounter Distance",
	}
}

func tableCollectionFor(etype string) (*tttable.TableCollection, error) {
	tc, err := tttable.NewTableCollection()
	if err != nil {
		return nil, fmt.Errorf("failed to create table collection: %v", err)
	}
	for _, table := range tables(etype) {
		if table == nil {
			return nil, fmt.Errorf("null table detected")
		}
		if err := tc.AddTable(table); err != nil {
			return nil, fmt.Errorf("failed to add table to collection: %v", table.Name)
		}
	}
	tc.SetRoller(dice.New(""))

	return tc, err
}

func startPoints(etype string) []string {
	switch etype {
	default:
		return []string{}
	}
}

func tables(etype string) []*tttable.Table {
	switch etype {
	default:
		return []*tttable.Table{}
	}

}
