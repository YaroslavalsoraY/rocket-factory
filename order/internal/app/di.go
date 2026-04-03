package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	api "order/internal/api/order/v1"
	client "order/internal/client/grpc"
	inventoryClient "order/internal/client/grpc/inventory/v1"
	paymentClient "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/repository"
	repo "order/internal/repository/order"
	"order/internal/service"
	orderService "order/internal/service/order"
	"platform/pkg/closer"
	order_v1 "shared/pkg/openapi/order/v1"
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
		di.orderService = orderService.NewService(di.OrderRepository(ctx), di.InventoryClient(ctx), di.PaymentClient(ctx))
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
		cc, err := grpc.NewClient(config.AppConfig().PaymentGRPC.PaymentAdress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
