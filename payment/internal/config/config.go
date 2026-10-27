package config

import (
	"github.com/joho/godotenv"
	"payment/internal/config/env"
)

var appConfig *config

type config struct {
	PaymentGRPC PaymentGRPCConfig
	Logger      LoggerConfig
	OtelTracer  OtelTracerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	paymentGRPCConfig, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	otelTracerConfig, err := env.NewOtelTracerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		PaymentGRPC: paymentGRPCConfig,
		Logger:      loggerConfig,
		OtelTracer:  otelTracerConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
