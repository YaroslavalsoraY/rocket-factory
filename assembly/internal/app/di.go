package app

import (
	"context"
	"fmt"

	"assembly/internal/config"
	kafkaConverter "assembly/internal/converter/kafka"
	decoder "assembly/internal/converter/kafka/decoder"
	"assembly/internal/service"
	assebmle_consumer "assembly/internal/service/consumer/assemble_consumer"
	"assembly/internal/service/producer/order_producer"
	"github.com/IBM/sarama"
	"platform/pkg/closer"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/kafka/producer"
	"platform/pkg/logger"
	middlewareKafka "platform/pkg/middleware/kafka"
)

type diContainer struct {
	producerService service.ProducerService
	consumerService service.ConsumerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer kafka.Consumer

	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	syncProducer          sarama.SyncProducer
	shipAssembledProducer kafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ProducerService() service.ProducerService {
	if d.producerService == nil {
		d.producerService = order_producer.NewService(d.ShipAssembledProducer())
	}

	return d.producerService
}

func (d *diContainer) ConsumerService() service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = assebmle_consumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.ProducerService())
	}

	return d.consumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().KafkaConsumer.GroupID(),
			config.AppConfig().KafkaConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidConsumer() kafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = consumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().KafkaConsumer.TopicName(),
			},
			logger.Logger(),
			middlewareKafka.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().KafkaProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) ShipAssembledProducer() kafka.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = producer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().KafkaProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.shipAssembledProducer
}
