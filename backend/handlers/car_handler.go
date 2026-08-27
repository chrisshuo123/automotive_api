package handlers

import (
	"automotiveApi/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CarHandler struct {
	DB *gorm.DB
}

func NewCarHandler(db *gorm.DB) *CarHandler {
	return &CarHandler{DB: db}
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 5MB max
	err := r.ParseMultipartForm(5 << 20) // 5MB
	if err != nil {
		http.Error(w, "File too large. Max Size: 5MB", http.StatusBadRequest)
		return
	}

	// Get car ID from URL or form FIRST
	carID := r.URL.Query().Get("id")
	if carID == "" {
		carID = r.FormValue("id")
	}
	if carID == "" {
		http.Error(w, "Car ID Required", http.StatusBadRequest)
		return
	}

	// Get existing car to preserve image if no new file
	var existingCar models.Cars
	if err := h.DB.Where("idCars = ?", carID).First(&existingCar).Error; err != nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Prepare update data (start with existing image)
	updateData := map[string]interface{}{
		"nama_mobil":  r.FormValue("nama_mobil"),
		"horse_power": r.FormValue("horse_power"),
		"idStatus_fk": r.FormValue("idStatus_fk"),
		"imageCar":    existingCar.ImageCar, // Keep existing image by default
	}

	// Try to get file
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		allowedTypes := []string{"image/png", "image/jpg", "image/jpeg", "image/webp"}
		isValid := false
		for _, t := range allowedTypes {
			if contentType == t {
				isValid = true
				break
			}
		}
		if !isValid {
			http.Error(w, "Invalid file type.  Allowed: PNG, JPG, JPEG, WEBP", http.StatusBadRequest)
			return
		}

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String(), ext)
		// fileName = "1734567890_abc123def.jpg" <- Ini yang disimpan ke dalam DB.

		// Create upload directory
		uploadDir := "./frontend/public/img/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}

		// Save file with UNIQUE filename
		dst, err := os.Create(filepath.Join(uploadDir, fileName))
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to copy file", http.StatusInternalServerError)
			return
		}

		// LOG INI UNTUK DEBUG
		fmt.Println("=== DEBUG FILE UPLOAD ===")
		fmt.Printf("Header.Filename: %s\n", header.Filename) // "C:\fakepath\hello-car.jpg"
		fmt.Printf("Ext only: %s\n", ext)                    // ".jpg"
		fmt.Printf("Generated fileName: %s\n", fileName)     // "1734567890_abc123def.jpg"
		fmt.Printf("Saving to DB: %s\n", fileName)           // "1734567890_abc123def.jpg"

		// Update imageCar with new filename (without the Path)
		updateData["imageCar"] = fileName

	} else if err != http.ErrMissingFile {
		http.Error(w, "Failed to get image", http.StatusBadRequest)
		return
	}

	// Update database
	result := h.DB.Model(&models.Cars{}).
		Where("idCars = ?", carID).
		Updates(updateData)

	if result.Error != nil {
		http.Error(w, "Failed to update car", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Car Updated Successfully",
		"image":   updateData["imageCar"].(string),
	})
}
