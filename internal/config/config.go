package config

import (
	// lscconfig "github.com/Galdoba/cepheus/internal/config/loot-share-calculator"
	// "github.com/Galdoba/cepheus/internal/declare"
	"github.com/Galdoba/gogacon"
)

type configProvider struct {
	app string
	Cfg gogacon.Serializer
}

type ConfigProvider interface {
	Provide() gogacon.Serializer
}

// func GetDefault(app string) gogacon.Serializer {
// 	switch app {
// 	case declare.APP_LOOT_CALCULATOR:
// 		return lscconfig.Default()
// 	}
// 	return nil
// }
