package main

import (
	"automotiveApi/configs"
	"automotiveApi/models"
	"automotiveApi/routes"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize Database
	db := configs.InitDatabase()
	fmt.Println("Database connection details:", configs.DB.Dialector.Name())

	// Optional: Auto-migrate models (if needed)
	configs.InitMigration(db)

	if err := db.Migrator().DropTable(&models.Cars{}); err == nil {
		db.AutoMigrate(&models.Cars{})
	}

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
