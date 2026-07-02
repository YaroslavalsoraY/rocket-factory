package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledConsumerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID   string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledConsumerConfig struct {
	raw orderAssembledConsumerEnvConfig
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &orderAssembledConsumerConfig{raw: raw}, nil
}

func (oac *orderAssembledConsumerConfig) TopicName() string {
	return oac.raw.TopicName
}

func (oac *orderAssembledConsumerConfig) GroupID() string {
	return oac.raw.GroupID
}

func (oac *orderAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
