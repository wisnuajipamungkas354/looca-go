package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"looca-go/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	UploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{UploadDir: uploadDir}
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "0000"
	}
	return hex.EncodeToString(bytes)
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	if header.Size > 1*1024*1024 {
		utils.Error(c, http.StatusBadRequest, "Image exceeds 1MB")
		return
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		utils.Error(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	contentType := http.DetectContentType(buf[:n])
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	default:
		utils.Error(c, http.StatusBadRequest, "Invalid image format. Allowed: jpeg, png, webp")
		return
	}

	// Seek back to start
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to process file")
		return
	}

	imagesDir := filepath.Join(h.UploadDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randomHex(4), ext)
	targetPath := filepath.Join(imagesDir, filename)

	out, err := os.Create(targetPath)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to write file")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"url": fmt.Sprintf("/uploads/images/%s", filename),
	})
}

func (h *UploadHandler) UploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	if header.Size > 50*1024*1024 {
		utils.Error(c, http.StatusBadRequest, "Video file size exceeds 50MB limit")
		return
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		utils.Error(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	contentType := http.DetectContentType(buf[:n])
	clientExt := strings.ToLower(filepath.Ext(header.Filename))

	var ext string
	if strings.Contains(contentType, "mp4") || clientExt == ".mp4" {
		ext = ".mp4"
	} else if strings.Contains(contentType, "webm") || clientExt == ".webm" {
		ext = ".webm"
	} else {
		utils.Error(c, http.StatusBadRequest, "Invalid video format. Allowed: mp4, webm")
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to process file")
		return
	}

	videosDir := filepath.Join(h.UploadDir, "videos")
	if err := os.MkdirAll(videosDir, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randomHex(4), ext)
	targetPath := filepath.Join(videosDir, filename)

	out, err := os.Create(targetPath)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to write file")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"url": fmt.Sprintf("/uploads/videos/%s", filename),
	})
}
