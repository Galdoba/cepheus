package paths

import "github.com/Galdoba/appcontext/xdg"

const (
	appName = "cepheus"
)

func AssetsDirectory(subs ...string) string {
	dirs := []string{"assets"}
	dirs = append(dirs, subs...)
	return xdg.Location(
		xdg.ForData(),
		xdg.WithProgramName(appName),
		xdg.WithSubDir(dirs),
	)
}
