package config

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

func (l *Logging) UnmarshalTOML(data interface{}) (err error) {
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

	return err
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
