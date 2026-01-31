package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Uploadfile(c fiber.Ctx, file *multipart.FileHeader, folder string) (string, error) {
	// Validate Extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".png": true, ".jpeg": true, ".gif": true}
	if !allowed[ext] {
		return "", fmt.Errorf("only images are allowed")
	}

	// 1. Tạo folder chung (VD: ./uploads/signatures)
	dirPath := filepath.Join("./uploads", folder)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	// 2. Tạo tên file unique
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(dirPath, filename)

	// 3. Lưu file
	if err := c.SaveFile(file, destPath); err != nil {
		return "", err
	}

	// Return path để lưu vào DB (Nên dùng dấu / để chuẩn web path)
	return fmt.Sprintf("/uploads/%s/%s", folder, filename), nil
}
