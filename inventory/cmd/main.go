package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	"slices"
	"sync"
	"syscall"

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
		if (slices.Contains(req.GetFilter().GetUuids(), v.Uuid) || len(req.GetFilter().GetUuids()) == 0) && 
		   (slices.Contains(req.GetFilter().GetNames(), v.Name) || len(req.GetFilter().GetNames()) == 0) && 
		   (slices.Contains(req.GetFilter().GetCategories(), v.Category) || len(req.GetFilter().GetCategories()) == 0) && 
		   (slices.Contains(req.GetFilter().GetManufacturerCountries(), v.Manufacturer.GetCountry()) || len(req.GetFilter().GetManufacturerCountries()) == 0) && 
		   slices.Equal(req.GetFilter().GetTags(), v.Tags) {
			result = append(result, v)
		}
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
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

	s := grpc.NewServer()

	service := &inventoryService{
		parts: make(map[string]*inventory_v1.Part),
	}

	inventory_v1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
