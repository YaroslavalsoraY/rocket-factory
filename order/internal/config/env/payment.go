package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentConfig struct {
	raw paymentEnvConfig
}

func NewPaymentConfig() (*paymentConfig, error) {
	var raw paymentEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &paymentConfig{raw: raw}, nil
}

func (cfg *paymentConfig) PaymentAdress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
