package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"stream2/proto/uploadpb"

	"google.golang.org/grpc"
)

type uploadServer struct {
	uploadpb.UnimplementedFileUploadServiceServer
}

func (s *uploadServer) Upload(stream uploadpb.FileUploadService_UploadServer) error {
	var fileName string
	var fileData []byte

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// Save the received file
			if err := os.WriteFile(fileName, fileData, 0644); err != nil {
				return stream.SendAndClose(&uploadpb.UploadStatus{
					Success: false,
					Message: fmt.Sprintf("Failed to save file: %v", err),
				})
			}
			log.Println("File uploaded:", fileName)
			return stream.SendAndClose(&uploadpb.UploadStatus{
				Success: true,
				Message: fmt.Sprintf("Uploaded %s successfully", fileName),
			})
		}
		if err != nil {
			return err
		}
		fileName = chunk.Filename
		fileData = append(fileData, chunk.Content...)
	}
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	uploadpb.RegisterFileUploadServiceServer(grpcServer, &uploadServer{})
	log.Println("File upload server running: 50051")
	log.Fatal(grpcServer.Serve(lis))
}
