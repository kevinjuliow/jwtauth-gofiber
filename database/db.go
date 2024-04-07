package database

import (
	"github.com/joho/godotenv"
	"github.com/kevinjuliow/jwtauth-gofiber/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func OpenConnection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db_username := os.Getenv("DB_USERNAME")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	dsn := db_username + ":" + db_pass + "@tcp(127.0.0.1:3306)/?" + db_name + "charset=utf8mb4&parseTime=True&loc=Local"
	db, errDB := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDB != nil {
		panic(errDB)
	}
	if err := db.AutoMigrate(&models.Users{}); err != nil {
		panic(err)
	}
}
