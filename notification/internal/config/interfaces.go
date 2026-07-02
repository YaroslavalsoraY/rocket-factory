package config

import "github.com/IBM/sarama"

type OrderPaidConsumerConfig interface {
	TopicName() string
	GroupID() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	TopicName() string
	GroupID() string
	Config() *sarama.Config
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type TelegramBotConfig interface {
	Token() string
}
