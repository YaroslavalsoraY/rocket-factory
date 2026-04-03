package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryConfig struct {
	raw inventoryEnvConfig
}

func NewInventoryConfig() (*inventoryConfig, error) {
	var raw inventoryEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryConfig{raw: raw}, nil
}

func (cfg *inventoryConfig) InventoryAdress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
