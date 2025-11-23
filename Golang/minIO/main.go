package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// Load environment variables from .env if exists
	_ = godotenv.Load()

	// MinIO configuration
	endpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	accessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	bucket := getEnv("MINIO_BUCKET", "streamhub")
	useSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalf("Failed to check bucket: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
		fmt.Printf("Created bucket %s\n", bucket)
	}

	// Fiber app
app := fiber.New(fiber.Config{
    BodyLimit: 5 * 1024 * 1024 * 1024, // 5 GB
})

app.Post("/upload", func(c *fiber.Ctx) error {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).SendString("Missing file: " + err.Error())
    }

    file, err := fileHeader.Open()
    if err != nil {
        return c.Status(500).SendString(err.Error())
    }
    defer file.Close()

    objectName := c.FormValue("path")
    if objectName == "" {
        objectName = "uploads/" + fileHeader.Filename
    }

    size := fileHeader.Size
    contentType := fileHeader.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    // Upload directly to MinIO
    info, err := minioClient.PutObject(ctx, bucket, objectName, file, size, minio.PutObjectOptions{
        ContentType: contentType,
    })
    if err != nil {
        return c.Status(500).SendString("Upload failed: " + err.Error())
    }

    return c.JSON(fiber.Map{
        "ok":          true,
        "object_name": objectName,
        "size":        info.Size,
    })
})


	// Health check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})



	// Generate presigned URL
	app.Get("/signed-url", func(c *fiber.Ctx) error {
		object := c.Query("object")
		if object == "" {
			return c.Status(400).SendString("object required")
		}

		expiry := int64(600) // default 10 minutes
		if q := c.Query("expiry"); q != "" {
			var v int64
			fmt.Sscan(q, &v)
			if v > 0 {
				expiry = v
			}
		}

		url, err := minioClient.PresignedGetObject(ctx, bucket, object, time.Duration(expiry)*time.Second, nil)
		if err != nil {
			return c.Status(500).SendString("Failed to generate URL: " + err.Error())
		}

		return c.JSON(fiber.Map{"url": url.String()})
	})

	// Start server
	port := getEnv("HTTP_PORT", "8080")
	log.Printf("Fiber server running on :%s", port)
	app.Listen(":" + port)
}

// Helper to get env with default
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
