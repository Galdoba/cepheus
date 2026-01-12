package jsonstorage

import (
	"fmt"
	"testing"
)

type TestObject struct {
	A string             `json:"a"`
	B int                `json:"b"`
	C bool               `json:"c"`
	D TestInternalObject `json:"d"`
}

type TestInternalObject struct {
	E int64   `json:"e"`
	F float64 `json:"f"`
}

func TestStorage(t *testing.T) {
	s, err := NewStorage[*TestObject](`c:\Users\pemaltynov\go\src\github.com\Galdoba\cepheus\internal\infrastructure\jsonstorage\database.json`)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create storage: %v", err))
		return
	}
	fmt.Println(s.Create("aa", &TestObject{
		A: "aaa",
		B: 1,
		C: false,
		D: TestInternalObject{
			E: -5,
			F: 3.14,
		},
	}))
	s.Create("bb", &TestObject{
		A: "abb",
		B: 2,
		C: true,
		D: TestInternalObject{
			E: 15,
			F: 153.14,
		},
	})
	s.Commit()
	fmt.Println(s.Read("bb"))
	s.Update("bb", &TestObject{
		A: "abb-updated",
		B: 100,
		C: false,
		D: TestInternalObject{},
	})
	fmt.Println(s.Read("bb"))
	s.Discard("bb")
	fmt.Println(s.Read("bb"))

}
