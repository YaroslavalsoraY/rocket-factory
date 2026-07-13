package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"iam/internal/app"
	"iam/internal/config"
	"platform/pkg/closer"
	"platform/pkg/logger"

	"go.uber.org/zap"
)

const configPath = "./deploy/compose/iam/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		if logger.Logger() != nil {
			logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
		} else {
			log.Printf("failed to close resources: %v\n", err)
		}
	}
}
