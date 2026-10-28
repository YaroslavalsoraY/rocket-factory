package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"assembly/internal/config"
	"go.uber.org/zap"
	"platform/pkg/closer"
	"platform/pkg/logger"
	"platform/pkg/metrics"
)

type App struct {
	diContainer    *diContainer
	metricsHandler http.Handler
	metricsServer  *http.Server
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
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runConsumer(ctx); err != nil && ctx.Err() == nil {
			errCh <- fmt.Errorf("consumer: %w", err)
		}
	}()

	go func() {
		if err := a.runMetricsServer(ctx); err != nil {
			errCh <- fmt.Errorf("metrics HTTP server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initMetrics,
		a.initMetricsServer,
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

func (a *App) initMetrics(_ context.Context) error {
	handler, err := metrics.InitProvider("assembly")
	if err != nil {
		return err
	}

	a.metricsHandler = handler

	closer.AddNamed("Metrics", metrics.Shutdown)

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	err := a.diContainer.ConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initMetricsServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", a.metricsHandler)

	a.metricsServer = &http.Server{
		Addr:              config.AppConfig().Metrics.Address(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	closer.AddNamed("Assembly metrics HTTP", a.metricsServer.Shutdown)
	return nil
}

func (a *App) runMetricsServer(ctx context.Context) error {
	logger.Info(ctx, "Starting Assembly metrics HTTP server", zap.String("address", a.metricsServer.Addr))

	err := a.metricsServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
