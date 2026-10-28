package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	ServiceName string `env:"METRICS_SERVICE_NAME,required"`
	Host        string `env:"METRICS_HOST,required"`
	Port        string `env:"METRICS_PORT,required"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

func (mc *metricsConfig) ServiceName() string {
	return mc.raw.ServiceName
}

func (mc *metricsConfig) Address() string {
	return net.JoinHostPort(mc.raw.Host, mc.raw.Port)
}
