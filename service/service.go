package service

import (
	"bytes"
	"fmt"
	"image"
	"image-server/storage"
	"io"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"  // WebP encoder
	_ "golang.org/x/image/webp" // WebP decoder
)

type ImageMetadata struct {
	ID       string
	Filename string
	Format   string
}

type ImageService struct {
	storage storage.ImageStorage
}

func NewImageService(storage storage.ImageStorage) *ImageService {
	return &ImageService{
		storage: storage,
	}
}

func (s *ImageService) GetImage(entityType string, uuid string, fileName string) (io.ReadCloser, error) {
	return s.storage.GetImage(entityType, uuid, fileName)
}

func (s *ImageService) SaveImage(entityType string, uuid string, imageData io.Reader) (*ImageMetadata, error) {
	// Read the image data
	imgBytes, err := io.ReadAll(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Add this to help debug format detection
	format := http.DetectContentType(imgBytes)

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (format: %s): %w", format, err)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 85}); err != nil {
		return nil, fmt.Errorf("failed to encode to WebP: %w", err)
	}

	// Save the image with the new structure
	filename := "main.webp"
	_, err = s.storage.SaveImage(entityType, uuid, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	return &ImageMetadata{
		ID:       uuid,
		Filename: filename,
		Format:   "webp",
	}, nil

}

func (s *ImageService) DeleteImage(entityType string, uuid string) error {
	return s.storage.DeleteImage(entityType, uuid)
}
