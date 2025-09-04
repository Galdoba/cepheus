package main

import (
	"fmt"

	"github.com/Galdoba/appcontext/configmanager"
	"github.com/Galdoba/appcontext/logmanager"
	"github.com/Galdoba/appcontext/logmanager/colorizer"
	"github.com/Galdoba/appcontext/xdg"
	"github.com/Galdoba/cepheus/cmd/travellermap/internal/config"
	"github.com/Galdoba/cepheus/internal/declare"
)

type AppContext struct {
	Config    *config.Config
	Logger    *logmanager.Logger
	PathMaker *xdg.ProgramPaths
}

func Initiate() (*AppContext, error) {
	app := declare.APP_TRAVELLERMAP
	pathman := xdg.New(app)
	actx := AppContext{}
	cm, err := configmanager.New(app,
		config.Config{
			App: config.AppData{
				Version: "",
			},
			Logger: config.Logger{
				Enabled:       false,
				Filepath:      "",
				Level:         "",
				ConsoleOutput: false,
				Color:         false,
			},
			Files: config.Files{
				DataDirectory: pathman.PersistentDataDir(),
				WorkSpaces:    pathman.StateDir(),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup config")
	}
	if err := cm.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config")
	}
	actx.Config = cm.Config()
	logHandlers := []*logmanager.MessageHandler{}
	logOpts := actx.Config.Logger
	if logOpts.ConsoleOutput {
		logHandlers = append(logHandlers, logmanager.NewHandler(logmanager.Stderr, logmanager.StringToLevel(logOpts.Level), logmanager.NewTextFormatter(
			logmanager.WithColor(logOpts.Color), logmanager.WithLevelTag(true), logmanager.WithTimePrecision(3),
		)))
	}
	logHandlers = append(logHandlers, logmanager.NewHandler(pathman.LogFile(), logmanager.StringToLevel(logOpts.Level), logmanager.NewTextFormatter(
		logmanager.WithLevelTag(true), logmanager.WithTimePrecision(3),
	)))
	colorizer.DefaultColorizer()
	logger := logmanager.New(logmanager.WithHandlers(logHandlers...))
	actx.Logger = logger
	actx.PathMaker = pathman
	return &actx, nil
}
