package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateCarsController(c echo.Context) error {
	var carsRequest models.Cars
	if err := c.Bind(&carsRequest); err != nil { // Cars Request
		return c.JSON(400, err.Error())
	}

	// Manually load relationships before saving
	if carsRequest.MerekID != nil {
		var merek models.Merek
		if err := configs.DB.First(&merek, *carsRequest.MerekID).Error; err == nil {
			carsRequest.Merek = &merek
		}
	}
	if carsRequest.JenisID != nil {
		var jenis models.Jenis
		if err := configs.DB.First(&jenis, *carsRequest.JenisID).Error; err == nil {
			carsRequest.Jenis = &jenis
		}
	}

	// Checking for MerekID availabality that connects with cars table
	if carsRequest.MerekID != nil {
		var merek models.Merek
		if err := configs.DB.First(&merek, *carsRequest.MerekID).Error; err != nil {
			return c.JSON(400, "Invalid brand ID")
		}
	}

	// Checking for JenisID availabality that connects with cars table
	if carsRequest.JenisID != nil {
		var jenis models.Jenis
		if err := configs.DB.First(&jenis, *carsRequest.JenisID).Error; err != nil {
			return c.JSON(400, "Invalid jenis ID")
		}
	}

	result := configs.DB.Create(&carsRequest) // Cars Request

	// Reload the car with relationships
	var newCar models.Cars
	configs.DB.
		Preload("Merek").
		Preload("Jenis").
		First(&newCar, carsRequest.CarsID)

	if result.Error != nil {
		//return c.JSON(500, result.Error.Error())
		return c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Message: result.Error.Error(),
			Status:  false,
			Data:    nil,
		})
	}
	//return c.JSON(http.StatusOK, carsRequest) // Cars Request
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menambahkan data",
		Status:  true,
		Data:    carsRequest,
	})
}

func UpdateCarController(c echo.Context) error {
	// Get ID from URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "false",
			"error":  "Invalid ID",
		})
	}

	// 1. Get Existing Cars, check if cars exists
	var car models.Cars
	if err := configs.DB.First(&car, id).Error; err != nil {
		return c.JSON(404, map[string]interface{}{
			"status": false,
			"error":  "Car not found",
		})
	}

	// 2. Bind updated data (partial updates allowed)
	//var updateData models.Cars
	/*if err := c.Bind(&car); err != nil {
		return c.JSON(400, err.Error())
	}*/
	if err := c.Bind(&car); err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	}

	// 3. Save Changes
	if err := configs.DB.Save(&car).Error; err != nil {
		return c.JSON(500, map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	}

	// Update only non-zero fields
	/*result := configs.DB.Model(&existingCars).Updates(updateData)
	if result.Error != nil {
		return c.JSON(500, result.Error.Error())
	}*/

	// 4. Force Reload with relationships
	var updatedCar models.Cars
	if err := configs.DB.
		Preload("Merek").
		Preload("Jenis").
		First(&updatedCar, id).Error; err != nil {
		return c.JSON(500, map[string]interface{}{
			"status": false,
			"error":  err.Error(),
		})
	}

	return c.JSON(200, map[string]interface{}{
		"status":  true,
		"message": "Car updated successfully",
		"data":    updatedCar, // Directly return the car object
	})
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
	}

	// Debug: Log the first car's relationships
	if len(cars) > 0 {
		log.Printf("First car relationships - Merek: %+v, Jenis: %+v",
			cars[0].Merek,
			cars[0].Jenis)
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
		return c.JSON(400, map[string]string{"error": "Invalid ID Format"})
	}

	var car models.Cars
	result := configs.DB.
		Preload("Merek").
		Preload("Jenis").
		First(&car, id) // Use First() for single records

	if result.Error != nil {
		return c.JSON(404, map[string]string{"error": "Car not found"})
	}
	return c.JSON(200, car)
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
		Message: "Berhasil menampilkan data",
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
		Message: "Berhasil menampilkan data",
		Status:  true,
		Data:    jenis,
	})
}
