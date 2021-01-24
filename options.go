package agora

import (
	"fmt"
	"github.com/wgentry22/agora/modules/api"
	"github.com/wgentry22/agora/types/config"
	"os"
)

type Options struct {
	configFilePath string
	configFileName string
}

func (o *Options) resolvedConfigPath() string {
	return fmt.Sprintf("%s/%s", o.configFilePath, o.configFileName)
}

func defaultOptions() *Options {
	return &Options{
		configFilePath: ".",
		configFileName: "app.toml",
	}
}

func New(confs ...Configuration) Application {
	opts := defaultOptions()

	for _, conf := range confs {
		conf(opts)
	}

	decoder := config.NewTOMLFileParser(opts.resolvedConfigPath())

	appConfig := decoder()

	application := Application{
		errors: make(chan error, 1),
		quit:   make(chan os.Signal, 1),
		conf:   appConfig,
		router: api.NewRouter(appConfig.API()),
	}

	application.Setup()

	return application
}
