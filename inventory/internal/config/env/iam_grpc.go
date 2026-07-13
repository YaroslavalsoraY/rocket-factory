package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamGRPCEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type iamGRPCConfig struct {
	raw iamGRPCEnvConfig
}

func NewIamGRPCConfig() (*iamGRPCConfig, error) {
	var raw iamGRPCEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &iamGRPCConfig{raw: raw}, nil
}

func (cfg *iamGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
