package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL,required"`
	AsJson bool   `env:"LOGGER_AS_JSON,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (lc *loggerConfig) AsJson() bool {
	return lc.raw.AsJson
}

func (lc *loggerConfig) Level() string {
	return lc.raw.Level
}
