package database

import (
	"BLOG_APP/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=123456 dbname=BLOG_APP_DB port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	// Auto-migration des modèles
	err = DB.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{}, &models.Favorite{}, &models.Follow{})
	if err != nil {
		return err
	}

	return nil
}
