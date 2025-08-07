package dice

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func Test_parseGoal(t *testing.T) {
	fmt.Println("Test_parseGoal:")
	type args struct {
		goal string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		//Test cases.
		{
			name: "simple",
			args: args{
				goal: "3",
			},
			want: []int{3},
		},
		{
			name: "multiple",
			args: args{
				goal: "3,5",
			},
			want: []int{3, 5},
		},
		{
			name: "tn and more positive",
			args: args{
				goal: "5+",
			},
			want: andMoreResult(5),
		},
		{
			name: "tn and more negative",
			args: args{
				goal: "-5+",
			},
			want: andMoreResult(-5),
		},
		{
			name: "tn and less positive",
			args: args{
				goal: "2-",
			},
			want: andLessResult(2),
		},
		{
			name: "tn and less negative",
			args: args{
				goal: "-2-",
			},
			want: andLessResult(-2),
		},
		{
			name: "tn from +X to +Y",
			args: args{
				goal: "2...5",
			},
			want: fromToResult(2, 5),
		},
		{
			name: "tn from -X to +Y",
			args: args{
				goal: "-2...2",
			},
			want: fromToResult(-2, 2),
		},
		{
			name: "tn from -X to -Y",
			args: args{
				goal: "-5...-2",
			},
			want: fromToResult(-5, -2),
		},
		{
			name: "complex",
			args: args{
				goal: "2,3,6...8",
			},
			want: []int{2, 3, 6, 7, 8},
		},
		{
			name: "NaN",
			args: args{
				goal: "a",
			},
			want: []int{},
		},
		{
			name: "repeating",
			args: args{
				goal: "2,3,2",
			},
			want: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseGoal(tt.args.goal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGoal() = %v, want %v", got, tt.want)
			} else {
				fmt.Printf("pass: '%v'\n     arg: '%v'\n  result: '%v'\n", tt.name, tt.args.goal, got)
			}
		})
	}

}

func andMoreResult(i int) []int {
	out := []int{}
	for v := i; v <= 30; v++ {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

func andLessResult(i int) []int {
	out := []int{}
	for v := i; v >= -30; v-- {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

func fromToResult(i, j int) []int {
	out := []int{}
	for v := i; v <= j; v++ {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}
