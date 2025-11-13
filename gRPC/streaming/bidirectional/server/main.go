package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	chatpb "chat-streaming/proto/chatpb"
)

type ChatServer struct {
	chatpb.UnimplementedChatServiceServer
}

func (s *ChatServer) ChatStream(stream chatpb.ChatService_ChatStreamServer) error {
	log.Println("💬 New client connected")

	for {
		// Receive message from client
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("client disconnected")
			return nil
		}
		if err != nil {
			log.Println("recv error:", err)
			return err
		}

		log.Printf("📨 Received from %s: %s", msg.User, msg.Text)

		// Echo back with prefix
		reply := &chatpb.ChatMessage{
			User: "Server",
			Text: fmt.Sprintf("Hi %s, got your message: %s", msg.User, msg.Text),
		}
		if err := stream.Send(reply); err != nil {
			log.Println("send error:", err)
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chatpb.RegisterChatServiceServer(s, &ChatServer{})
	log.Println("🚀 Chat server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
