package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	chatpb "chat-streaming/proto/chatpb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := chatpb.NewChatServiceClient(conn)

	// Open bidirectional stream
	stream, err := client.ChatStream(context.Background())
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// Goroutine: receive messages from server
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("recv error: %v", err)
				break
			}
			fmt.Printf("\n💬 %s: %s\n", msg.User, msg.Text)
		}
	}()

	// Main thread: send messages
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your name:")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1]

	for {
		fmt.Print("👉 You: ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		if text == "exit" {
			fmt.Println("👋 Bye!")
			break
		}

		stream.Send(&chatpb.ChatMessage{
			User: name,
			Text: text,
		})
		time.Sleep(time.Millisecond * 100)
	}
}
