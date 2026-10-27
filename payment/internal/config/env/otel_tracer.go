package env

import "github.com/caarlos0/env/v11"

type otelTracerEnvConfig struct {
	OtelExporterEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,required"`
	OtelServiceName      string `env:"OTEL_SERVICE_NAME,required"`
	OtelEnvironment      string `env:"OTEL_ENVIRONMENT,required"`
	OtelServiceVersion   string `env:"OTEL_SERVICE_VERSION,required"`
}

type otelTracerConfig struct {
	raw otelTracerEnvConfig
}

func NewOtelTracerConfig() (*otelTracerConfig, error) {
	var raw otelTracerEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &otelTracerConfig{raw: raw}, nil
}

func (cfg *otelTracerConfig) CollectorEndpoint() string {
	return cfg.raw.OtelExporterEndpoint
}

func (cfg *otelTracerConfig) ServiceName() string {
	return cfg.raw.OtelServiceName
}

func (cfg *otelTracerConfig) Environment() string {
	return cfg.raw.OtelEnvironment
}

func (cfg *otelTracerConfig) ServiceVersion() string {
	return cfg.raw.OtelServiceVersion
}
