package config

import (
	"github.com/joho/godotenv"
	"notification/internal/config/env"
)

var appConfig *config

type config struct {
	Kafka                  KafkaConfig
	Logger                 LoggerConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	TelegramBot            TelegramBotConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	kafka, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumer, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumer, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	telegramBot, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Kafka:                  kafka,
		OrderPaidConsumer:      orderPaidConsumer,
		OrderAssembledConsumer: orderAssembledConsumer,
		TelegramBot:            telegramBot,
		Logger:                 logger,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
