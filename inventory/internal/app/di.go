package app

import (
	"context"
	"fmt"

	api "inventory/internal/api/inventory/v1"
	clientInterface "inventory/internal/client/grpc"
	iamClient "inventory/internal/client/grpc/iam/v1"
	"inventory/internal/config"
	"inventory/internal/repository"
	repo "inventory/internal/repository/part"
	"inventory/internal/service"
	partService "inventory/internal/service/part"
	"platform/pkg/closer"
	iam_v1 "shared/pkg/proto/iam/v1"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	interceptor "platform/pkg/middleware/grpc"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	invRepsitory repository.InventoryRepository

	invService service.InventoryService

	invAPI inventory_v1.InventoryServiceServer

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	iamClient clientInterface.IamClient
	authInterceptor *interceptor.AuthInterceptor
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) IamClient() clientInterface.IamClient {
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

func (di *diContainer) AuthInterceptor() *interceptor.AuthInterceptor {
	if di.authInterceptor == nil {
		di.authInterceptor = interceptor.NewAuthInterceptor(di.IamClient())
	}

	return di.authInterceptor
}

func (di *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if di.mongoDBClient == nil {
		client, err := mongo.Connect(
			ctx,
			options.Client().ApplyURI(config.AppConfig().Mongo.URI()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to mongo: %v", err))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping mongo: %v", err))
		}

		closer.AddNamed("mongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		di.mongoDBClient = client
	}

	return di.mongoDBClient
}

func (di *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if di.mongoDBHandle == nil {
		di.mongoDBHandle = di.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return di.mongoDBHandle
}

func (di *diContainer) InvRepsitory(ctx context.Context) repository.InventoryRepository {
	if di.invRepsitory == nil {
		di.invRepsitory = repo.NewInventory(di.MongoDBHandle(ctx))
	}

	return di.invRepsitory
}

func (di *diContainer) InvService(ctx context.Context) service.InventoryService {
	if di.invService == nil {
		di.invService = partService.NewService(di.InvRepsitory(ctx))
	}

	return di.invService
}

func (di *diContainer) InvAPI(ctx context.Context) inventory_v1.InventoryServiceServer {
	if di.invAPI == nil {
		di.invAPI = api.NewApi(di.InvService(ctx))
	}

	return di.invAPI
}
