package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host          string `env:"REDIS_HOST,required"`
	Port          string `env:"EXTERNAL_REDIS_PORT,required"`
	MaxIdle       int `env:"REDIS_MAX_IDLE,required"`
	Timeout       time.Duration `env:"REDIS_CONNECTION_TIMEOUT,required"`
	IdleTimeout   time.Duration `env:"REDIS_IDLE_TIMEOUT,required"`
}

type redisConfig struct {
	raw redisEnvConfig
}

func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &redisConfig{raw: raw}, nil
}

func (cfg *redisConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *redisConfig) MaxIdle() int {
	return cfg.raw.MaxIdle
}

func (cfg *redisConfig) Timeout() time.Duration {
	return cfg.raw.Timeout
}

func (cfg *redisConfig) IdleTimeout() time.Duration {
	return cfg.raw.IdleTimeout
}
