package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)
const maxFileSize = 10 << 20 
const uploads = "/uploads/"
const avatar  = "avatar/"
const versionAPI = "/api/v1"

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func HandleFileUpload(ctx *gin.Context, formFieldName, uploadDir string) (string, error) {
	file, err := ctx.FormFile(formFieldName)
	if err != nil {
		return "", fmt.Errorf("failed to get upload file: %s", err)
	}

	if file.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds %dMB limit", maxFileSize >> 20)
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	// Validate file extension
	if !allowedExts[fileExt] {
		return "", fmt.Errorf("only .jpg, .jpeg or .png files are allowed")
	}

	// Generate unique filename
	filename := generateFileNameUpload(filepath.Base(file.Filename), fileExt)

	  // Create upload directory if not exists
	if _, err := createUploadDir(filepath.Join(uploadDir,avatar)); err != nil {
		return "", err
	}

	// Save file
	savePath := getAvatarFilePath(uploadDir, filename)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return generateAvatarURL(ctx,filename), nil
}

func HandleFileDeleted(fileName string, uploadDir string) error {
	if fileName == "" {
		return nil
	}

	fileUrl := getAvatarFilePath(uploadDir, fileName)
	log.Printf("Attempting to delete file: %s", fileUrl)
	fileInfo, err := os.Stat(fileUrl)
    if os.IsNotExist(err) {
        return nil // File doesn't exist is not an error
    }

    if err != nil {
        return fmt.Errorf("failed to stat file: %w", err)
    }

    if fileInfo.IsDir() {
        return fmt.Errorf("path is a directory: %s", fileUrl)
    }

	os.Chmod(fileUrl, 077) // Ensure the file is writable before deletion
    if err := os.Remove(fileUrl); err != nil {
        return fmt.Errorf("failed to delete file: %w", err)
    }

    return nil
}

func generateFileNameUpload(fileName, fileExt string) string {
	nameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, fileExt)
}

func generateAvatarURL(ctx *gin.Context, fileName string) string{
	return ctx.Request.Host+versionAPI+uploads+avatar+fileName
}

func createUploadDir(uploadDir string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}
	return uploadDir, nil
}

func getAvatarFilePath(uploadDir, fileURL string) string {
	fileName := path.Base(fileURL)
	baseURL := filepath.Join(uploadDir,avatar)
	log.Printf("Base URL: %s, File Name: %s", baseURL, fileName)
	return filepath.Join(baseURL, fileName)
}