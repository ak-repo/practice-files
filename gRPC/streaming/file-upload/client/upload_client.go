package main

import (
	"context"
	"io"
	"log"
	"os"
	"stream2/proto/uploadpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := uploadpb.NewFileUploadServiceClient(conn)

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("Failed to create upload stream: %v", err)
	}

	file, err := os.Open("sample.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, 1024)

	// Read and send first chunk with filename
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatalf("File read error: %v", err)
	}

	if n > 0 {
		err = stream.Send(&uploadpb.FileChunk{
			Filename: "uploaded_sample.txt",
			Content:  buf[:n],
		})
		if err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
		}
	}

	// Continue sending remaining chunks WITHOUT filename
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("File read error: %v", err)
		}

		err = stream.Send(&uploadpb.FileChunk{
			Content: buf[:n],
		})
		if err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
		}
	}

	// Finish upload
	status, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	log.Println("Upload result:", status.Message)
}
