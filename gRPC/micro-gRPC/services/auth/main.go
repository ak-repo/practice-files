package main

import (
	"context"
	"log"
	pb "micro-gRPC/pkg/proto/authpb"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
	users map[string]string
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id := req.Username + "_id"
	s.users[req.Username] = req.Password
	return &pb.RegisterResponse{
		Id:      id,
		Message: "User registered",
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	pass, ok := s.users[req.Username]
	if !ok || pass != req.Password {
		return &pb.LoginResponse{Message: "Invalid credentials"}, nil
	}

	token := "token_" + req.Username
	return &pb.LoginResponse{
		Token:   token,
		Message: "Login successful",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{users: make(map[string]string)}
	pb.RegisterAuthServiceServer(s, srv)

	log.Println("Auth service running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server : %v", err)
	}

}
