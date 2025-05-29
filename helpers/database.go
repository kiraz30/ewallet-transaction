package helpers

import (
	"ewallet-transaction/internal/models"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupMySQL() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", GetEnv("DB_USER", "root"), GetEnv("DB_PASSWORD", ""), GetEnv("DB_HOST", "localhost"), GetEnv("DB_PORT", "3306"), GetEnv("DB_NAME", "ewallet"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logrus.Info("Successfully connected to database")

	DB.AutoMigrate(&models.Transaction{})
}
