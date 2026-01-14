package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	api "inventory/internal/api/inventory/v1"
	repository "inventory/internal/repository/part"
	service "inventory/internal/service/part"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Извлекаем имя метода из полного пути
		method := path.Base(info.FullMethod)

		// Логируем начало вызова метода
		log.Printf("💨 Started gRPC method %s\n", method)

		// Засекаем время начала выполнения
		startTime := time.Now()

		// Вызываем обработчик
		resp, err := handler(ctx, req)

		// Вычисляем длительность выполнения
		duration := time.Since(startTime)

		// Форматируем сообщение в зависимости от результата
		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration)
		} else {
			log.Printf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration)
		}

		return resp, err
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer(grpc.UnaryInterceptor(LoggerInterceptor()))

	repo := repository.NewInventory()
	service := service.NewService(repo)
	api := api.NewApi(service)

	inventory_v1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 Inventory gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down inventory gRPC server...")
	s.GracefulStop()
	log.Println("✅ Inventory server stopped")
}
