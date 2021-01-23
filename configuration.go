package agora

import (
	"errors"
	"strings"
)

var (
	ErrConfigPathMustBeAbsolute = errors.New("config path must be absolute or current wd")
	ErrConfigMustBeTOML         = errors.New("only TOML configuration files are supported")
)

type Configuration func(*Options)

func ConfigPath(configPath string) Configuration {
	if configPath != "." || !strings.HasPrefix(configPath, "/") {
		panic(ErrConfigPathMustBeAbsolute)
	}

	return func(o *Options) {
		o.configFilePath = configPath
	}
}

func ConfigName(configFileName string) Configuration {
	if !strings.HasSuffix(configFileName, ".toml") {
		panic(ErrConfigMustBeTOML)
	}

	return func(o *Options) {
		o.configFileName = configFileName
	}
}
