package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHttpEnvConfig struct {
	Host            string        `env:"HTTP_HOST,required"`
	Port            string        `env:"HTTP_PORT,required"`
	HeaderTimeout   time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type orderHttpConfig struct {
	raw orderHttpEnvConfig
}

func NewOrderHttpConfig() (*orderHttpConfig, error) {
	var raw orderHttpEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &orderHttpConfig{raw: raw}, nil
}

func (cfg *orderHttpConfig) HttpAdress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderHttpConfig) HeaderTimeout() time.Duration {
	return cfg.raw.HeaderTimeout
}

func (cfg *orderHttpConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}
