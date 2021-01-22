package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

var (
	ErrSemanticVersionFormat = errors.New("expected SemanticVersion to be in format x.x.x")
	ErrPartCannotBe          = func(part, in string) error {
		return fmt.Errorf("cannot use `%s` as `%s`", in, part)
	}
	semanticVersionPartsLength = 3
)

type SemanticVersion struct {
	Major int `json:"major" yaml:"major" toml:"major"`
	Minor int `json:"minor" yaml:"minor" toml:"minor"`
	Patch int `json:"patch" yaml:"patch" toml:"patch"`
}

func (s SemanticVersion) IsZero() bool {
	return s.Major == 0 && s.Minor == 0 && s.Patch == 0
}

func (s SemanticVersion) IsStrictlyLessThan(other SemanticVersion) bool {
	if s.Major < other.Major {
		return true
	} else if s.Major == other.Major && s.Minor < other.Minor {
		return true
	} else if s.Major == other.Major && s.Minor == other.Minor {
		return s.Patch < other.Patch
	}

	return false
}

func ParseSemanticVersion(in string) (*SemanticVersion, error) {
	parts := strings.Split(in, ".")
	if len(parts) != semanticVersionPartsLength {
		return nil, ErrSemanticVersionFormat
	}

	version := SemanticVersion{}

	if major, err := strconv.Atoi(parts[0]); err == nil {
		version.Major = major
	} else {
		return nil, ErrPartCannotBe("major", parts[0])
	}

	if minor, err := strconv.Atoi(parts[1]); err == nil {
		version.Minor = minor
	} else {
		return nil, ErrPartCannotBe("minor", parts[1])
	}

	if patch, err := strconv.Atoi(parts[2]); err == nil {
		version.Patch = patch
	} else {
		return nil, ErrPartCannotBe("patch", parts[2])
	}

	return &version, nil
}

func (s SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func (s SemanticVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SemanticVersion) UnmarshalJSON(data []byte) error {
	var in string

	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	version, err := ParseSemanticVersion(in)
	if err != nil {
		return err
	}

	*s = *version

	return nil
}

func (s *SemanticVersion) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == "!!str" {
		version, err := ParseSemanticVersion(node.Value)
		if err != nil {
			return err
		}

		*s = *version

		return nil
	}

	return errors.New("[yaml] expected semantic version to be a string")
}

func (s *SemanticVersion) UnmarshalText(data []byte) error {
	version, err := ParseSemanticVersion(string(data))
	if err != nil {
		return err
	}

	*s = *version

	return nil
}

func NewVersion() SemanticVersion {
	return SemanticVersion{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}
}
