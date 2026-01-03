package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	"slices"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer

	mu sync.RWMutex
	parts map[string]*inventory_v1.Part
}

func (inv *inventoryService) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	part, ok := inv.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventory_v1.GetPartResponse{Part: part}, nil
}

func (inv *inventoryService) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	
	result := make([]*inventory_v1.Part, 0)

	if len(req.GetFilter().GetUuids()) == 0 && 
	   len(req.GetFilter().GetNames()) == 0 && 
	   len(req.GetFilter().GetCategories()) == 0 && 
	   len(req.GetFilter().GetManufacturerCountries()) == 0 && 
	   len(req.GetFilter().GetTags()) == 0 {
		for _, v := range inv.parts {
			result = append(result, v)
		}
		
		return &inventory_v1.ListPartsResponse{Parts: result}, nil
	}

	for _, v := range inv.parts {
		if isInFilters(req.GetFilter(), v) {
			result = append(result, v)
		}
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
}

func isInFilters(filters *inventory_v1.PartsFilter, part *inventory_v1.Part) bool {
	if  (slices.Contains(filters.GetUuids(), part.Uuid) || len(filters.GetUuids()) == 0) && 
		(slices.Contains(filters.GetNames(), part.Name) || len(filters.GetNames()) == 0) && 
		(slices.Contains(filters.GetCategories(), part.Category) || len(filters.GetCategories()) == 0) && 
		(slices.Contains(filters.GetManufacturerCountries(), part.Manufacturer.GetCountry()) || len(filters.GetManufacturerCountries()) == 0) && 
		isInTags(filters.GetTags(), part.Tags) {
			return true
		}
	return false 
}

func isInTags(filterTags []string, partTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, filterTag := range partTags {
		if !slices.Contains(partTags, filterTag) {
			return false
		}
	}
	return true 
}

// LoggerInterceptor создает серверный унарный интерцептор, который логирует
// информацию о времени выполнения методов gRPC сервера.
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

	service := &inventoryService{
		parts: make(map[string]*inventory_v1.Part),
	}

	service.parts["123e4567-e89b-12d3-a456-426614174000"] = &inventory_v1.Part{
		Uuid:          "123e4567-e89b-12d3-a456-426614174000",
		Name:          "Тормозной диск передний",
		Description:   "Высококачественный вентилируемый тормозной диск для передних колес",
		Price:         12499.99,
		StockQuantity: 25,
		Category:      inventory_v1.Category_CATEGORY_PORTHOLE,
		Tags:      []string{"тормоза", "диск", "передний", "вентилируемый"},
	}

	inventory_v1.RegisterInventoryServiceServer(s, service)

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
