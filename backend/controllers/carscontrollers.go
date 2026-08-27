package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateCarsController(c echo.Context) error {
	var carsRequest models.Cars

	carsRequest.NamaMobil = c.FormValue("nama_mobil")

	if hp, err := strconv.Atoi(c.FormValue("horse_power")); err == nil {
		carsRequest.HorsePower = uint(hp)
	}

	if v := c.FormValue("idMerek_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			carsRequest.MerekID = &uid
		}
	}
	if v := c.FormValue("idJenis_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			carsRequest.JenisID = &uid
		}
	}
	if v := c.FormValue("idStatus_fk"); v != "" {
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

	// Optional image upload
	fileHeader, err := c.FormFile("image")
	fmt.Println("FormFile err: ", err)
	if err == nil {
		fmt.Println("Got file: ", fileHeader.Filename)
		fileName, err := saveUploadedImage(fileHeader)
		if err != nil {
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
	if v := c.FormValue("idMerek_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			car.MerekID = &uid
		}
	}
	if v := c.FormValue("idJenis_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			car.JenisID = &uid
		}
	}
	if v := c.FormValue("idStatus_fk"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			uid := uint(id)
			car.StatusID = &uid
		}
	}

	// Only replace the image if a new file was actually sent
	fileHeader, err := c.FormFile("image")
	if err == nil {
		fileName, err := saveUploadedImage(fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BaseResponse{
				Message: err.Error(),
				Status:  false,
			})
		}
		car.ImageCar = fileName
	}

	// 3. Save Changes
	if err := configs.DB.Save(&car).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: err.Error(), Status: false,
		})
	}

	// 4. Force Reload with relationships
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
	var cars []models.Cars

	// Set Response Header
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

	// Debug GORM Error output Mitigation
	fmt.Println("GORM is using table:", configs.DB.NamingStrategy.TableName("Merek"))

	// First load cars with relationships
	result := configs.DB.
		Preload("Merek", func(db *gorm.DB) *gorm.DB {
			return db.Select("idMerek, merek") // Only load necessary fields
		}).
		Preload("Jenis", func(db *gorm.DB) *gorm.DB {
			return db.Select("idJenis, jenis")
		}).
		Preload("Status", func(db *gorm.DB) *gorm.DB {
			return db.Select("idStatus, status")
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

	// var car []models.Cars // Single record, not slice
	var car models.Cars // Single struct, NOT slice

	// idCars := c.QueryParam("idCars")

	// Build query with Preload and filter by ID
	result := configs.DB.
		Preload("Merek").
		Preload("Jenis").
		Preload("Status").
		First(&car, id) // Use First() for single records

	// result := query.Find(&cars)

	if result.Error != nil {
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

func GetMerekController(c echo.Context) error {
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

func GetJenisController(c echo.Context) error {
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

	result := configs.DB.Delete(&models.Cars{}, id)
	if result.Error != nil {
		return c.JSON(500, map[string]interface{}{
			"status": false,
			"error":  result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(404, map[string]interface{}{
			"status": false,
			"error":  "Car not found",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"status":  true,
		"message": "Car deleted successfully!",
	})
}
