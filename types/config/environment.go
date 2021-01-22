package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Environment int8

const (
	Development Environment = iota
	QualityAssurance
	Staging
	Production
	EnvUnknown
)

var (
	ErrUnknownEnvironment = func(in string) error {
		return fmt.Errorf("unknown environment `%s`", in)
	}
	envDisplay = []string{"dev", "qa", "staging", "prod", "unknown"}
	envLookup  = map[string]Environment{
		"dev":     Development,
		"qa":      QualityAssurance,
		"staging": Staging,
		"prod":    Production,
		"unknown": EnvUnknown,
	}
)

func (e Environment) String() string {
	return envDisplay[e]
}

func (e Environment) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func ParseEnvironment(in string) (Environment, error) {
	e, ok := envLookup[in]
	if ok {
		return e, nil
	}

	return EnvUnknown, ErrUnknownEnvironment(in)
}

func (e *Environment) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == "!!str" {
		env, err := ParseEnvironment(node.Value)
		if err != nil {
			return err
		}

		*e = env

		return nil
	}

	return errors.New("[yaml] expected environment to be a string")
}

func (e *Environment) UnmarshalText(data []byte) error {
	env, err := ParseEnvironment(string(data))

	if err != nil {
		return err
	}

	*e = env

	return nil
}
