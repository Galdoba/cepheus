package check

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/pkg/dice"
)

func TestNew(t *testing.T) {
	dp := dice.New("")
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		opts []CheckOption
		want *Check
	}{
		{
			name: "1",
			opts: []CheckOption{WithCode("2d6")},
			want: &Check{
				checkType:  Raw,
				code:       "2d6",
				dms:        make(map[checkModifier]bool),
				finalDM:    0,
				difficulty: 0,
				effect:     0,
				result:     0,
				resolution: "",
				useBounds:  false,
				lowBound:   0,
				highBound:  0,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.opts...)
			// TODO: update the condition below to compare got with tt.want.
			// if true {
			// 	t.Errorf("New() = %v, want %v", got, tt.want)
			// }
			fmt.Println(got.Resolve(dp))
			fmt.Println(got.Outcome())
		})

	}
}
