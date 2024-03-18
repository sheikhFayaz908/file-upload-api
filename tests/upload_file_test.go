package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"file-upload-api/database"
	"file-upload-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var uploadId string

func TestUploadHandler(t *testing.T) {
	router := gin.New()
	database.InitDB()
	router.POST("/api/upload", handlers.UploadFile)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the file field to the multipart request
	part, err := writer.CreateFormFile("file", "example.csv")
	if err != nil {
		t.Fatal(err)
	}
	//records with varying fields and quoting
	part.Write([]byte(`Name,Description,Price
	"Product 1","This is a multi-line description that spans across
	multiple lines. It includes details about the features
	and benefits of the product.","$99.99"
	"Product 2","Another product with a detailed description
	that is split across lines due to its length.","$149.99"
	"Product 3","A short description.","$49.99"`))

	writer.Close()

	req, err := http.NewRequest("POST", "/api/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	//For Async process to complete
	time.Sleep(time.Second)
	// Assert the response status code
	var responseBody map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	assert.Equal(t, http.StatusAccepted, w.Code)
	uploadId = responseBody["id"]
}

func TestDownloadHandler(t *testing.T) {
	router := gin.New()
	database.InitDB()
	router.GET("/api/download/:id", handlers.DownloadHandler)

	req, err := http.NewRequest("GET", "/api/download/"+uploadId, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, "attachment; filename=data.csv", w.Header().Get("Content-Disposition"))
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))

}
