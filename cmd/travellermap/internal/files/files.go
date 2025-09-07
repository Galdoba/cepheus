package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/cepheus/cmd/travellermap/internal/infra"
	"github.com/Galdoba/cepheus/iiss/survey"
)

func AssertCanonicalData(actx *infra.Container) error {
	dir := actx.Config.Files.DataDirectory
	fi, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read canonical data directory: %v", err)
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
		return fmt.Errorf("no canonical files detected")
	}
	for _, file := range canons {
		data, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return fmt.Errorf("failed to read canonical data: %v", err)
		}
		imported := survey.Import{}
		if err := json.Unmarshal(data, &imported); err != nil {
			return fmt.Errorf("failed to unmarshal canonical data: %v", err)
		}
	}

	return nil
}
