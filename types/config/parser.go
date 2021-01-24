package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type Parser func() Application

func NewTOMLFileParser(file string) Parser {
	var conf Application

	data, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	return func() Application {
		if err := toml.Unmarshal(data, &conf); err != nil {
			panic(err)
		}

		return conf
	}
}
