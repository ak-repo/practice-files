package main

import (
	"context"
	"log"
	pb "micro-gRPC/pkg/proto/orderpb"
	prodpb "micro-gRPC/pkg/proto/productpb"
	"strconv"
	"time"

	"net"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	mu      sync.Mutex
	orders  map[string]*pb.Order
	counter int
	prodCli prodpb.ProductServiceClient // product service client
}

// CreateOrder implementation
func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	// check the product
	ctxProd, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	productResp, err := s.prodCli.ListProducts(ctxProd, &prodpb.Empty{})
	if err != nil {
		return nil, err
	}
	productExists := false

	for _, p := range productResp.Products {
		if p.Id == req.ProductId {
			productExists = true
			break
		}
	}

	if !productExists {
		return &pb.CreateOrderResponse{
			Id:      "",
			Message: "Product doest not exist",
		}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	id := "O" + strconv.Itoa(s.counter)
	order := &pb.Order{
		Id:        id,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	}
	s.orders[id] = order

	return &pb.CreateOrderResponse{
		Id:      id,
		Message: "Order created successfully",
	}, nil
}

// ListOrders implementation
func (s *server) ListOrders(ctx context.Context, req *pb.Empty) (*pb.OrderList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []*pb.Order
	for _, o := range s.orders {
		list = append(list, o)
	}
	return &pb.OrderList{Orders: list}, nil
}

func main() {

	// connect product service
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Product service: %v", err)
	}

	defer conn.Close()
	prodClinet := prodpb.NewProductServiceClient(conn)

	// start order service
	lis, err := net.Listen("tcp", ":50053") // Order service on port 50053
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{
		orders:  make(map[string]*pb.Order),
		counter: 0,
		prodCli: prodClinet,
	})

	log.Println("Order service running on port 50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
