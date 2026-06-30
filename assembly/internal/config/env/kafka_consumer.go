package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type kafkaConsumerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID   string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type kafkaConsumerConfig struct {
	raw kafkaConsumerEnvConfig
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	var raw kafkaConsumerEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumerConfig{raw: raw}, nil
}

func (kcc *kafkaConsumerConfig) TopicName() string {
	return kcc.raw.TopicName
}

func (kcc *kafkaConsumerConfig) GroupID() string {
	return kcc.raw.GroupID
}

func (kcc *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
