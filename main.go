package main

import (
	"context"
	"file-upload-api/database"
	"file-upload-api/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Set maximum request body size
	r.MaxMultipartMemory = 32 << 20

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	database.InitDB()
	v1 := r.Group("/api")
	{
		v1.POST("/upload", handlers.UploadFile)
		v1.GET("/download/:id", handlers.DownloadHandler)
	}
	//Health Endpoint
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]bool{"alive": true})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	shutDown(server)
}

func shutDown(server *http.Server) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server shutdown complete.")
}
