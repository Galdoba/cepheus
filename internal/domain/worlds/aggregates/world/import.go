package world

import (
	"fmt"

	"github.com/Galdoba/cepheus/internal/domain/worlds/valueobject/uwp"
)

func (w *World) importUWP(uwpData string) error {
	if uwpData == "" || w.mainworldUWP != "" {
		return nil
	}
	u, err := uwp.New(uwpData)
	if err != nil {
		fmt.Printf("failed to create uwp from canonical data: %v\n", err)
	}

	w.mainworldUWP = u

	return nil
}
