package gametable

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestNewTable(t *testing.T) {
	tab2, err := NewTable("tab 2", "1d2", NewRollResult("1..2", "second", nil))
	tab, err := NewTable("test", "2d6",
		NewRollResult("0..4", "becomes wormer but undamaged", nil),
		NewRollResult("5..7", "is burned, receive 1D damage", tab2),
		NewRollResult("8", "is burned, receive 6 damage", nil),
		NewRollResult("9+", "Suffers 2D damage", nil),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	dp := dice.NewDicepool()
	res, err := tab.Roll(dp)
	fmt.Println(res, err)
	// tests := []struct {
	// 	name string // description of this test case
	// 	// Named input parameters for target function.
	// 	name    string
	// 	options []*gametable.RollResult
	// 	want    *gametable.GameTable
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, gotErr := gametable.NewTable(tt.name, tt.options)
	// 		if gotErr != nil {
	// 			if !tt.wantErr {
	// 				t.Errorf("NewTable() failed: %v", gotErr)
	// 			}
	// 			return
	// 		}
	// 		if tt.wantErr {
	// 			t.Fatal("NewTable() succeeded unexpectedly")
	// 		}
	// 		// TODO: update the condition below to compare got with tt.want.
	// 		if true {
	// 			t.Errorf("NewTable() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
