package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// For the typescontrollers.go
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
