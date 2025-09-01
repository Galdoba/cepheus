package coordinates

import (
	"reflect"
	"testing"

	"github.com/Galdoba/cepheus/pkg/grid/coordinates/cube"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/global"
	"github.com/Galdoba/cepheus/pkg/grid/coordinates/local"
)

func TestNewSpaceCoordinates(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
		want SpaceCoordinates
	}{
		{name: "start global", args: args{values: []int{0, 0}}, want: SpaceCoordinates{cube: cube.Cube{Q: 0, R: 0, S: 0}, global: global.SpaceGlobal{X: 0, Y: 0}, local: local.SpaceSectorLocal{SectorX: 0, SectorY: 0, X: 1, Y: 40}}},
		{name: "start cube", args: args{values: []int{0, 0, 0}}, want: SpaceCoordinates{cube: cube.Cube{Q: 0, R: 0, S: 0}, global: global.SpaceGlobal{X: 0, Y: 0}, local: local.SpaceSectorLocal{SectorX: 0, SectorY: 0, X: 1, Y: 40}}},
		{name: "start local", args: args{values: []int{0, 0, 1, 40}}, want: SpaceCoordinates{cube: cube.Cube{Q: 0, R: 0, S: 0}, global: global.SpaceGlobal{X: 0, Y: 0}, local: local.SpaceSectorLocal{SectorX: 0, SectorY: 0, X: 1, Y: 40}}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpaceCoordinates(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v: NewSpaceCoordinates(%v) = %v, want %v", tt.name, tt.args.values, got, tt.want)
			}

		})
	}
}
