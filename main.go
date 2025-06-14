package main

import (
	"automotiveApi/configs"
	"automotiveApi/routes"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm/schema"
)

func main() {
	// Initialize Database
	db := configs.InitDatabase()
	fmt.Println("Database connection details:", configs.DB.Dialector.Name())

	// Disable automatic pluralization (Use when GORM Error Readings Occurs. In this case, is 'Merek' read as 'Mereks')
	db.Config.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   "",
		SingularTable: true, // <-- This is the key setting
		NoLowerCase:   true,
	}

	// Optional: Auto-migrate models (if needed)
	configs.InitMigration(db)

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
