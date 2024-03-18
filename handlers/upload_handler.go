package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"file-upload-api/uploadmanager"

	"github.com/gin-gonic/gin"
)

var uploadManager *uploadmanager.UploadManager

func init() {
	uploadManager = uploadmanager.NewUploadManager()
	uploadManager.ProcessUploads()
}

// UploadFile handles the upload endpoint
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("File error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ext := filepath.Ext(file.Filename)
	if ext != ".csv" {
		c.AbortWithStatusJSON(400, gin.H{"error": "Only CSV files are allowed"})
		return
	}

	uploadID := uploadManager.Push(file)
	if uploadID == "" {
		log.Print("Failed to process upload")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process upload"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"id": uploadID})
}
