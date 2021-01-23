package config

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"strings"
)

var (
	defaultApplicationName    = "agora-app"
	defaultApplicationVersion = NewVersion()
	defaultApplicationEnv     = Development
	defaultAPIPathPrefix      = "/v1"
	defaultAPIPort            = 8123
	defaultLoggingLevel       = "debug"
	defaultLoggingOutputPaths = []string{"stdout"}
)

type Application struct {
	info    Info
	api     API
	logging Logging
	db      DB
}

func (a Application) Logging() Logging {
	return a.logging.WithFields(a.info.fields())
}

func (a Application) Info() Info {
	return a.info
}

func (a Application) API() API {
	return a.api.withInfo(a.info)
}

func (a Application) DB() DB {
	return a.db
}

func (a *Application) UnmarshalTOML(data interface{}) error {
	var err error
	dataMap := data.(map[string]interface{})

	if api, ok := dataMap["api"]; ok {
		var apiConfig API
		if apiErr := apiConfig.UnmarshalTOML(api); apiErr != nil {
			err = errwrap.Wrap(apiErr, err)
		} else {
			a.api = apiConfig
		}
	} else {
		a.api = defaultAPIServer()
	}

	if info, ok := dataMap["info"]; ok {
		var appInfo Info
		if infoErr := appInfo.UnmarshalTOML(info); infoErr != nil {
			err = errwrap.Wrap(infoErr, err)
		} else {
			a.info = appInfo
		}
	} else {
		a.info = defaultInfo()
	}

	if logging, ok := dataMap["logging"]; ok {
		var loggingConfig Logging
		if logErr := loggingConfig.UnmarshalTOML(logging); logErr != nil {
			err = errwrap.Wrap(logErr, err)
		} else {
			a.logging = loggingConfig
		}
	} else {
		a.logging = defaultLoggingWithFields(a.info.fields())
	}

	if db, ok := dataMap["db"]; ok {
		var dbConfig DB
		if dbErr := dbConfig.UnmarshalTOML(db); dbErr != nil {
			err = errwrap.Wrap(dbErr, err)
		} else {
			a.db = dbConfig
		}
	}

	return err
}

type Info struct {
	Name    string          `json:"name" yaml:"name" toml:"name"`
	Version SemanticVersion `json:"version" yaml:"name" toml:"name"`
	Env     Environment     `json:"env" yaml:"env" toml:"env"`
}

func defaultInfo() Info {
	return Info{
		Name:    defaultApplicationName,
		Version: defaultApplicationVersion,
		Env:     defaultApplicationEnv,
	}
}

func (i *Info) UnmarshalTOML(data interface{}) error {
	var err error
	dataMap := data.(map[string]interface{})

	if name, ok := dataMap["name"].(string); ok && name != "" {
		i.Name = name
	} else {
		i.Name = defaultApplicationName
	}

	if version, ok := dataMap["version"].(string); ok {
		semVer, semVerErr := ParseSemanticVersion(version)
		if semVerErr != nil {
			err = errwrap.Wrap(semVerErr, err)
		} else {
			i.Version = *semVer
		}
	} else {
		i.Version = defaultApplicationVersion
	}

	if environment, ok := dataMap["env"].(string); ok {
		env, envErr := ParseEnvironment(environment)

		if envErr != nil {
			err = errwrap.Wrap(envErr, err)
		} else {
			i.Env = env
		}
	} else {
		i.Env = defaultApplicationEnv
	}

	return err
}

func (i Info) fields() map[string]interface{} {
	fields := make(map[string]interface{})

	if i.Name == "" {
		fields["name"] = "agora-app"
	} else {
		fields["name"] = i.Name
	}

	if i.Version.IsZero() {
		fields["version"] = NewVersion().String()
	} else {
		fields["version"] = i.Version.String()
	}

	fields["env"] = i.Env.String()

	return fields
}

type API struct {
	Port       int    `json:"port" yaml:"port" toml:"port"`
	PathPrefix string `json:"pathPrefix" yaml:"pathPrefix" toml:"pathPrefix"`
	info       *Info  `toml:"-"`
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

	return nil
}

type Logging struct {
	Level       string                 `json:"level" yaml:"level" toml:"level"`
	OutputPaths []string               `json:"outputPaths" yaml:"outputPaths" toml:"outputPaths"`
	Fields      map[string]interface{} `json:"fields" yaml:"fields" toml:"fields"`
}

func defaultLoggingWithFields(fields map[string]interface{}) Logging {
	return Logging{
		Level:       defaultLoggingLevel,
		OutputPaths: defaultLoggingOutputPaths,
		Fields:      fields,
	}
}

func (l *Logging) UnmarshalTOML(data interface{}) error {
	dataMap := data.(map[string]interface{})

	if level, ok := dataMap["level"].(string); ok && level != "" {
		l.Level = level
	} else {
		l.Level = defaultLoggingLevel
	}

	if paths, ok := dataMap["outputPaths"].([]interface{}); ok && len(paths) > 0 {
		out := make([]string, len(paths))
		for i, path := range paths {
			out[i] = path.(string)
		}

		l.OutputPaths = out
	} else {
		l.OutputPaths = defaultLoggingOutputPaths
	}

	if fields, ok := dataMap["fields"].(map[string]interface{}); ok {
		l.Fields = fields
	} else {
		l.Fields = make(map[string]interface{})
	}

	return nil
}

func (l *Logging) WithFields(fields map[string]interface{}) Logging {
	newFields := make(map[string]interface{})

	for k, v := range l.Fields {
		newFields[k] = v
	}

	for k, v := range fields {
		newFields[k] = v
	}

	return Logging{l.Level, l.OutputPaths, newFields}
}

type DB struct {
	vendor   DBVendor
	user     string
	password string
	host     string
	port     int
	name     string
	args     map[string]interface{}
}

func (d *DB) UnmarshalTOML(data interface{}) error {
	var err error
	dataMap := data.(map[string]interface{})

	if vendor, ok := dataMap["vendor"].(string); ok && vendor != "" {
		parsed, vendorErr := ParseDBVendor(vendor)
		if vendorErr != nil {
			err = errwrap.Wrap(vendorErr, err)
		} else {
			d.vendor = parsed
		}
	}

	if user, ok := dataMap["user"].(string); ok && user != "" {
		d.user = user
	}

	if password, ok := dataMap["password"].(string); ok && password != "" {
		d.password = password
	}

	if host, ok := dataMap["host"].(string); ok && host != "" {
		d.host = host
	}

	if name, ok := dataMap["name"].(string); ok && name != "" {
		d.name = name
	}

	if port, ok := dataMap["port"].(int64); ok && port > 0 {
		d.port = int(port)
	}

	if args, ok := dataMap["args"].(map[string]interface{}); ok {
		d.args = args
	}

	return err
}
