package action

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Galdoba/cepheus/uwp"
	"github.com/Galdoba/cepheus/uwp/descriptor"
	"github.com/urfave/cli/v3"
)

func SetupDescription(ctx context.Context, c *cli.Command) (context.Context, error) {
	path, err := descriptionFilePath(c.Name)
	if err != nil {
		return context.Background(), fmt.Errorf("failed to get description file path: %v", err)
	}
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ds := descriptor.New(path)
			fillDescriptorWithDefaultValues(ds)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return context.Background(), fmt.Errorf("failed to create data directory: %v", err)
			}
			fnew, err := os.Create(path)
			if err != nil {
				return context.Background(), fmt.Errorf("failed to create description file: %v", err)
			}
			defer fnew.Close()
			if err := ds.Save(); err != nil {
				return context.Background(), fmt.Errorf("failed to save descriptor: %v", err)
			}
			return context.Background(), nil
		}
		return context.Background(), fmt.Errorf("unexpected error occured: %v", err)
	}
	defer f.Close()
	return context.Background(), nil
}

func descriptionFilePath(appName string) (string, error) {
	hm, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home space: %v", err)
	}
	return filepath.ToSlash(filepath.Join(hm, ".local", "share", appName, "uwp_descriptions.json")), nil
}

func fillDescriptorWithDefaultValues(ds *descriptor.Descriptor) {
	for _, category := range []string{
		uwp.Port,
		uwp.Size,
		uwp.Atmo,
		uwp.Hydr,
		uwp.Pops,
		uwp.Govr,
		uwp.Laws,
		uwp.TL,
	} {
		codes := []string{}
		switch category {
		case uwp.Port:
			codes = []string{"A", "B", "C", "D", "E", "X", "F", "G", "H", "Y", "?"}
		case uwp.Size:
			codes = []string{"0", "R", "S", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "?"}
		case uwp.Atmo:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "?"}
		case uwp.Hydr:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "?"}
		case uwp.Pops:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "E", "F", "?"}
		case uwp.Govr:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "W", "X", "?"}
		case uwp.Laws:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "S", "?"}
		case uwp.TL:
			codes = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "?"}
		}
		for _, code := range codes {
			ds.AddDescription(category, code, uwp.Description(category, code)["en"])
		}
	}
}
