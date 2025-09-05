package files

import (
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
)

func CanonicalData(actx *infra.Container) ([]string, error) {
	dir := actx.Config.Files.DataDirectory
	fi, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read canonical data direcory: %v", err)
	}
	canons := []string{}
	for _, f := range fi {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		canons = append(canons, f.Name())
	}
	if len(canons) == 0 {
		return nil, fmt.Errorf("no canonical files detected: run 'travellermap update' to download Traveller OTU data")
	}
	return canons, nil
}
