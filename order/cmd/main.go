package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	order_v1 "shared/pkg/openapi/order/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

const (
	errorPaid      = "Order is already paid"
	errorNoOrders  = "No orders were found"
	errorCancelled = "Order was cancelled"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*order_v1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*order_v1.OrderDto),
	}
}

func (s *OrderStorage) GetOrder(uuid uuid.UUID) *order_v1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) PostOrder(order *order_v1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

func (s *OrderStorage) CancelOrder(uuid uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[uuid]
	if !ok {
		return errors.New(errorNoOrders)
	}

	if s.orders[uuid].Status == order_v1.OrderStatusPAID {
		return errors.New(errorPaid)
	}

	s.orders[uuid].Status = order_v1.OrderStatusCANCELLED

	return nil
}

func (s *OrderStorage) PayOrder(uuid uuid.UUID, paymentMethod order_v1.PaymentMethod, transactionUUID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[uuid]
	if !ok {
		return errors.New(errorNoOrders)
	}

	if s.orders[uuid].Status == order_v1.OrderStatusPAID {
		return errors.New(errorPaid)
	}

	if s.orders[uuid].Status == order_v1.OrderStatusCANCELLED {
		return errors.New(errorCancelled)
	}

	s.orders[uuid].Status = order_v1.OrderStatusPAID
	s.orders[uuid].PaymentMethod = paymentMethod
	s.orders[uuid].TransactionalUUID = transactionUUID
	
	return nil
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	err := h.storage.CancelOrder(params.OrderUUID)

	if err != nil {
		switch err.Error() {
		case errorNoOrders:
			return &order_v1.NotFoundError{
				Code: 404,
				Message: "Order was not found",
			}, nil
		case errorPaid:
			return &order_v1.ConflictError{
				Code: 409,
				Message: "Cannot cancel paid order",
			}, nil
		}
	}

	return &order_v1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) CreateNewOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateNewOrderRes, error) {
	orderUUID := uuid.New()
	newOrder := order_v1.OrderDto{
		OrderUUID: orderUUID,
		UserUUID: req.UserUUID,
		PartUuids: req.PartUuids,
		PaymentMethod: order_v1.PaymentMethodUNKNOWN,
		Status: order_v1.OrderStatusPENDINGPAYMENT,
	}

	h.storage.PostOrder(&newOrder)

	return &order_v1.CreateOrderResponse{
		OrderUUID: orderUUID,
		}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)

	if order == nil {
		return &order_v1.NotFoundError{
			Code: 404,
			Message: "Order was not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	transactionalUUID := uuid.New()
	
	err := h.storage.PayOrder(params.OrderUUID, req.GetPaymentMethod(), transactionalUUID)


	if err != nil {
		return &order_v1.NotFoundError{
			Code: 404,
			Message: "Order was not found",
		}, nil
	}

	return &order_v1.PayOrderResponse{TransactionUUID: transactionalUUID}, nil
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code: order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}

func main() {
	// Создаем хранилище для данных о погоде
	storage := NewOrderStorage()

	// Создаем обработчик API погоды
	orderHandler := NewOrderHandler(storage)

	// Создаем OpenAPI сервер
	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера заказов: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, 
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер заказов запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера заказов: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
