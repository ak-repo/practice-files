package main

import (
	"file/api/filepb"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type FileServer struct {
	filepb.UnimplementedFileServiceServer
}

func (s *FileServer) Upload(stream filepb.FileService_UploadServer) error {
	var meta *filepb.UploadMetadata

	// 1️⃣ Metadata
	req, _ := stream.Recv()
	meta = req.GetMetadata()

	fileID := uuid.NewString()
	savePath := filepath.Join("uploads", fileID+"_"+meta.Filename)

	f, _ := os.Create(savePath)
	defer f.Close()

	// 2️⃣ Receive chunks
	for {
		chunkReq, err := stream.Recv()
		if err == io.EOF {
			break
		}

		chunk := chunkReq.GetChunk()
		f.Write(chunk.Content)
	}

	// 3️⃣ Response
	return stream.SendAndClose(&filepb.UploadResponse{
		FileId:      fileID,
		Filename:    meta.Filename,
		Filesize:    meta.Filesize,
		DownloadUrl: "/api/v1/files/" + fileID,
	})
}

func (s *FileServer) Download(req *filepb.DownloadRequest, stream filepb.FileService_DownloadServer) error {
	path := filepath.Join("uploads", req.FileId+"_*") // handle wildcard
	matches, _ := filepath.Glob(path)
	if len(matches) == 0 {
		return nil
	}

	f, _ := os.Open(matches[0])
	defer f.Close()

	buf := make([]byte, 1024*64)
	part := int64(1)

	for {
		n, err := f.Read(buf)
		if n > 0 {
			stream.Send(&filepb.DownloadChunk{
				Content: buf[:n],
				Part:    part,
			})
			part++
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}

func main() {
	os.MkdirAll("uploads", os.ModePerm)

	grpcServer := grpc.NewServer()
	filepb.RegisterFileServiceServer(grpcServer, &FileServer{})

	lis, _ := net.Listen("tcp", ":50052")
	grpcServer.Serve(lis)
}
