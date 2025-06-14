package routes

import (
	"automotiveApi/controllers"

	"github.com/labstack/echo/v4"
)

func InitRoute(e *echo.Echo) {
	e.POST("/cars", controllers.CreateCarsController)
	e.GET("/cars", controllers.GetCarsController)
}
