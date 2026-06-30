package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type kafkaProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type kafkaProducerConfig struct {
	raw kafkaProducerEnvConfig
}

func NewKafkaProducerConfig() (*kafkaProducerConfig, error) {
	var raw kafkaProducerEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &kafkaProducerConfig{raw: raw}, nil
}

func (kpc *kafkaProducerConfig) TopicName() string {
	return kpc.raw.TopicName
}

func (kpc *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
