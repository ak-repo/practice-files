package main

import (
	"context"
	"log"
	"net"

	pb "basics/gRPC/calculatorpb" // replace with your module path

	"google.golang.org/grpc"
)

// server struct implements CalculatorServer interface
type server struct {
	pb.UnimplementedCalculatorServer
}

// Add method implementation
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	result := req.A + req.B
	return &pb.AddResponse{Result: result}, nil
}

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})

	log.Println("gRPC server listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
