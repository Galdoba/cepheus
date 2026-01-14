package app

import (
	"fmt"
	"os"

	"github.com/Galdoba/appcontext/configmanager"
	"github.com/Galdoba/cepheus/internal/declare"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
)

type TrvWorldsInfrastructure struct {
	Config *config.TrvWorldsCfg
}

func InitTrvWorlds() (*TrvWorldsInfrastructure, error) {
	appname := declare.APP_TRV_WORLDS
	inf := TrvWorldsInfrastructure{}

	cfgman, err := configmanager.New(appname, config.DefaultTrvWorldsConfig())
	if err != nil {
		return nil, fmt.Errorf("config init failed: %v", err)
	}
	if err := cfgman.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "config loading failed: %v\nfallback to default configuration...\n", err)
	}
	inf.Config = cfgman.Config()

	return &inf, nil
}
