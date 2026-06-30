package config

import (
	"assembly/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Kafka         KafkaConfig
	Logger        LoggerConfig
	KafkaProducer KafkaProducerConfig
	KafkaConsumer KafkaConsumerConfig
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

	kafkaProducer, err := env.NewKafkaProducerConfig()
	if err != nil {
		return err
	}

	kafkaConsumer, err := env.NewKafkaConsumerConfig()
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Kafka:         kafka,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
		Logger:        logger,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
