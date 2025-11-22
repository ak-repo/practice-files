package main

import (
	"context"
	"file/api/filepb"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
)

type FileHandler struct {
	GRPC filepb.FileServiceClient
}

func (h *FileHandler) Upload(c *fiber.Ctx) error {
	// Read file from multipart request
	file, _ := c.FormFile("file")
	f, _ := file.Open()
	defer f.Close()

	stream, _ := h.GRPC.Upload(context.Background())

	// 1️⃣ Send metadata
	stream.Send(&filepb.UploadRequest{
		Data: &filepb.UploadRequest_Metadata{
			Metadata: &filepb.UploadMetadata{
				Filename: file.Filename,
				Filesize: file.Size,
				Mimetype: file.Header.Get("Content-Type"),
				UserId:   "123",
			},
		},
	})

	// 2️⃣ Send chunks
	buf := make([]byte, 1024*32)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			stream.Send(&filepb.UploadRequest{
				Data: &filepb.UploadRequest_Chunk{
					Chunk: &filepb.UploadChunk{
						Content: buf[:n],
					},
				},
			})
		}
		if err == io.EOF {
			break
		}
	}

	// 3️⃣ Receive response
	res, _ := stream.CloseAndRecv()

	return c.JSON(res)
}

func (h *FileHandler) Download(c *fiber.Ctx) error {
	fileID := c.Params("id")

	stream, _ := h.GRPC.Download(context.Background(), &filepb.DownloadRequest{
		FileId: fileID,
	})

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileID))

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		c.Write(chunk.Content)
	}

	return nil
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100 MB
	})

	// 1️⃣ FIX CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Content-Type, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	}))

	// 3️⃣ gRPC client
	conn, _ := grpc.Dial("localhost:50052", grpc.WithInsecure())
	grpcClient := filepb.NewFileServiceClient(conn)

	fileHandler := &FileHandler{GRPC: grpcClient}

	// 4️⃣ Routes
	app.Post("/api/v1/files/upload", fileHandler.Upload)
	app.Get("/api/v1/files/:id", fileHandler.Download)

	app.Listen(":8080")
}
