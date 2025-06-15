package routes

import (
	"automotiveApi/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	e.Static("/", "frontend/views")

	e.POST("/api/cars", controllers.CreateCarsController)
	e.GET("/api/cars", controllers.GetCarsController)
}
