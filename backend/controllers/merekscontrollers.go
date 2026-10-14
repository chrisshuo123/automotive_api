package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// For the merekscontrollers.go
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
