package config

import "github.com/IBM/sarama"

type KafkaConfig interface {
	Brokers() []string
}

type KafkaProducerConfig interface {
	TopicName() string
	Config() *sarama.Config
}

type KafkaConsumerConfig interface {
	TopicName() string
	GroupID() string
	Config() *sarama.Config
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
