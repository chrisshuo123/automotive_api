package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateCarsController(c echo.Context) error {
	var carsRequest models.Cars
	c.Bind(&carsRequest) // Cars Request

	result := configs.DB.Create(&carsRequest) // Cars Request
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

func GetCarsController(c echo.Context) error {
	var cars []models.Cars

	// Set Response Header
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

	// Debug GORM Error output Mitigation
	fmt.Println("GORM is using table:", configs.DB.NamingStrategy.TableName("Merek"))

	result := configs.DB.
		Preload("Merek").
		Preload("Jenis").
		Find(&cars)

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

	//return c.JSON(http.StatusOK, cars)
	return c.JSON(http.StatusOK, models.BaseResponse{
		Message: "Berhasil menampilkan data",
		Status:  true,
		Data:    cars,
	})
}
