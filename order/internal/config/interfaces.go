package config

import (
	"time"

	"github.com/IBM/sarama"
)

type OrderHttpConfig interface {
	HttpAdress() string
	HeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type InventoryGRPCConfig interface {
	InventoryAdress() string
}

type PaymentGRPCConfig interface {
	PaymentAdress() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

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
