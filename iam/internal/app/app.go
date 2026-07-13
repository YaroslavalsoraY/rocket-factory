package app

import (
	"context"
	"fmt"
	"net"

	"iam/internal/config"
	"platform/pkg/closer"
	"platform/pkg/grpc/health"
	"platform/pkg/logger"
	"platform/pkg/migrator/pg"
	iam_v1 "shared/pkg/proto/iam/v1"
	user_v1 "shared/pkg/proto/user/v1"

	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initMigrations,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDi(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().IamGRPC.Address())
	if err != nil {
		return err
	}

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	user_v1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserApi(ctx))

	iam_v1.RegisterIamServiceServer(a.grpcServer, a.diContainer.IamApi(ctx))

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	conn, err := a.diContainer.Pool(ctx).Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	migrator := pg.NewMigrator(stdlib.OpenDB(*conn.Conn().Config().Copy()), config.AppConfig().Postgres.MigrationsDir())

	if err = migrator.Up(); err != nil {
		return err
	}

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IamService server listening on %s", config.AppConfig().IamGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
