package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	api "order/internal/api/order/v1"
	inventoryClient "order/internal/client/grpc/inventory/v1"
	paymentClient "order/internal/client/grpc/payment/v1"
	repository "order/internal/repository/order"
	service "order/internal/service/order"
	"os"
	"os/signal"
	"syscall"
	"time"

	order_v1 "shared/pkg/openapi/order/v1"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	payment_v1 "shared/pkg/proto/payment/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort          = "8080"
	inventoryAdress   = "localhost:50051"
	paymentAdress     = "localhost:50052"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(inventoryAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ошибка создания клиента сервиса хранения: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := inventoryClient.NewClient(inventory_v1.NewInventoryServiceClient(inventoryConn))

	paymentConn, err := grpc.NewClient(paymentAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("ошибка создания клиента сервиса оплаты: %v", err)
	}
	defer paymentConn.Close()
	paymentClient := paymentClient.NewClient(payment_v1.NewPaymentServiceClient(paymentConn))
	
	repository := repository.NewRepository()
	service := service.NewService(repository, inventoryClient, paymentClient)
	api := api.NewApi(service)

	orderServer, err := order_v1.NewServer(api)
	if err != nil {
		log.Fatalf("ошибка создания сервера заказов: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, 
	}

	go func() {
		log.Printf("🚀 HTTP-сервер заказов запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера заказов: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
