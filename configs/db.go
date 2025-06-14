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
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//panic(err)
		log.Fatal("failed to connect database", err)
	}
	return DB
}

func InitMigration(db *gorm.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS=0")

	// Migrate tables without foreign keys
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.Merek{},
		&models.Jenis{},
		&models.Merek{},
	)

	// Re-Enable Foreign Key Checks
	db.Exec("SET FOREIGN_KEY_CHECKS=1")

	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}
}
