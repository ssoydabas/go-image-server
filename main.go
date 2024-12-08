package main

import (
	"image-server/config"
	"image-server/handler"
	"image-server/service"
	"image-server/storage"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fileStorage, err := storage.NewFileStorage(cfg.Storage.BasePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize service
	imageService := service.NewImageService(fileStorage)

	// Initialize handler
	imageHandler := handler.NewImageHandler(imageService, cfg.Server.MaxFileSize)

	// Setup router
	mux := http.NewServeMux()
	imageHandler.RegisterRoutes(mux)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
