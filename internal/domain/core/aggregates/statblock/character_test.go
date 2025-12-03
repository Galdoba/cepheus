package statblock

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStatblockFull(t *testing.T) {
	data, err := json.MarshalIndent(Leeroy, "", "  ")

	fmt.Println(err)
	fmt.Println(string(data))
	fmt.Println(Leeroy.View())

}
