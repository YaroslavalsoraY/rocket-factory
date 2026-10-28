package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	api "order/internal/api/order/v1"
	client "order/internal/client/grpc"
	iamClient "order/internal/client/grpc/iam/v1"
	inventoryClient "order/internal/client/grpc/inventory/v1"
	paymentClient "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	kafkaConverter "order/internal/converter/kafka"
	"order/internal/converter/kafka/decoder"
	"order/internal/repository"
	repo "order/internal/repository/order"
	"order/internal/service"
	order_consumer "order/internal/service/consumer"
	orderService "order/internal/service/order"
	producerService "order/internal/service/producer"
	"platform/pkg/closer"
	"platform/pkg/kafka"
	"platform/pkg/kafka/consumer"
	"platform/pkg/kafka/producer"
	"platform/pkg/logger"
	authMiddleware "platform/pkg/middleware/http"
	middlewareKafka "platform/pkg/middleware/kafka"
	"platform/pkg/tracing"
	order_v1 "shared/pkg/openapi/order/v1"
	iam_v1 "shared/pkg/proto/iam/v1"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	payment_v1 "shared/pkg/proto/payment/v1"
)

type diContainer struct {
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
	pool            *pgxpool.Pool
	orderRepository repository.OrderRepository
	orderService    service.OrderService
	orderAPI        order_v1.Handler

	producerService service.ProducerService
	consumerService service.ConsumerService

	consumerGroup         sarama.ConsumerGroup
	shipAssembledConsumer kafka.Consumer

	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
	syncProducer         sarama.SyncProducer
	orderPaidProducer    kafka.Producer

	iamClient      client.IamClient
	authMiddleware *authMiddleware.AuthMiddleware

	ordersCounter metric.Int64Counter
	ordersRevenue metric.Float64Counter
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if di.orderRepository == nil {
		di.orderRepository = repo.NewRepository(di.Pool(ctx))
	}

	return di.orderRepository
}

func (di *diContainer) OrderService(ctx context.Context) service.OrderService {
	if di.orderService == nil {
		di.orderService = orderService.NewService(
			di.OrderRepository(ctx),
			di.InventoryClient(ctx),
			di.PaymentClient(ctx),
			di.ProducerService(),
			di.OrdersCounter(),
			di.OrdersRevenue(),
		)
	}

	return di.orderService
}

func (di *diContainer) OrderAPI(ctx context.Context) order_v1.Handler {
	if di.orderAPI == nil {
		di.orderAPI = api.NewApi(di.OrderService(ctx))
	}

	return di.orderAPI
}

func (di *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if di.inventoryClient == nil {
		cc, err := grpc.NewClient(config.AppConfig().InventoryGRPC.InventoryAdress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to iventory service: %v", err))
		}

		closer.AddNamed("connection to inventory service", func(ctx context.Context) error {
			return cc.Close()
		})

		inventoryClient := inventoryClient.NewClient(inventory_v1.NewInventoryServiceClient(cc))

		di.inventoryClient = inventoryClient
	}

	return di.inventoryClient
}

func (di *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if di.paymentClient == nil {
		cc, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.PaymentAdress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(tracing.UnaryClientInterceptor("payment-service")),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment service: %v", err))
		}

		closer.AddNamed("connection to payment service", func(ctx context.Context) error {
			return cc.Close()
		})

		paymentClient := paymentClient.NewClient(payment_v1.NewPaymentServiceClient(cc))

		di.paymentClient = paymentClient
	}

	return di.paymentClient
}

func (di *diContainer) Pool(ctx context.Context) *pgxpool.Pool {
	if di.pool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to postgres: %v", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping postgres: %v", err))
		}

		closer.AddNamed("pgxpool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		di.pool = pool
	}

	return di.pool
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = order_consumer.NewService(d.ShipAssembledConsumer(), d.ShipAssembledDecoder(), d.OrderRepository(ctx))
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

func (d *diContainer) ShipAssembledConsumer() kafka.Consumer {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = consumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().KafkaConsumer.TopicName(),
			},
			logger.Logger(),
			middlewareKafka.Logging(logger.Logger()),
		)
	}

	return d.shipAssembledConsumer
}

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
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

func (d *diContainer) OrderPaidProducer() kafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = producer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().KafkaProducer.TopicName(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
}

func (d *diContainer) ProducerService() service.ProducerService {
	if d.producerService == nil {
		d.producerService = producerService.NewService(d.OrderPaidProducer())
	}

	return d.producerService
}

func (di *diContainer) IamClient() client.IamClient {
	if di.iamClient == nil {
		cc, err := grpc.NewClient(config.AppConfig().IamGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to iam service: %v", err))
		}

		closer.AddNamed("connection to iam service", func(ctx context.Context) error {
			return cc.Close()
		})

		iamClient := iamClient.NewClient(iam_v1.NewIamServiceClient(cc))

		di.iamClient = iamClient
	}

	return di.iamClient
}

func (di *diContainer) OrdersCounter() metric.Int64Counter {
	if di.ordersCounter == nil {
		counter, err := otel.Meter("order").Int64Counter("orders")
		if err != nil {
			panic(fmt.Errorf("create orders counter: %w", err))
		}
		di.ordersCounter = counter
	}

	return di.ordersCounter
}

func (di *diContainer) OrdersRevenue() metric.Float64Counter {
	if di.ordersRevenue == nil {
		counter, err := otel.Meter("order").Float64Counter("orders_revenue")
		if err != nil {
			panic(fmt.Errorf("create orders revenue counter: %w", err))
		}
		di.ordersRevenue = counter
	}

	return di.ordersRevenue
}

func (di *diContainer) AuthMiddleware() *authMiddleware.AuthMiddleware {
	if di.authMiddleware == nil {
		di.authMiddleware = authMiddleware.NewAuthMiddleware(di.IamClient())
	}

	return di.authMiddleware
}
