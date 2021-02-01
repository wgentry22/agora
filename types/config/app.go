package config

import (
	"time"

	"github.com/hashicorp/errwrap"
)

var (
	defaultApplicationName    = "agora-app"
	defaultApplicationVersion = NewVersion()
	defaultApplicationEnv     = Development
	defaultAPIPathPrefix      = "/v1"
	defaultAPIPort            = 8123
	defaultLoggingLevel       = "debug"
	defaultLoggingOutputPaths = []string{"stdout"}
	defaultTimeoutDuration    = 5000 * time.Millisecond
)

type Application struct {
	info      Info
	api       API
	logging   Logging
	db        DB
	heartbeat Heartbeat
	broker    Broker
}

func (a Application) Heartbeat() Heartbeat {
	return a.heartbeat
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

func (a Application) Broker() Broker {
	return a.broker
}

func (a *Application) UnmarshalTOML(data interface{}) (err error) {
	dataMap := data.(map[string]interface{})

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

	a.api = a.api.withInfo(a.info)

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

	if heartbeat, ok := dataMap["heartbeat"]; ok {
		var heartbeatConfig Heartbeat
		if heartbeatErr := heartbeatConfig.UnmarshalTOML(heartbeat); heartbeatErr != nil {
			err = errwrap.Wrap(heartbeatErr, err)
		} else {
			a.heartbeat = heartbeatConfig
		}
	} else {
		a.heartbeat = defaultHeartbeatConfig()
	}

	a.heartbeat = a.heartbeat.WithInfo(a.info)

	if db, ok := dataMap["db"]; ok {
		var dbConfig DB
		if dbErr := dbConfig.UnmarshalTOML(db); dbErr != nil {
			err = errwrap.Wrap(dbErr, err)
		} else {
			a.db = dbConfig
		}
	}

	if broker, ok := dataMap["broker"]; ok {
		var conf Broker
		if brokerErr := conf.UnmarshalTOML(broker); brokerErr != nil {
			err = errwrap.Wrap(brokerErr, err)
		} else {
			a.broker = conf
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

func (i *Info) UnmarshalTOML(data interface{}) (err error) {
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
