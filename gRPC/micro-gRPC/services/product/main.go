package main

import (
	"context"
	"log"
	pb "micro-gRPC/pkg/proto/productpb"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProductServiceServer
	mu       sync.Mutex
	products map[string]*pb.Product
	counter  int
}

func (s *server) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	id := "p" + string(s.counter)
	product := &pb.Product{
		Id:    id,
		Name:  req.Name,
		Price: req.Price,
	}

	s.products[id] = product

	return &pb.AddProductResponse{
		Id:      id,
		Message: "Product added successfully",
	}, nil
}

func (s *server) ListProducts(ctx context.Context, req *pb.Empty) (*pb.ProductList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []*pb.Product
	for _, p := range s.products {
		list = append(list, p)
	}

	return &pb.ProductList{Products: list}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")

	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{
		products: make(map[string]*pb.Product),
		counter:  0,
	})

	log.Println("Product service running on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
