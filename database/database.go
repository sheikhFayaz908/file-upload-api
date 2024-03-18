package database

import (
	"file-upload-api/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&models.Uploads{}, &models.UploadedData{})
	DB = db
	return DB
}

func SaveJob(job *models.Uploads) error {
	err := DB.Create(job).Error
	if err != nil {
		log.Print("DB error", err)
	}
	return err
}
func UpdateJob(job *models.Uploads) error {
	err := DB.Save(job).Error
	if err != nil {
		log.Print("DB error", err)
	}
	return err
}

func SaveUploadedData(data []*models.UploadedData) error {
	err := DB.Create(data).Error
	if err != nil {
		log.Print("DB error", err)
	}
	return err
}

func FetchUploadedDataById(id string) (*models.Uploads, error) {
	var uploadData models.Uploads
	err := DB.Preload("Data").First(&uploadData, "id = ?", id).Error
	if err != nil {
		log.Print("DB error", err)
	}
	return &uploadData, err
}
