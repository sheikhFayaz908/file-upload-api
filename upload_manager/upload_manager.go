package uploadmanager

import (
	"file-upload-api/database"
	fileprocessor "file-upload-api/file_processor"
	"file-upload-api/models"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UploadManager struct {
	uploadQueue chan *UploadJob
	uploadsMap  map[string]*UploadJob
}

type UploadJob struct {
	ID       string
	File     *multipart.FileHeader
	Uploaded time.Time
}

func NewUploadManager() *UploadManager {
	return &UploadManager{
		uploadQueue: make(chan *UploadJob),
		uploadsMap:  make(map[string]*UploadJob),
	}
}

func (um *UploadManager) Push(file *multipart.FileHeader) string {
	uploadID := uuid.New().String()

	uploadJob := &UploadJob{
		ID:       uploadID,
		File:     file,
		Uploaded: time.Now(),
	}
	um.uploadsMap[uploadID] = uploadJob
	err := database.SaveJob(&models.Uploads{Status: models.UploadStatusStarted, ID: uploadID, FileName: file.Filename})
	if err != nil {
		log.Print(err)
		return ""
	}
	um.uploadQueue <- uploadJob
	return uploadID
}

func (um *UploadManager) ProcessUploads() {
	go func() {
		for job := range um.uploadQueue {
			go um.processUpload(job)
		}
	}()
}

func (um *UploadManager) processUpload(job *UploadJob) {
	processor, err := fileprocessor.GetFileProcessor("csv")
	if err != nil {
		return
	}

	data, err := processor.ProcessFile(job.File)
	if err != nil {
		database.UpdateJob(&models.Uploads{ID: job.ID, Status: models.UploadStatusError})
		return
	}
	batchSize := 100
	dataBatch := make([]*models.UploadedData, 0, batchSize)

	for idx, columnValue := range data {
		dataBatch = append(dataBatch, &models.UploadedData{ColumnValue: columnValue, UploadsId: job.ID})

		if len(dataBatch) == batchSize || idx == len(data)-1 {
			err := database.SaveUploadedData(dataBatch)
			if err != nil {
				log.Printf("Error saving batch: %v", err)
				database.UpdateJob(&models.Uploads{ID: job.ID, Status: models.UploadStatusError})
				return
			}
			// Clear the batch for the next iteration
			dataBatch = make([]*models.UploadedData, 0, batchSize)
		}
	}

	database.UpdateJob(&models.Uploads{ID: job.ID, Status: models.UploadStatusCompleted})
}
