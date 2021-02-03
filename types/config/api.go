package config

import (
  "fmt"
  "strings"
  "time"

  "github.com/gin-contrib/cors"
)

type API struct {
  Port       int            `json:"port" yaml:"port" toml:"port"`
  PathPrefix string         `json:"pathPrefix" yaml:"pathPrefix" toml:"pathPrefix"`
  info       *Info          `toml:"-"`
  Timeout    TimeoutOptions `toml:"timeout"`
  Cors       CORS           `toml:"cors"`
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

  if cors, ok := dataMap["cors"]; ok {
    var corsConf CORS
    if err := corsConf.UnmarshalTOML(cors); err != nil {
      return err
    }

    a.Cors = corsConf
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

func (a API) ShouldRegisterCors() bool {
  return len(a.Cors.AllowOrigins) > 0 ||
    len(a.Cors.AllowMethods) > 0 ||
    len(a.Cors.AllowHeaders) > 0 ||
    len(a.Cors.ExposeHeaders) > 0
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

type CORS struct {
  AllowOrigins     []string `toml:"allow-origins"`
  AllowMethods     []string `toml:"allow-methods"`
  AllowHeaders     []string `toml:"allow-headers"`
  ExposeHeaders    []string `toml:"expose-headers"`
  AllowCredentials bool     `toml:"allow-credentials"`
}

func (c CORS) ToGinConfig() cors.Config {
  conf := cors.Config{}

  if len(c.AllowOrigins) > 0 {
    conf.AllowOrigins = c.AllowOrigins
  }

  if len(c.AllowMethods) > 0 {
    conf.AllowMethods = c.AllowMethods
  }

  if len(c.AllowHeaders) > 0 {
    conf.AllowHeaders = c.AllowHeaders
  }

  if len(c.ExposeHeaders) > 0 {
    conf.ExposeHeaders = c.ExposeHeaders
  }

  conf.AllowCredentials = c.AllowCredentials

  return conf
}

func (c *CORS) UnmarshalTOML(data interface{}) error {
  dataMap := data.(map[string]interface{})

  c.AllowOrigins = getStringSliceFromMap("allow-origins", dataMap)
  c.AllowMethods = getStringSliceFromMap("allow-methods", dataMap)
  c.AllowHeaders = getStringSliceFromMap("allow-headers", dataMap)
  c.ExposeHeaders = getStringSliceFromMap("expose-headers", dataMap)

  if allow, ok := dataMap["allow-credentials"].(bool); ok {
    c.AllowCredentials = allow
  }

  return nil
}

func getStringSliceFromMap(key string, m map[string]interface{}) []string {
  if vals, ok := m[key].([]interface{}); ok {
    slice := make([]string, len(vals))

    for i, val := range vals {
      slice[i] = val.(string)
    }

    return slice
  }

  return []string{}
}
