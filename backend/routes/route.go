package routes

import (
	"automotiveApi/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowOrigins:     []string{"*"}, // Allow all origins during development
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Static("/", "frontend/views")

	// Cars CRUD
	e.POST("/api/cars", controllers.CreateCarsController)
	e.GET("/api/cars", controllers.GetCarsController)    // List all cars
	e.GET("/api/cars/:id", controllers.GetCarController) // For single car, outmostly supporting the update func
	e.PUT("/api/cars/:id", controllers.UpdateCarController)
	e.DELETE("api/cars/:id", controllers.DeleteCarController)
	// Brands CRUD
	e.GET("/api/brands", controllers.GetMereksController)
	e.GET("/api/brands/:id", controllers.GetMerekController)
	// Types CRUD
	e.GET("/api/types", controllers.GetTypesController)
	e.GET("/api/types/:id", controllers.GetTypeController)
	// Status CRUD
	e.GET("/api/status", controllers.GetStatusController)
}
