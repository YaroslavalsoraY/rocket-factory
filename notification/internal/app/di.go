package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
	"notification/internal/client/http"
	telegramClient "notification/internal/client/http/telegram"
	"notification/internal/config"
	converter "notification/internal/converter/kafka"
	"notification/internal/converter/kafka/decoder"
	"notification/internal/service"
	"notification/internal/service/consumer/order_assembled_consumer"
	"notification/internal/service/consumer/order_paid_consumer"
	"notification/internal/service/telegram"
	"platform/pkg/closer"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"
	middlewareKafka "platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramService service.TelegramService
	telegramClient  http.TelegramClient
	telegramBot     *bot.Bot

	orderPaidConsumerService service.ConsumerService
	orderPaidConsumerGroup   sarama.ConsumerGroup
	orderPaidConsumer        kafka.Consumer
	orderPaidDecoder         converter.OrderPaidDecoder

	orderAssembledConsumerService service.ConsumerService
	orderAssembledConsumerGroup   sarama.ConsumerGroup
	orderAssembledConsumer        kafka.Consumer
	orderAssembledDecoder         converter.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegram.NewService(
			d.TelegramClient(ctx),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) http.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = order_paid_consumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.TelegramService(ctx))
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() kafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = consumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.TopicName(),
			},
			logger.Logger(),
			middlewareKafka.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() converter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = order_assembled_consumer.NewService(d.OrderAssembledConsumer(), d.OrderAssembledDecoder(), d.TelegramService(ctx))
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		orderAssembledConsumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = orderAssembledConsumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledConsumer() kafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = consumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.TopicName(),
			},
			logger.Logger(),
			middlewareKafka.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() converter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}
