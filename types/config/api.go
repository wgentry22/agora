package config

import (
  "fmt"
  "strings"
  "time"
)

type API struct {
  Port       int            `json:"port" yaml:"port" toml:"port"`
  PathPrefix string         `json:"pathPrefix" yaml:"pathPrefix" toml:"pathPrefix"`
  info       *Info          `toml:"-"`
  Timeout    TimeoutOptions `toml:"timeout"`
}

func defaultAPIServer() API {
  return API{
    Port:       defaultAPIPort,
    PathPrefix: defaultAPIPathPrefix,
  }
}

func (a API) ListenAddr() string {
  return fmt.Sprintf(":%d", a.Port)
}

func (a API) withInfo(info Info) API {
  return API{
    Port:       a.Port,
    PathPrefix: a.PathPrefix,
    info:       &info,
  }
}

func (a API) Info() Info {
  if a.info == nil {
    return defaultInfo()
  }

  return *a.info
}

func (a *API) UnmarshalTOML(data interface{}) error {
  dataMap := data.(map[string]interface{})

  if pathPrefix, ok := dataMap["pathPrefix"].(string); ok {
    if !strings.HasPrefix(pathPrefix, "/") {
      pathPrefix = fmt.Sprintf("/%s", pathPrefix)
    }

    a.PathPrefix = pathPrefix
  } else {
    a.PathPrefix = defaultAPIPathPrefix
  }

  if port, ok := dataMap["port"].(int64); ok {
    a.Port = int(port)
  } else {
    a.Port = defaultAPIPort
  }

  if timeout, ok := dataMap["timeout"]; ok {
    var opts TimeoutOptions
    if err := opts.UnmarshalTOML(timeout); err != nil {
      return err
    }

    a.Timeout = opts
  } else {
    a.Timeout = defaultTimeoutOptions()
  }

  return nil
}

func (a API) WithInfo(info Info) API {
  return API{
    Port:       a.Port,
    PathPrefix: a.PathPrefix,
    info:       &info,
  }
}

func (a API) WithDefaultInfo() API {
  info := defaultInfo()

  return API{
    Port:       a.Port,
    PathPrefix: a.PathPrefix,
    info:       &info,
  }
}

type TimeoutOptions struct {
  Read  time.Duration `toml:"read"`
  Write time.Duration `toml:"write"`
}

func defaultTimeoutOptions() TimeoutOptions {
  return TimeoutOptions{
    Read:  defaultTimeoutDuration,
    Write: defaultTimeoutDuration,
  }
}

func (t *TimeoutOptions) UnmarshalTOML(data interface{}) error {
  dataMap := data.(map[string]interface{})

  if read, ok := dataMap["read"].(int64); ok {
    if read <= 0 {
      read = 5000
    }

    t.Read = time.Duration(read) * time.Millisecond
  } else {
    t.Read = defaultTimeoutDuration
  }

  if write, ok := dataMap["write"].(int64); ok {
    if write <= 0 {
      write = 5000
    }

    t.Write = time.Duration(write) * time.Millisecond
  } else {
    t.Write = defaultTimeoutDuration
  }

  return nil
}
