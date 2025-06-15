package configs

import (
	"automotiveApi/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/automotive_api?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		//panic(err)
		log.Fatal("failed to connect database", err)
	}
	return DB
}

func InitMigration(db *gorm.DB) {
	// Migrate tables without foreign keys
	err := db.AutoMigrate(
		&models.Merek{},
		&models.Jenis{},
		&models.Cars{},
	)

	if err != nil {
		log.Println("Migration Warning (non-fatal): ", err)
	}
}
