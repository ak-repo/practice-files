package main

import (
	"context"
	"log"
	"time"

	pb "basics/gRPC/calculatorpb" // replace with your module path

	"google.golang.org/grpc"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	// Call Add method
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Add(ctx, &pb.AddRequest{A: 5, B: 7})
	if err != nil {
		log.Fatalf("error calling Add: %v", err)
	}

	log.Printf("Result: %d", resp.Result)
}
