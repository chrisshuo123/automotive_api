package main

import (
	"automotiveApi/configs"
	"automotiveApi/routes"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize Database
	db := configs.InitDatabase()

	// Optional: Auto-migrate models (if needed)
	configs.InitMigration(db)

	// Fix FKs without dropping tables
	if err := configs.FixForeignKeys(db); err != nil {
		log.Fatal("Failed to Fix Foreign Keys: ", err)
	}

	// Create Echo instance
	e := echo.New()

	// Register routes (example)
	/*e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})*/

	// Register Routes
	routes.InitRoute(e)

	//configs.DB.InitRoute()

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
