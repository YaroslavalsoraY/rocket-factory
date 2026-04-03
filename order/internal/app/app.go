package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	"order/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"
	"platform/pkg/migrator/pg"
	order_v1 "shared/pkg/openapi/order/v1"
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
	return a.runHttpServer(ctx)
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
