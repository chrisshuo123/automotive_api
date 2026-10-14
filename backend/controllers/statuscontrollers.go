package controllers

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// For the statuscontrollers.go
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
