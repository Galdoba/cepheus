package au

import (
	"testing"
)

func TestToOrbitNumber(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		auInput AU
		want    float64
	}{
		{
			name:    "book test",
			auInput: 3.4,
			want:    5.25,
		}, // TODO: Add test cases.
		{
			name:    "book test",
			auInput: 338,
			want:    12.1,
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.auInput.OrbitNumber()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("ToOrbitNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
