package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image-server/handler/middlewares/cors"
	"image-server/service"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ImageHandler struct {
	imageService *service.ImageService
	maxFileSize  int64
}

func NewImageHandler(imageService *service.ImageService, maxFileSize int64) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
		maxFileSize:  maxFileSize,
	}
}

func (h *ImageHandler) RegisterRoutes(mux *http.ServeMux) {
	handler := http.HandlerFunc(h.handleImageRequests)
	mux.Handle("/images/", cors.CorsMiddleware(handler))
}

func (h *ImageHandler) handleImageRequests(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	entityType := parts[1]

	switch r.Method {
	case http.MethodGet:
		h.getImage(w, r, entityType, parts[2], parts[3])
	case http.MethodPost:
		if strings.HasSuffix(r.URL.Path, "/batch") {
			h.uploadMultipleImages(w, r, entityType)
		} else {
			h.uploadImage(w, r, entityType)
		}
	case http.MethodDelete:
		h.deleteImage(w, r, entityType, parts[2])
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ImageHandler) getImage(w http.ResponseWriter, r *http.Request, entityType, uuid string, fileName string) {
	img, err := h.imageService.GetImage(entityType, uuid, fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get image: %v", err), http.StatusNotFound)
		return
	}
	defer img.Close()

	// Detect content type
	buffer := make([]byte, 512)
	_, err = img.(io.ReadSeeker).Read(buffer)
	if err != nil {
		http.Error(w, "failed to detect content type", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buffer)

	// Seek back to start of file
	_, err = img.(io.ReadSeeker).Seek(0, 0)
	if err != nil {
		http.Error(w, "failed to reset file pointer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, "", time.Time{}, img.(io.ReadSeeker))
}

func (h *ImageHandler) uploadImage(w http.ResponseWriter, r *http.Request, entityType string) {
	uuid := uuid.New().String()

	// Parse with the actual max file size limit
	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		http.Error(w, fmt.Sprintf("file size exceeds maximum limit of %d MB", h.maxFileSize/1024/1024), http.StatusRequestEntityTooLarge)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "failed to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file to check size
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	// Check if single file exceeds limit
	if int64(len(fileBytes)) > h.maxFileSize {
		http.Error(w, fmt.Sprintf("file size %d MB exceeds maximum limit of %d MB", len(fileBytes)/1024/1024, h.maxFileSize/1024/1024), http.StatusRequestEntityTooLarge)
		return
	}

	fileReader := bytes.NewReader(fileBytes)
	metadata, err := h.imageService.SaveImage(entityType, uuid, fileReader)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save image: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response with metadata
	response := map[string]interface{}{
		"id":       metadata.ID,
		"filename": metadata.Filename,
		"format":   metadata.Format,
		"url":      fmt.Sprintf("/images/%s/%s/%s", entityType, metadata.ID, metadata.Filename),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ImageHandler) uploadMultipleImages(w http.ResponseWriter, r *http.Request, entityType string) {
	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		http.Error(w, fmt.Sprintf("upload size exceeds maximum limit of %d MB", h.maxFileSize/1024/1024), http.StatusRequestEntityTooLarge)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "no images provided", http.StatusBadRequest)
		return
	}

	var responses []map[string]interface{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to open file %s: %v", fileHeader.Filename, err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read and check each file's size individually
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}

		// Check if individual file exceeds limit
		if int64(len(fileBytes)) > h.maxFileSize {
			http.Error(w, fmt.Sprintf("file %s size %d MB exceeds maximum limit of %d MB",
				fileHeader.Filename, len(fileBytes)/1024/1024, h.maxFileSize/1024/1024),
				http.StatusRequestEntityTooLarge)
			return
		}

		uuid := uuid.New().String()
		fileReader := bytes.NewReader(fileBytes)
		metadata, err := h.imageService.SaveImage(entityType, uuid, fileReader)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to save image %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
			return
		}

		responses = append(responses, map[string]interface{}{
			"id":       metadata.ID,
			"filename": metadata.Filename,
			"format":   metadata.Format,
			"url":      fmt.Sprintf("/images/%s/%s/%s", entityType, metadata.ID, metadata.Filename),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses)
}

func (h *ImageHandler) deleteImage(w http.ResponseWriter, r *http.Request, entityType, uuid string) {
	if err := h.imageService.DeleteImage(entityType, uuid); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete image: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
