package config

import (
	"github.com/joho/godotenv"
	"inventory/internal/config/env"
)

var appConfig *config

type config struct {
	InventoryGRPC InventoryGRPCConfig
	Mongo         MongoConfig
	Logger        LoggerConfig
	IamGRPC       IamGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	inventoryGRPC, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongo, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	iam, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		InventoryGRPC: inventoryGRPC,
		Mongo:         mongo,
		Logger:        logger,
		IamGRPC:       iam,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
