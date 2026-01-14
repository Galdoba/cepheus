package convert

import (
	"fmt"
	"testing"
)

func Test_localToGlobal(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		sx    int
		sy    int
		lx    int
		ly    int
		want  int
		want2 int
	}{
		{name: "Reference", sx: 0, sy: 0, lx: 1, ly: 40, want: 0, want2: 0},
		{name: "Capital", sx: 0, sy: 0, lx: 21, ly: 18, want: 20, want2: -22},
		{name: "Regina", sx: -4, sy: -1, lx: 19, ly: 10, want: -110, want2: -70},
		{name: "Tepav", sx: 4, sy: 2, lx: 9, ly: 8, want: 136, want2: 48},
		{name: "Shardi", sx: -1, sy: 0, lx: 32, ly: 36, want: -1, want2: -4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := Local_to_global(tt.sx, tt.sy, tt.lx, tt.ly)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("localToGlobal(X) = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("localToGlobal(Y) = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_localToCube(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		sx    int
		sy    int
		lx    int
		ly    int
		want  int
		want2 int
		want3 int
	}{
		{name: "Reference", sx: 0, sy: 0, lx: 1, ly: 40, want: 0, want2: 0, want3: 0},
		{name: "Azpun", sx: 0, sy: 0, lx: 1, ly: 39, want: 0, want2: -1, want3: 1},
		{name: "Wonderland", sx: -1, sy: 1, lx: 30, ly: 1, want: -3, want2: 3, want3: 0},
		{name: "Shardi", sx: -1, sy: 0, lx: 32, ly: 36, want: -1, want2: -3, want3: 4},
		//sdl;fk; asd dl;kf	}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, got3 := Local_to_cube(tt.sx, tt.sy, tt.lx, tt.ly)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("localToCube(Q) %v = %v, want %v", tt.name, got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("localToCube(R) %v = %v, want %v", tt.name, got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("localToCube(S) %v = %v, want %v", tt.name, got3, tt.want3)
			}
		})
	}
}

func Test_cubeToGlobal(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		q     int
		r     int
		s     int
		want  int
		want2 int
	}{
		{name: "Markun", q: 2, r: -3, s: 1, want: 2, want2: -2}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := Cube_to_global(tt.q, tt.r, tt.s)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("cubeToGlobal(X) %v = %v, want %v", tt.name, got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("cubeToGlobal(Y) %v = %v, want %v", tt.name, got2, tt.want2)
			}
		})
	}
}

func Test_globalToCube(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		x     int
		y     int
		want  int
		want2 int
		want3 int
	}{
		{name: "Markun", x: 2, y: -2, want: 2, want2: -3, want3: 1},
		{name: "Regina", x: -110, y: -70, want: -110, want2: -15, want3: 125},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, got3 := Global_to_cube(tt.x, tt.y)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("globalToCube(Q) %v = %v, want %v", tt.name, got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("globalToCube(R) %v = %v, want %v", tt.name, got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("globalToCube(S) %v = %v, want %v", tt.name, got3, tt.want3)
			}
		})
	}
}

func TestRoundTrips(t *testing.T) {
	rad := 500
	coversions := 0
	for q := -rad; q <= rad; q++ {
		for r := -rad; r <= rad; r++ {
			for s := -rad; s <= rad; s++ {
				if q+r+s != 0 {
					continue
				}
				// if err := roundTrip_Cube(q, r, s); err != nil {
				// 	t.Errorf("%v", err)
				// }

				// if err2 := roundTrip_CubeBack(q, r, s); err2 != nil {
				// 	t.Errorf("%v", err2)
				// }
				if err := RoundTrip(q, r, s); err != nil {
					t.Errorf("%v", err)
				}
				coversions += 6
			}
		}
	}
	fmt.Println("total", coversions, "conversions")
}
