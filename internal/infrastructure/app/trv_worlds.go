package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/Galdoba/appcontext/configmanager"
	"github.com/Galdoba/cepheus/internal/declare"
	"github.com/Galdoba/cepheus/internal/infrastructure/config"
)

type TrvWorldsInfrastructure struct {
	CfgPath string
	Config  config.TrvWorldsCfg
}

func InitTrvWorlds() (*TrvWorldsInfrastructure, error) {
	appname := declare.APP_TRV_WORLDS
	inf := TrvWorldsInfrastructure{}

	cfg := config.DefaultTrvWorldsConfig()
	cfgman, err := configmanager.New(appname, cfg)
	if err != nil {
		return nil, fmt.Errorf("config init failed: %v", err)
	}
	if wantSaveInitalConfig(cfgman.Path()) {
		fmt.Fprintf(os.Stderr, "init config... ")
		if err := cfgman.Save(); err != nil {
			return nil, fmt.Errorf("failed to save initial config: %v", err)
		}
		fmt.Fprintf(os.Stderr, "ok\n")
	}
	if err := cfgman.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "config loading failed: %v", err)
		os.Exit(1)
	}
	inf.Config = cfgman.Config()
	inf.CfgPath = cfgman.Path()
	return &inf, nil
}

func wantSaveInitalConfig(path string) bool {
	if _, err := os.ReadFile(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true
		}
	}
	return false
}
