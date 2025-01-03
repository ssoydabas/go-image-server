package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ImageStorage interface {
	SaveImage(entityType string, uuid string, imageData io.Reader) (string, error)
	GetImage(entityType string, uuid string, fileName string) (io.ReadCloser, error)
	DeleteImage(entityType string, uuid string) error
}

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &FileStorage{basePath: basePath}, nil
}

func (fs *FileStorage) SaveImage(entityType string, uuid string, imageData io.Reader) (string, error) {
	// Create directory path: basePath/entityType/uuid/
	dirPath := filepath.Join(fs.basePath, entityType, uuid)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create image directory: %w", err)
	}

	// Create file path: basePath/entityType/uuid/main.webp
	filename := "main.webp"
	filepath := filepath.Join(dirPath, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, imageData); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return filename, nil
}

func (fs *FileStorage) GetImage(entityType string, uuid string, fileName string) (io.ReadCloser, error) {
	absBasePath, err := filepath.Abs(fs.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve base path: %w", err)
	}

	entityType = filepath.Clean(entityType)
	uuid = filepath.Clean(uuid)
	fullPath := filepath.Join(fs.basePath, entityType, uuid, fileName)
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve full path: %w", err)
	}

	// Check if directory exists
	if _, err := os.Stat(filepath.Dir(absFullPath)); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", filepath.Dir(absFullPath))
	}

	// Security check with resolved absolute paths
	if !strings.HasPrefix(absFullPath, absBasePath) {
		return nil, fmt.Errorf("security error: attempted path traversal detected")
	}

	file, err := os.Open(absFullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	return file, nil
}

func (fs *FileStorage) DeleteImage(entityType string, entityID string) error {
	formats := []string{"jpg", "jpeg", "png", "webp"}
	var lastErr error
	deleted := false

	for _, format := range formats {
		filepath := filepath.Join(fs.basePath, entityType, fmt.Sprintf("%s.%s", entityID, format))
		if err := os.Remove(filepath); err == nil {
			deleted = true
		} else {
			lastErr = err
		}
	}

	if !deleted {
		if lastErr != nil {
			return fmt.Errorf("failed to delete image: %w", lastErr)
		}
		return fmt.Errorf("image not found")
	}
	return nil
}
