package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const VIDEO_DIR = "./videos/"

func main() {
	r := gin.Default()

	// Ensure folder exists
	os.MkdirAll(VIDEO_DIR, os.ModePerm)

	r.POST("/upload", uploadVideo)
	r.GET("/stream/:filename", streamVideo)

	r.Run(":8080")
}

// ------------------------
// 1. UPLOAD VIDEO
// ------------------------
func uploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(400, gin.H{"error": "No video file"})
		return
	}

	savePath := filepath.Join(VIDEO_DIR, file.Filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save"})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Uploaded successfully",
		"filename": file.Filename,
	})
}

// ------------------------
// 2. STREAM VIDEO
// ------------------------
func streamVideo(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(VIDEO_DIR, filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	// -------- HANDLE RANGE HEADER --------
	rangeHeader := c.GetHeader("Range")

	if rangeHeader == "" {
		// If no range header → send full file
		c.Header("Content-Type", "video/mp4")
		c.Header("Content-Length", fmt.Sprintf("%d", fileSize))
		c.Header("Accept-Ranges", "bytes")
		io.Copy(c.Writer, file)
		return
	}

	// Extract bytes from: "bytes=0-"
	rangeParts := strings.Split(strings.Replace(rangeHeader, "bytes=", "", 1), "-")
	start, _ := strconv.ParseInt(rangeParts[0], 10, 64)

	var end int64
	if rangeParts[1] != "" {
		end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
	} else {
		end = fileSize - 1
	}

	if start > end || end >= fileSize {
		c.JSON(416, gin.H{"error": "Invalid range"})
		return
	}

	chunkSize := (end - start) + 1

	// Seek the position
	file.Seek(start, 0)

	c.Status(http.StatusPartialContent)
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Length", fmt.Sprintf("%d", chunkSize))
	c.Header("Content-Range",
		fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Accept-Ranges", "bytes")

	// Stream chunk
	io.CopyN(c.Writer, file, chunkSize)
}
