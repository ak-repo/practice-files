package main

import (
	"context"
	"log"
	"time"

	authpb "micro-gRPC/pkg/proto/authpb"
	orderpb "micro-gRPC/pkg/proto/orderpb"
	productpb "micro-gRPC/pkg/proto/productpb"

	"google.golang.org/grpc"
)

func main() {
	// -------------------
	// 1️⃣ Connect to Auth Service
	// -------------------
	authConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Auth service: %v", err)
	}
	defer authConn.Close()
	authClient := authpb.NewAuthServiceClient(authConn)

	// Register a user
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	registerResp, err := authClient.Register(ctx, &authpb.RegisterRequest{
		Username: "john",
		Password: "pass123",
	})
	if err != nil {
		log.Fatalf("Register error: %v", err)
	}
	log.Println("Register Response:", registerResp.Message, "ID:", registerResp.Id)

	// Login user
	loginResp, err := authClient.Login(ctx, &authpb.LoginRequest{
		Username: "john",
		Password: "pass123",
	})
	if err != nil {
		log.Fatalf("Login error: %v", err)
	}
	log.Println("Login Response:", loginResp.Message, "Token:", loginResp.Token)

	// -------------------
	// 2️⃣ Connect to Product Service
	// -------------------
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Product service: %v", err)
	}
	defer productConn.Close()
	productClient := productpb.NewProductServiceClient(productConn)

	// Add a product
	addProdResp, err := productClient.AddProduct(ctx, &productpb.AddProductRequest{
		Name:  "Laptop",
		Price: "1000",
	})
	if err != nil {
		log.Fatalf("AddProduct error: %v", err)
	}
	log.Println("AddProduct Response:", addProdResp.Message, "ID:", addProdResp.Id)

	// List products
	listProdResp, err := productClient.ListProducts(ctx, &productpb.Empty{})
	if err != nil {
		log.Fatalf("ListProducts error: %v", err)
	}
	log.Println("List of Products:")
	for _, p := range listProdResp.Products {
		log.Printf("ID: %s, Name: %s, Price: %.2f", p.Id, p.Name, p.Price)
	}

	// -------------------
	// 3️⃣ Connect to Order Service
	// -------------------
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Order service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	// Create an order using product ID
	createOrderResp, err := orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		ProductId: addProdResp.Id,
		Quantity:  2,
	})
	if err != nil {
		log.Fatalf("CreateOrder error: %v", err)
	}
	log.Println("CreateOrder Response:", createOrderResp.Message, "Order ID:", createOrderResp.Id)

	// List orders
	listOrderResp, err := orderClient.ListOrders(ctx, &orderpb.Empty{})
	if err != nil {
		log.Fatalf("ListOrders error: %v", err)
	}
	log.Println("List of Orders:")
	for _, o := range listOrderResp.Orders {
		log.Printf("OrderID: %s, ProductID: %s, Quantity: %d", o.Id, o.ProductId, o.Quantity)
	}
}
