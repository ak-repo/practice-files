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
    var (
        file    *os.File
        written bool
        fileName string
    )

    for {
        chunk, err := stream.Recv()
        if err == io.EOF {
            if !written {
                return stream.SendAndClose(&uploadpb.UploadStatus{
                    Success: false,
                    Message: "No file data received",
                })
            }

            return stream.SendAndClose(&uploadpb.UploadStatus{
                Success: true,
                Message: fmt.Sprintf("Uploaded %s successfully", fileName),
            })
        }

        if err != nil {
            return err
        }

        if file == nil {
            // First message
            if chunk.Filename == "" {
                return fmt.Errorf("missing filename in first chunk")
            }
            fileName = chunk.Filename

            // Open file for writing immediately
            file, err = os.Create(fileName)
            if err != nil {
                return fmt.Errorf("cannot create file %s: %v", fileName, err)
            }
            defer file.Close()
        }

        if len(chunk.Content) > 0 {
            _, err := file.Write(chunk.Content)
            if err != nil {
                return fmt.Errorf("failed writing to file: %v", err)
            }
            written = true
        }
    }
}


func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    uploadpb.RegisterFileUploadServiceServer(grpcServer, &uploadServer{})

    log.Println("File upload server running on port 50051")

    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
