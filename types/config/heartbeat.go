package config

import (
	"fmt"
	"strings"
)

type Heartbeat struct {
	PathPrefix string         `json:"pathPrefix" yaml:"pathPrefix" toml:"pathPrefix"`
	info       *Info          `toml:"-"`
	Timeout    TimeoutOptions `json:"timeout"`
}

func (h *Heartbeat) UnmarshalTOML(data interface{}) error {
	dataMap := data.(map[string]interface{})

	if path, ok := dataMap["pathPrefix"].(string); ok {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		h.PathPrefix = path
	} else {
		h.PathPrefix = "/heartbeat"
	}

	if timeout, ok := dataMap["timeout"]; ok {
		var opts TimeoutOptions
		if err := opts.UnmarshalTOML(timeout); err != nil {
			return err
		}

		h.Timeout = opts
	} else {
		h.Timeout = defaultTimeoutOptions()
	}

	return nil
}

func (h Heartbeat) WithInfo(info Info) Heartbeat {
	return Heartbeat{
		PathPrefix: h.PathPrefix,
		info:       &info,
		Timeout:    h.Timeout,
	}
}

func (h Heartbeat) Info() Info {
	if h.info == nil {
		return defaultInfo()
	}

	return *h.info
}

func defaultHeartbeatConfig() Heartbeat {
	return Heartbeat{
		PathPrefix: "/heartbeat",
		info:       nil,
		Timeout:    defaultTimeoutOptions(),
	}
}
