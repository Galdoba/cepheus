package travellermap

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	// tests := []struct {
	// 	name string // description of this test case
	// 	// Named input parameters for target function.
	// 	url     string
	// 	want    []byte
	// 	wantErr bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, gotErr := travellermap.Get(tt.url)
	// 		if gotErr != nil {
	// 			if !tt.wantErr {
	// 				t.Errorf("Get() failed: %v", gotErr)
	// 			}
	// 			return
	// 		}
	// 		if tt.wantErr {
	// 			t.Fatal("Get() succeeded unexpectedly")
	// 		}
	// 		// TODO: update the condition below to compare got with tt.want.
	// 		if true {
	// 			t.Errorf("Get() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
	fmt.Println("start test")
	data, err := Get("")
	fmt.Println(string(data), err)
	fmt.Println(len(data), "bytes received")
}
