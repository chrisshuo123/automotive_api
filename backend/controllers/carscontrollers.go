package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Base folder where car images are saved and served from.
// Relative to the Go process's working directory.
var carImageDir string

func init() {
	dir := os.Getenv("CAR_IMAGE_DIR")
	if dir == "" {
		// Fallback for local dev - adjust to taste
		dir = "../frontend/public/img"
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	carImageDir = abs
	fmt.Println("Car image dir: ", carImageDir)
}

func saveUploadedImageWithName(fileHeader *multipart.FileHeader, fileName string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Cari root project dari working directory saat ini
	wd, _ := os.Getwd()
	// File Path Upload Directory
	uploadDir := filepath.Join(wd, "..", "frontend", "public", "img")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("Failed to create upload dir: %w", err)
	}

	dstPath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// deleteCarImage removes a previously saved car image from disk.
// Returns nil if the file is already gone (idempotent)
func deleteCarImage(fileName string) error {
	if fileName == "" {
		return nil
	}

	// Safety: only allow a bare filename, no path seperators or ".."
	clean := filepath.Clean(fileName)
	if clean != fileName || strings.ContainsAny(fileName, `/\`) || fileName == "." || fileName == ".." {
		return fmt.Errorf("refusing to delete suspicious filename: %q", fileName)
	}

	fullPath := filepath.Join(carImageDir, clean)

	abs, _ := filepath.Abs(fullPath)
	fmt.Println("[deletecarImage] attempting to remove: ", abs)

	if err := os.Remove(fullPath); err != nil {
		fmt.Println("[deleteCarImage] os.Remove error: ", err)
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	fmt.Println("[deleteCarImage] removed OK")
	return nil
}

func CreateCarsController(c echo.Context) error {
	var carsRequest models.Cars

	carsRequest.NamaMobil = c.FormValue("nama_mobil")

	if hp, err := strconv.Atoi(c.FormValue("horse_power")); err == nil {
		carsRequest.HorsePower = uint(hp)
	}

	if v := c.FormValue("idmerek_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			carsRequest.MerekID = &uid
		}
	}
	if v := c.FormValue("idjenis_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			carsRequest.JenisID = &uid
		}
	}
	if v := c.FormValue("idstatus_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			carsRequest.StatusID = &uid
		}
	}

	// Validate foreign keys exits
	if carsRequest.MerekID != nil {
		var merek models.Merek
		if err := configs.DB.First(&merek, *carsRequest.MerekID).Error; err != nil {
			return c.JSON(400, models.BaseResponse{
				Message: "Invalid brand ID",
				Status:  false,
			})
		}
	}
	if carsRequest.JenisID != nil {
		var jenis models.Jenis
		if err := configs.DB.First(&jenis, *carsRequest.JenisID).Error; err != nil {
			return c.JSON(400, models.BaseResponse{
				Message: "Invalid jenis ID",
				Status:  false,
			})
		}
	}
	if carsRequest.StatusID != nil {
		var status models.Status
		if err := configs.DB.First(&status, *carsRequest.StatusID).Error; err != nil {
			return c.JSON(400, models.BaseResponse{
				Message: "Invalid Status ID",
				Status:  false,
			})
		}
	}

	// ===== DEBUG: dump every multipart key the server received =====
	if form, err := c.MultipartForm(); err == nil {
		for k, v := range form.Value {
			fmt.Printf("[CREATE] value key: %q = %v\n", k, v)
		}
		for k, v := range form.File {
			for _, fh := range v {
				fmt.Printf("[CREATE] file key: %q => filename=%q size=%d\n", k, fh.Filename, fh.Size)
			}
		}
	} else {
		fmt.Println("[CREATE] MultipartForm err: ", err)
	}
	// ===============================================================

	// Optional image upload
	fileHeader, err := c.FormFile("image")
	fmt.Println("FormFile err: ", err)
	if err == nil {
		fmt.Println("Got file: ", fileHeader.Filename)

		// Ambil extension dari file asli (mis. ".jpg")
		ext := filepath.Ext(fileHeader.Filename)

		// Generate 7 byte random: "car_123124124qsafwr" jadi 14 karakter hex
		b := make([]byte, 7) // 7 byte x 2 = 14 char hex
		if _, err := rand.Read(b); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: "Failed to generate random name",
				Status:  false,
			})
		}

		// Gabungkan: car_ + hex + ext
		fileName := "car_" + hex.EncodeToString(b) + ext
		// Hasil: car_6a9bdbc37fbfb.jpg

		fmt.Println("Generated filename: ", fileName)

		// Simpan file dengan nama baru
		if err := saveUploadedImageWithName(fileHeader, fileName); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: err.Error(),
				Status:  false,
			})
		}
		carsRequest.ImageCar = fileName
	}

	result := configs.DB.Create(&carsRequest) // Cars Request
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: result.Error.Error(), Status: false,
		})
	}

	// Reload the car with relationships
	var newCar models.Cars
	configs.DB.
		Preload("Merek").
		Preload("Jenis").
		Preload("Status").
		First(&newCar, carsRequest.CarsID)

	//return c.JSON(http.StatusOK, carsRequest) // Cars Request
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menambahkan data",
		Status:  true,
		Data:    newCar,
	})
}

func UpdateCarController(c echo.Context) error {
	fmt.Println("[UPDATE] handler entered, id = ", c.Param("id"))
	// Get ID from URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: "Invalid ID", Status: false,
		})
	}

	// 1. Get Existing Cars, check if cars exists
	var car models.Cars
	if err := configs.DB.First(&car, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, models.BaseResponse{
			Message: "Car not found", Status: false,
		})
	}

	if v := c.FormValue("nama_mobil"); v != "" {
		car.NamaMobil = v
	}
	if v := c.FormValue("horse_power"); v != "" {
		if hp, err := strconv.Atoi(v); err == nil {
			car.HorsePower = uint(hp)
		}
	}
	if v := c.FormValue("idmerek_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			car.MerekID = &uid
		}
	}
	if v := c.FormValue("idjenis_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			car.JenisID = &uid
		}
	}
	if v := c.FormValue("idstatus_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			car.StatusID = &uid
		}
	}

	var oldFileName string
	newFileUploaded := false

	// Only replace the image if a new file was actually sent
	fileHeader, err := c.FormFile("image") // "image" is actually "nama_foto", taken from handlers.js
	fmt.Println("[UPDATE] FormFile err: ", err)

	if err == nil {
		fmt.Println("Got file: ", fileHeader.Filename)

		// Remember the old filename BEFORE we overwrite it
		oldFileName = car.ImageCar

		// Ambil extension dari file asli (misal: ".jpg .png, dkk...")
		ext := filepath.Ext(fileHeader.Filename)

		// Generate 7 byte random: "car_1234123412qwqrs" jadi 14 karakter hex
		b := make([]byte, 7) // 7 byte x 2 = 14 char hex
		if _, err := rand.Read(b); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: "Failed to generate random name",
				Status:  false,
			})
		}

		// Gabungkan: car_ + hex + ext
		fileName := "car_" + hex.EncodeToString(b) + ext
		// Hasil: car_6a9bdbc37fbfb.jpg

		fmt.Println("Generated filename: ", fileName)

		// Simpan file dengan nama baru
		if err := saveUploadedImageWithName(fileHeader, fileName); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: err.Error(),
				Status:  false,
			})
		}

		// Point the row at the new file
		car.ImageCar = fileName
		newFileUploaded = true

		// Save to DB
		if err := configs.DB.Save(&car).Error; err != nil {
			if newFileUploaded {
				// Rollback the new file so we don't leave an orphan
				_ = deleteCarImage(fileName) // roll back the new file
				car.ImageCar = oldFileName   // (optional) restore in-memory value
			}
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: err.Error(), Status: false,
			})
		}

		if newFileUploaded && oldFileName != "" {
			fmt.Println("[UPDATE] deleting old file: ", oldFileName)
			// DB save succeeded -> now it's safe to delete the old file
			if err := deleteCarImage(oldFileName); err != nil {
				// Non-fatal: Log it, but don't fail the request
				fmt.Printf("warning: failed to delete old image %q: %v\n", oldFileName, err)
			}
		}
	} else {
		fmt.Println("[UPDATE] no new file uploaded - reason: ", err)
		// No new file - just save any other field changes

		form, mErr := c.MultipartForm()
		fmt.Println("[UPDATE] MultipartForm err: ", mErr)
		if mErr == nil {
			fmt.Println("[UPDATE] --- form.Value keys ---")
			for k, v := range form.Value {
				fmt.Printf("[UPDATE] value key: %q = %v\n", k, v)
			}
			fmt.Println("[UPDATE] --- form.File keys ---")
			for k, v := range form.File {
				for _, fh := range v {
					fmt.Printf("[UPDATE] file key: %q -> filename=%q size=%d\n", k, fh.Filename, fh.Size)
				}
			}
			if len(form.File) == 0 {
				fmt.Println("[UPDATE]	(no file parts at all)")
			}
		} else {
			fmt.Println("[UPDATE] MultipartForm err: ", mErr)
		}

		if err := configs.DB.Save(&car).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: err.Error(), Status: false,
			})
		}
	}

	// Force Reload with relationships
	var updatedCar models.Cars
	if err := configs.DB.
		Preload("Merek").
		Preload("Jenis").
		Preload("Status").
		First(&updatedCar, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: err.Error(), Status: false,
		})
	}

	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Car updated successfully", Status: true, Data: updatedCar,
	})
}

// saveUploadedImage validates and saves an uploaded image file, returning the generated filename (not the full path) to store in the DB.
func saveUploadedImage(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file")
	}
	defer src.Close()

	// Randomize the File Name
	// ext := filepath.Ext(fileHeader.Filename)
	// fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// Want to keep the actual file name? here:
	fileName := fileHeader.Filename

	uploadDir := "../frontend/public/img"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory")
	}

	dst, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		return "", fmt.Errorf("failed to save file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to write file")
	}

	return fileName, nil
}

func GetCarsController(c echo.Context) error {
	var count int64
	configs.DB.Table("cars").Count(&count)
	fmt.Println("Total rows in cars table: ", count)

	var cars []models.Cars

	// Set Response Header
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

	// Debug GORM Error output Mitigation
	fmt.Println("GORM is using table:", configs.DB.NamingStrategy.TableName("Merek"))

	// First load cars with relationships
	result := configs.DB.
		Preload("Merek", func(db *gorm.DB) *gorm.DB {
			return db.Select("idmerek, namamerek") // Only load necessary fields
		}).
		Preload("Jenis", func(db *gorm.DB) *gorm.DB {
			return db.Select("idjenis, namajenis")
		}).
		Preload("Status", func(db *gorm.DB) *gorm.DB {
			return db.Select("idstatus, namastatus")
		}).
		Find(&cars)
		//First(&cars, id)

	/*if result.Error != nil {
		return c.JSON(500, result.Error.Error())
	}*/

	if result.Error != nil {
		//return c.JSON(500, result.Error.Error())
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: "Failed to load cars: " + result.Error.Error(),
			Status:  false,
			Data:    nil,
		})
	}

	// Debug: Check what's actually being loaded
	for i, car := range cars {
		// Recently Merek Column only Shows null in the Postman API:
		if car.Merek == nil && car.MerekID != nil {
			fmt.Printf("Car %d: MerekID Exists (%d) but Merek is\n", i, *car.MerekID)
		}
	}

	// Verify and Clean null relationships
	for i := range cars {
		if cars[i].Merek == nil && cars[i].MerekID != nil {
			var merek models.Merek
			if err := configs.DB.First(&merek, *cars[i].MerekID).Error; err == nil {
				cars[i].Merek = &merek
			}
		}

		if cars[i].Jenis == nil && cars[i].JenisID != nil {
			var jenis models.Jenis
			if err := configs.DB.First(&jenis, *cars[i].JenisID).Error; err == nil {
				cars[i].Jenis = &jenis
			}
		}

		if cars[i].Status == nil && cars[i].StatusID != nil {
			var status models.Status
			if err := configs.DB.First(&status, *cars[i].StatusID).Error; err == nil {
				cars[i].Status = &status
			}
		}
	}

	// Debug: Log the first car's relationships
	if len(cars) > 0 {
		log.Printf("First car relationships - Merek: %+v, Jenis: %+v, Status: %+v",
			cars[0].Merek,
			cars[0].Jenis,
			cars[0].Status)
	}

	//return c.JSON(http.StatusOK, cars)
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menampilkan data",
		Status:  true,
		Data:    cars,
	})
}
func GetCarController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// return c.JSON(400, map[string]string{"error": "Invalid ID Format"})
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: "Invalid ID format",
			Status:  false,
			Data:    nil,
		})
	}

	log.Printf("🔵 GetCarController: id=%d", id)

	// var car []models.Cars // Single record, not slice
	var car models.Cars // Single struct, NOT slice

	// idCars := c.QueryParam("idCars")

	// Build query with Preload and filter by ID
	result := configs.DB.
		Preload("Merek").
		Preload("Jenis").
		Preload("Status").
		Where("idcars = ?", id).
		First(&car) // Use First() for single records

	log.Printf("🔵 Query error: %v", result.Error)
	log.Printf("🔵 Rows affected: %d", result.RowsAffected)
	log.Printf("🔵 Car: %+v", car)

	// result := query.Find(&cars)

	if result.Error != nil {
		log.Println("GetCar error: ", result.Error)
		log.Println("Looking for idcars: ", id)
		// return c.JSON(404, map[string]string{"error": "Car not found"})
		return c.JSON(http.StatusNotFound, models.BaseResponse{
			Message: "Cars not found",
			Status:  false,
			Data:    nil,
		})
	}
	// return c.JSON(200, car)
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menampilkan data car",
		Status:  true,
		Data:    car,
	})
}

func GetMereksController(c echo.Context) error {
	var merek []models.Merek

	// Set Response Header
	//c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

	result := configs.DB.Find(&merek)

	if result.Error != nil {
		//return c.JSON(500, result.Error.Error())
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: "Failed to load merek: " + result.Error.Error(),
			Status:  false,
			Data:    nil,
		})
	}

	//return c.JSON(http.StatusOK, cars)
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menampilkan data merek",
		Status:  true,
		Data:    merek,
	})
}

func GetMerekController(c echo.Context) error {
	// 1. Take id from URL param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// return c.JSON(400, map[string]string{"error": "Invalid ID Format"})
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: "Invalid ID format",
			Status:  false,
			Data:    nil,
		})
	}

	// 2. Query to Database
	var merek models.Merek // Single struct, NOT slice
	result := configs.DB.First(&merek, "idmerek = ?", id)

	// 3. Check whether the data already founded or not yet
	// return c.JSON(200, car)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, models.BaseResponse{
			Message: "Brand Not Founded",
			Status:  true,
			Data:    merek,
		})
	}

	// 4. Return the result
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Successfully shows the Brand data",
		Status:  true,
		Data:    merek,
	})
}

func GetTypesController(c echo.Context) error {
	var jenis []models.Jenis

	result := configs.DB.Find(&jenis)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: "Failed to load jenis: " + result.Error.Error(),
			Status:  false,
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menampilkan data jenis",
		Status:  true,
		Data:    jenis,
	})
}

func GetTypeController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: "Invalid ID format",
			Status:  false,
			Data:    nil,
		})
	}

	var jenis models.Jenis
	result := configs.DB.First(&jenis, "idjenis = ?", id)

	if result.Error != nil {
		return c.JSON(http.StatusNotFound, models.BaseResponse{
			Message: "Type Not Founded",
			Status:  false,
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Successfully shows the Type data",
		Status:  true,
		Data:    jenis,
	})
}

func GetStatusController(c echo.Context) error {
	var status []models.Status

	result := configs.DB.Find(&status)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: "Failed to load status: " + result.Error.Error(),
			Status:  false,
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menampilkan data status",
		Status:  true,
		Data:    status,
	})
}

func DeleteCarController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": false,
			"error":  "Invalid ID Format",
		})
	}

	// Load the car first to know the image filename
	var car models.Cars
	if err := configs.DB.First(&car, id).Error; err != nil {
		return c.JSON(404, map[string]interface{}{
			"status": false,
			"error":  "Car not found",
		})
	}

	// Delete DB row
	result := configs.DB.Delete(&models.Cars{}, id)
	if result.Error != nil {
		return c.JSON(500, map[string]interface{}{
			"status": false,
			"error":  result.Error.Error(),
		})
	}

	// Delete the image file (non-fatal if it fails)
	if car.ImageCar != "" {
		if err := deleteCarImage(car.ImageCar); err != nil {
			fmt.Printf("Warning: failed to delete image %q: %v\n", car.ImageCar, err)
		} else {
			fmt.Printf("Image %s deleted successfully \n", car.ImageCar)
		}
	}

	// if result.RowsAffected == 0 {
	// 	return c.JSON(404, map[string]interface{}{
	// 		"status": false,
	// 		"error":  "Car not found",
	// 	})
	// }

	return c.JSON(200, map[string]interface{}{
		"status":  true,
		"message": "Car deleted successfully!",
	})
}
