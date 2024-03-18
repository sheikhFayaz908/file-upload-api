package handlers

import (
	"errors"
	"file-upload-api/database"
	"file-upload-api/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func DownloadHandler(c *gin.Context) {
	uploadID := c.Param("id")

	uploadData, err := database.FetchUploadedDataById(uploadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
			return
		}
		log.Printf("Error in fetching data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upload data"})
		return
	}

	switch uploadData.Status {
	case models.UploadStatusStarted:
		c.JSON(http.StatusConflict, gin.H{"error": "Upload is still processing"})
	case models.UploadStatusCompleted:
		c.Header("Content-Disposition", "attachment; filename=data.csv")
		c.Header("Content-Type", "text/csv")

		// Send file data in chunks
		chunkSize := 2
		for i := 0; i < len(uploadData.Data); i += chunkSize {
			end := i + chunkSize
			if end > len(uploadData.Data) {
				end = len(uploadData.Data)
			}
			var chunk strings.Builder
			for j := i; j < end; j++ {
				chunk.WriteString(uploadData.Data[j].ColumnValue + "\n")
			}
			c.Writer.WriteString(chunk.String())
			c.Writer.Flush()
		}
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown status"})
	}
}
