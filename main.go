package main

import (
	"context"
	"image-server/config"
	"image-server/handler"
	"image-server/service"
	"image-server/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with configured timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
