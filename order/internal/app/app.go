package app

import (
	"context"
	"fmt"
	"net/http"

	"order/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
	"platform/pkg/migrator/pg"
	order_v1 "shared/pkg/openapi/order/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer crashed: %v", err)
		}
	}()

	// HTTP сервер
	go func() {
		if err := a.runHttpServer(ctx); err != nil {
			errCh <- fmt.Errorf("grpc server crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initHttpServer,
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

func (a *App) initHttpServer(ctx context.Context) error {
	orderServer, err := order_v1.NewServer(a.diContainer.OrderAPI(ctx))
	if err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(config.AppConfig().OrderHttp.ShutdownTimeout()))
	r.Use(middleware.Logger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.AppConfig().OrderHttp.HttpAdress(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHttp.HeaderTimeout(),
	}

	a.httpServer = server

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	conn, err := a.diContainer.pool.Acquire(ctx)
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

func (a *App) runHttpServer(ctx context.Context) error {
	closer.AddNamed("http server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	logger.Info(ctx, fmt.Sprintf("🚀 HTTP-сервер заказов запущен по адресу %s", config.AppConfig().OrderHttp.HttpAdress()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 ShipAssembled Kafka consumer running")

	err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}