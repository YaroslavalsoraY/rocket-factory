package config

import (
	"github.com/joho/godotenv"
	"order/internal/config/env"
)

var appConfig *config

type config struct {
	OrderHttp     OrderHttpConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
	Logger        LoggerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	orderHttp, err := env.NewOrderHttpConfig()
	if err != nil {
		return err
	}

	postgres, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPC, err := env.NewInventoryConfig()
	if err != nil {
		return err
	}

	paymentGRPC, err := env.NewPaymentConfig()
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		OrderHttp:     orderHttp,
		Postgres:      postgres,
		InventoryGRPC: inventoryGRPC,
		PaymentGRPC:   paymentGRPC,
		Logger:        logger,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
