package models

import "github.com/jinzhu/gorm"

type UploadStatus string

const (
	UploadStatusStarted   UploadStatus = "STARTED"
	UploadStatusCompleted UploadStatus = "COMPLETED"
	UploadStatusError     UploadStatus = "ERROR"
)

// Table: uploads
type Uploads struct {
	ID       string `gorm:"primary_key"`
	Status   UploadStatus
	FileName string         `gorm:"size:255"`
	Data     []UploadedData //one to many Relationship with uploaded_data Table
}

// Table: uploaded_data
type UploadedData struct {
	gorm.Model
	UploadsId   string
	ColumnValue string `gorm:"size:255"`
}
