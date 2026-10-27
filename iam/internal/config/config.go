package config

import (
	"github.com/joho/godotenv"
	"iam/internal/config/env"
)

var appConfig *config

type config struct {
	IamGRPC  IamGRPCConfig
	Logger   LoggerConfig
	Redis    RedisConfig
	Session  SessionConfig
	Postgres PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	IamGRPCConfig, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionConfig, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		IamGRPC:  IamGRPCConfig,
		Logger:   loggerConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		Session:  sessionConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
