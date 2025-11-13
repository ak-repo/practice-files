package main

import (
	"context"
	"io"
	"log"
	"os"
	"stream2/proto/uploadpb"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()
	client := uploadpb.NewFileUploadServiceClient(conn)

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.Open("sample.txt")
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		stream.Send(&uploadpb.FileChunk{
			Filename: "uploaded_sample.txt",
			Content:  buf[:n],
		})
		time.Sleep(500 * time.Millisecond)
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Upload result:", status.Message)
}
